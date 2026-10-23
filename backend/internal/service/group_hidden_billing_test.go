//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/stretchr/testify/require"
)

func TestHiddenRateBilling_RecordUsageChargeAndPublicSnapshot(t *testing.T) {
	for _, gateway := range []string{"gateway", "openai"} {
		for _, tt := range []struct {
			name       string
			hidden     *float64
			override   *float64
			wantActual float64
			wantPublic float64
		}{
			{"public default", nil, nil, 3.75, 3.75},
			{"public user override", nil, testPtrFloat64(0.75), 2.25, 2.25},
			{"hidden overrides default", testPtrFloat64(2), nil, 6, 3.75},
			{"hidden overrides user", testPtrFloat64(2), testPtrFloat64(0.75), 6, 2.25},
			{"hidden zero stays free", testPtrFloat64(0), testPtrFloat64(0.75), 0, 2.25},
			{"hidden overrides free public", testPtrFloat64(2), testPtrFloat64(0), 6, 0},
		} {
			t.Run(gateway+"/"+tt.name, func(t *testing.T) {
				group := newPeakGroup(true, "14:00", "18:00", 3)
				group.ID, group.RateMultiplier, group.HiddenRateMultiplier = 901, 1.25, tt.hidden
				apiKey := &APIKey{ID: 501, GroupID: &group.ID, Group: group, Quota: 100}
				user := &User{ID: 601}
				apiKey.User = user
				usageRepo := &openAIRecordUsageLogRepoStub{inserted: true}
				billingRepo := &openAIRecordUsageBillingRepoStub{}
				userRepo := &openAIRecordUsageUserRepoStub{}
				rateRepo := &openAIUserGroupRateRepoStub{rate: tt.override}
				pricingAt := at(14, 30)

				if gateway == "gateway" {
					group.Platform = PlatformAnthropic
					svc := newGatewayRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, &openAIRecordUsageSubRepoStub{})
					svc.userGroupRateResolver = newUserGroupRateResolver(rateRepo, nil, resolveUserGroupRateCacheTTL(svc.cfg), nil, "hidden-rate-test")
					err := svc.RecordUsage(context.Background(), &RecordUsageInput{
						Result: &ForwardResult{RequestID: "hidden-gateway", Model: "claude-sonnet-4", Usage: ClaudeUsage{InputTokens: 1000, OutputTokens: 500}},
						APIKey: apiKey, User: user, Account: &Account{ID: 701}, PricingAt: pricingAt, APIKeyService: &openAIRecordUsageAPIKeyQuotaStub{},
					})
					require.NoError(t, err)
					// A hidden override must never replace the public rate in the resolver cache.
					wantPublicBase := 1.25
					if tt.override != nil {
						wantPublicBase = *tt.override
					}
					require.Equal(t, wantPublicBase, svc.ResolveUserGroupRateMultiplier(context.Background(), user.ID, group.ID, group.RateMultiplier))
				} else {
					group.Platform = PlatformOpenAI
					svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, &openAIRecordUsageSubRepoStub{}, rateRepo)
					err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
						Result: &OpenAIForwardResult{RequestID: "hidden-openai", Model: "gpt-5.1", Usage: OpenAIUsage{InputTokens: 1000, OutputTokens: 500}},
						APIKey: apiKey, User: user, Account: &Account{ID: 701}, PricingAt: pricingAt, APIKeyService: &openAIRecordUsageAPIKeyQuotaStub{},
					})
					require.NoError(t, err)
				}

				log := usageRepo.lastLog
				require.NotNil(t, log)
				require.Greater(t, log.TotalCost, 0.0)
				require.InDelta(t, tt.wantActual, log.RateMultiplier, 1e-12)
				require.InDelta(t, log.TotalCost*tt.wantActual, log.ActualCost, 1e-12)
				require.NotNil(t, log.DisplayRateMultiplier)
				require.InDelta(t, tt.wantPublic, *log.DisplayRateMultiplier, 1e-12, "display uses the public rate and peak at request time")
				require.False(t, log.CreatedAt.Equal(pricingAt))
				require.NotNil(t, billingRepo.lastCmd)
				require.InDelta(t, QuantizeUsageBillingAmount(log.ActualCost), billingRepo.lastCmd.BalanceCost, 1e-12)
				require.InDelta(t, QuantizeUsageBillingAmount(log.ActualCost), billingRepo.lastCmd.APIKeyQuotaCost, 1e-12)
				require.Zero(t, userRepo.deductCalls, "atomic billing owns the debit")
				// The stored display value must remain a snapshot after group edits.
				group.RateMultiplier = 99
				require.InDelta(t, tt.wantPublic, *log.DisplayRateMultiplier, 1e-12)
			})
		}
	}
}

