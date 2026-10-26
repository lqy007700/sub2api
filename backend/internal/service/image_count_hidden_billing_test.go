//go:build unit

package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageCountBillingHiddenRateSettlementAndPublicSnapshot(t *testing.T) {
	for _, gateway := range []string{"gateway", "openai"} {
		for _, independent := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/independent=%v", gateway, independent), func(t *testing.T) {
				group := newPeakGroup(true, "14:00", "18:00", 3)
				group.ID, group.Platform, group.RateMultiplier = 906, PlatformOpenAI, 1.25
				group.Status, group.Hydrated = StatusActive, true
				group.HiddenRateMultiplier = testPtrFloat64(2)
				group.AllowImageGeneration, group.ImageCountBilling = true, true
				group.ImagePrice1K = testPtrFloat64(0.2)
				group.ImageRateIndependent, group.ImageRateMultiplier = independent, 0.4
				group.ModelPricing = []ChannelModelPricing{{
					Models: []string{"gpt-image-2"}, BillingMode: BillingModeToken,
					InputPrice: testPtrFloat64(3e-6), OutputPrice: testPtrFloat64(15e-6),
				}}
				user := &User{ID: 601}
				key := &APIKey{ID: 501, GroupID: &group.ID, Group: group, User: user, Quota: 100}
				usageRepo := &openAIRecordUsageLogRepoStub{inserted: true}
				billingRepo := &openAIRecordUsageBillingRepoStub{}
				userRepo := &openAIRecordUsageUserRepoStub{}
				rateRepo := &openAIUserGroupRateRepoStub{rate: testPtrFloat64(0.75)}
				resolver := newOpenAITokenImageChannelPricingResolverForTest(t, group.ID, "gpt-image-2")

				if gateway == "gateway" {
					group.Platform = PlatformComposite
					svc := newGatewayRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, &openAIRecordUsageSubRepoStub{})
					svc.resolver, svc.billingService = resolver, resolver.billingService
					svc.userGroupRateResolver = newUserGroupRateResolver(rateRepo, nil, resolveUserGroupRateCacheTTL(svc.cfg), nil, "hidden-count-test")
					require.NoError(t, svc.RecordUsage(context.Background(), &RecordUsageInput{
						Result: &ForwardResult{RequestID: "hidden-count-gateway", Model: "gpt-image-2", ImageCount: 3, ImageSize: "1K", ImageEndpoint: true,
							Usage: ClaudeUsage{InputTokens: 1_000_000, OutputTokens: 10_000}},
						APIKey: key, User: user, Account: &Account{ID: 701, Platform: PlatformOpenAI}, PricingAt: at(14, 30),
						APIKeyService: &openAIRecordUsageAPIKeyQuotaStub{},
					}))
				} else {
					svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, &openAIRecordUsageSubRepoStub{}, rateRepo)
					svc.resolver, svc.billingService = resolver, resolver.billingService
					require.NoError(t, svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
						Result: &OpenAIForwardResult{RequestID: "hidden-count-openai", Model: "gpt-image-2", ImageCount: 3, ImageSize: "1K", ImageEndpoint: true,
							Usage: OpenAIUsage{InputTokens: 1_000_000, OutputTokens: 10_000}},
						APIKey: key, User: user, Account: &Account{ID: 701, Platform: PlatformOpenAI}, PricingAt: at(14, 30),
						APIKeyService: &openAIRecordUsageAPIKeyQuotaStub{},
					}))
				}

				wantRate, wantDisplay := 2.0, 0.75
				if independent {
					wantRate, wantDisplay = 0.4, 0.4
				}
				log := usageRepo.lastLog
				require.NotNil(t, log)
				require.Equal(t, 3, log.ImageCount)
				require.InDelta(t, 0.6, log.TotalCost, 1e-12, "three actual images override the token card")
				require.InDelta(t, wantRate, log.RateMultiplier, 1e-12, "image rates do not acquire the token peak factor")
				require.InDelta(t, 0.6*wantRate, log.ActualCost, 1e-12)
				require.NotNil(t, log.DisplayRateMultiplier)
				require.InDelta(t, wantDisplay, *log.DisplayRateMultiplier, 1e-12)
				require.NotNil(t, billingRepo.lastCmd)
				require.InDelta(t, QuantizeUsageBillingAmount(0.6*wantRate), billingRepo.lastCmd.BalanceCost, 1e-12)
				require.InDelta(t, QuantizeUsageBillingAmount(0.6*wantRate), billingRepo.lastCmd.APIKeyQuotaCost, 1e-12)
				group.RateMultiplier = 99
				require.InDelta(t, wantDisplay, *log.DisplayRateMultiplier, 1e-12, "display multiplier remains a request-time snapshot")
			})
		}
	}
}