func TestHiddenRateBilling_MediaIndependentRatesTakePrecedence(t *testing.T) {
	for _, gateway := range []string{"gateway", "openai"} {
		for _, media := range []string{"image", "video"} {
			if gateway == "gateway" && media == "video" {
				continue // Video final billing is handled by the OpenAI gateway.
			}
			for _, independent := range []bool{false, true} {
				name := gateway + "/" + media + "/inherited"
				if independent {
					name = gateway + "/" + media + "/independent"
				}
				t.Run(name, func(t *testing.T) {
					group := newPeakGroup(true, "14:00", "18:00", 3)
					group.ID, group.Platform, group.RateMultiplier = 902, PlatformGrok, 1.25
					group.HiddenRateMultiplier = testPtrFloat64(2)
					group.ImagePrice1K, group.VideoPrice480P = testPtrFloat64(0.2), testPtrFloat64(0.08)
					group.ImageRateIndependent, group.VideoRateIndependent = independent, independent
					group.ImageRateMultiplier, group.VideoRateMultiplier = 0.4, 0.25
					user := &User{ID: 601}
					key := &APIKey{ID: 501, GroupID: &group.ID, Group: group, User: user}
					usageRepo := &openAIRecordUsageLogRepoStub{inserted: true}
					userRepo := &openAIRecordUsageUserRepoStub{}
					rateRepo := &openAIUserGroupRateRepoStub{rate: testPtrFloat64(0.75)}
					if gateway == "gateway" {
						svc := newGatewayRecordUsageServiceForTest(usageRepo, userRepo, &openAIRecordUsageSubRepoStub{})
						svc.userGroupRateResolver = newUserGroupRateResolver(rateRepo, nil, resolveUserGroupRateCacheTTL(svc.cfg), nil, "hidden-media-test")
						require.NoError(t, svc.RecordUsage(context.Background(), &RecordUsageInput{
							Result: &ForwardResult{RequestID: "hidden-image", Model: "gpt-image-2", ImageCount: 1, ImageSize: "1K"},
							APIKey: key, User: user, Account: &Account{ID: 701}, PricingAt: at(14, 30),
						}))
					} else {
						svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, &openAIRecordUsageSubRepoStub{}, rateRepo)
						result := &OpenAIForwardResult{RequestID: "hidden-image", Model: "gpt-image-2", ImageCount: 1, ImageSize: "1K"}
						if media == "video" {
							result = &OpenAIForwardResult{RequestID: "hidden-video", Model: "grok-imagine-video-1.5", VideoCount: 1, VideoResolution: VideoBillingResolution480P, VideoDurationSeconds: 1}
						}
						require.NoError(t, svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
							Result: result, APIKey: key, User: user, Account: &Account{ID: 701, Platform: PlatformGrok}, PricingAt: at(14, 30),
						}))
					}
					wantRate, wantDisplay, wantTotal := 2.0, 0.75, 0.2
					if media == "video" {
						wantTotal = 0.08
					}
					if independent {
						wantRate = group.ImageRateMultiplier
						if media == "video" {
							wantRate = group.VideoRateMultiplier
						}
						wantDisplay = wantRate
					}
					log := usageRepo.lastLog
					require.NotNil(t, log)
					require.InDelta(t, wantTotal, log.TotalCost, 1e-12)
					require.InDelta(t, wantRate, log.RateMultiplier, 1e-12, "per-request media does not acquire the token peak factor")
					require.NotNil(t, log.DisplayRateMultiplier)
					require.InDelta(t, wantDisplay, *log.DisplayRateMultiplier, 1e-12)
					require.InDelta(t, wantTotal*wantRate, log.ActualCost, 1e-12)
					require.InDelta(t, QuantizeUsageBillingAmount(log.ActualCost), userRepo.lastAmount, 1e-12)
				})
			}
		}
	}
}

func TestHiddenRateBilling_InflightEstimateUsesActualRate(t *testing.T) {
	for _, gateway := range []string{"gateway", "openai"} {
		for _, tt := range []struct {
			name   string
			hidden *float64
			want   float64
		}{
			{"public override", nil, 0.375},
			{"hidden override", testPtrFloat64(2), 1},
			{"hidden free", testPtrFloat64(0), 0},
		} {
			t.Run(gateway+"/"+tt.name, func(t *testing.T) {
				group := &Group{ID: 903, Platform: PlatformAnthropic, RateMultiplier: 1.25, HiddenRateMultiplier: tt.hidden,
					ModelPricing: []ChannelModelPricing{{Models: []string{"hidden-request"}, BillingMode: BillingModePerRequest, PerRequestPrice: testPtrFloat64(0.5)}},
				}
				key := &APIKey{GroupID: &group.ID, Group: group, User: &User{ID: 601}}
				rateRepo := &openAIUserGroupRateRepoStub{rate: testPtrFloat64(0.75)}
				var estimate func(context.Context, *APIKey, InflightEstimateRequest) (float64, bool)
				if gateway == "gateway" {
					svc := newInflightEstimateGateway(t, nil)
					svc.userGroupRateResolver = newUserGroupRateResolver(rateRepo, nil, resolveUserGroupRateCacheTTL(svc.cfg), nil, "hidden-inflight-test")
					estimate = svc.EstimateInflightReservation
				} else {
					svc := newOpenAIRecordUsageServiceForTest(nil, nil, nil, rateRepo)
					svc.resolver = NewModelPricingResolver(nil, svc.billingService)
					estimate = svc.EstimateInflightReservation
				}
				cost, priced := estimate(context.Background(), key, InflightEstimateRequest{Model: "hidden-request", BodyBytes: 400})
				require.True(t, priced)
				require.InDelta(t, tt.want, cost, 1e-12)
			})
		}
	}
}

func TestHiddenRateBilling_ProfitGateUsesBillingGroupOverride(t *testing.T) {
	for _, gateway := range []string{"gateway", "openai"} {
		for _, composite := range []bool{false, true} {
			for _, hidden := range []float64{0, 2} {
				name := gateway + "/direct"
				if composite {
					name = gateway + "/composite"
				}
				if hidden == 0 {
					name += "/free"
				}
				t.Run(name, func(t *testing.T) {
					billingGroup := profitControlTestGroup(904, 0.2, 0.05)
					billingGroup.RateMultiplier, billingGroup.HiddenRateMultiplier = 1.25, &hidden
					billingGroup.SubscriptionType = SubscriptionTypeSubscription
					billingGroup.PeakRateEnabled, billingGroup.PeakStart, billingGroup.PeakEnd, billingGroup.PeakRateMultiplier = true, "14:00", "18:00", 3
					member := billingGroup
					if composite {
						billingGroup.Platform = PlatformComposite
						member = profitControlTestGroup(905, 0.2, 0.05)
						member.RateMultiplier, member.HiddenRateMultiplier = 99, testPtrFloat64(88)
					}
					ctx := context.WithValue(profitControlTestCtx(billingGroup), ctxkey.UserID, int64(601))
					ctx, _ = WithGatewayTokenRequestPricing(ctx)
					ctx = context.WithValue(ctx, openAIPricingAtCtxKey{}, at(14, 30))
					ctx = context.WithValue(ctx, gatewayTokenRequestPricingAtCtxKey{}, at(14, 30))
					rateRepo := &openAIUserGroupRateRepoStub{rate: testPtrFloat64(0.1)}
					var gate *openAIProfitControlGate
					if gateway == "gateway" {
						if composite {
							// Routing may replace the context group after billing ownership was frozen.
							ctx = context.WithValue(ctx, ctxkey.Group, member)
						}
						svc := &GatewayService{schedulerSnapshot: &SchedulerSnapshotService{groupRepo: profitControlGroupRepo{group: member}}}
						svc.userGroupRateResolver = newUserGroupRateResolver(rateRepo, nil, time.Minute, nil, "hidden-profit-test")
						ctx = svc.withGatewayProfitControlGate(ctx, &member.ID)
						gate, _ = ctx.Value(openAIProfitControlGateCtxKey{}).(*openAIProfitControlGate)
					} else {
						svc := newOpenAIRecordUsageServiceForTest(nil, nil, nil, rateRepo)
						svc.schedulerSnapshot = &SchedulerSnapshotService{groupRepo: profitControlGroupRepo{group: member}}
						gate = svc.resolveOpenAIProfitControlGate(ctx, &member.ID)
						ctx = context.WithValue(ctx, openAIProfitControlGateCtxKey{}, gate)
					}
					require.NotNil(t, gate)
					require.Equal(t, member.ID, gate.groupID)
					require.InDelta(t, hidden*3*0.75, gate.threshold, 1e-12)
					account := &Account{ID: 701, RateMultiplier: testPtrFloat64(0.5)}
					vetoed, _ := openAIProfitControlVetoReason(ctx, account)
					require.Equal(t, hidden == 0, vetoed, "admission follows the actual charge rather than the public user override or member hidden rate")
				})
			}
		}
	}
}

func TestHiddenRateBilling_ProfitPreviewIgnoresPublicOverridesWhenHidden(t *testing.T) {
	for _, tt := range []struct {
		name        string
		hidden      *float64
		wantDefault float64
		wantMin     float64
	}{
		{"public rates", nil, 3.75, 0.3},
		{"hidden rate", testPtrFloat64(2), 6, 6},
		{"hidden zero", testPtrFloat64(0), 0, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			group := profitControlTestGroup(906, 0.2, 0.05)
			group.RateMultiplier, group.HiddenRateMultiplier = 1.25, tt.hidden
			group.SubscriptionType = SubscriptionTypeSubscription
			group.PeakRateEnabled, group.PeakStart, group.PeakEnd, group.PeakRateMultiplier = true, "14:00", "18:00", 3
			reports := PreviewProfitAdmission([]ProfitPreviewGroupInput{{
				Group: group, UserOverrides: map[int64]float64{601: 0.1, 602: 8},
				Accounts: []*Account{{ID: 701, RateMultiplier: testPtrFloat64(0.5)}}, Models: []string{"gpt-5.1"},
			}}, at(14, 30))
			require.Len(t, reports, 1)
			report := reports[0]
			require.True(t, report.EffectiveGate)
			require.InDelta(t, tt.wantDefault, report.DefaultD, 1e-12)
			require.InDelta(t, tt.wantMin, report.MinEffectiveD, 1e-12)
			require.InDelta(t, tt.wantDefault*0.75, report.ThresholdDefault, 1e-12)
			require.InDelta(t, tt.wantMin*0.75, report.ThresholdMinD, 1e-12)
			if tt.wantDefault == 0 {
				require.Equal(t, ProfitPreviewClassRejectedThreshold, report.Verdicts[0].Class)
			} else {
				require.Equal(t, ProfitPreviewClassAdmitted, report.Verdicts[0].Class)
				require.Equal(t, tt.hidden == nil, report.Verdicts[0].RejectedUnderMinD)
			}
		})
	}
}
