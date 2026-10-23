//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHiddenRateBatch_SubmissionAndSettlementSnapshots(t *testing.T) {
	for _, tt := range []struct {
		name        string
		hidden      *float64
		userRate    *float64
		imageRate   *float64
		wantActual  float64
		wantDisplay float64
	}{
		{"public default", nil, nil, nil, 1.25, 0.625},
		{"public user override", nil, testPtrFloat64(0.75), nil, 0.75, 0.375},
		{"hidden overrides default", testPtrFloat64(2), nil, nil, 2, 0.625},
		{"hidden overrides user", testPtrFloat64(2), testPtrFloat64(0.75), nil, 2, 0.375},
		{"hidden zero", testPtrFloat64(0), testPtrFloat64(0.75), nil, 0, 0.375},
		{"hidden overrides free public", testPtrFloat64(2), testPtrFloat64(0), nil, 2, 0},
		{"independent image overrides hidden and user", testPtrFloat64(2), testPtrFloat64(0.75), testPtrFloat64(0.4), 0.4, 0.2},
		{"independent image overrides free hidden", testPtrFloat64(0), testPtrFloat64(0.75), testPtrFloat64(0.4), 0.4, 0.2},
		{"independent free image overrides hidden", testPtrFloat64(2), testPtrFloat64(0.75), testPtrFloat64(0), 0, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc, fakeRepo, queue, provider, _ := newTestBatchImagePublicService(true)
			repo := &hiddenBatchSnapshotRepository{fakeBatchImageRepository: fakeRepo}
			svc.Repo = repo
			group := &Group{
				ID: 907, Platform: PlatformGemini, RateMultiplier: 1.25, HiddenRateMultiplier: tt.hidden,
				AllowImageGeneration: true, AllowBatchImageGeneration: true,
				BatchImageDiscountMultiplier: 0.5, BatchImageHoldMultiplier: 0.6,
			}
			if tt.imageRate != nil {
				group.ImageRateIndependent, group.ImageRateMultiplier = true, *tt.imageRate
			}
			svc.GroupRepo = &publicBatchImageGroupRepo{groups: map[int64]*Group{group.ID: group}}
			userRates := &publicBatchImageUserGroupRateRepo{rates: map[int64]*float64{group.ID: tt.userRate}}
			svc.UserGroupRateRepo = userRates
			accountRepo := svc.AccountRepo.(*publicBatchImageAccountRepo)
			accountRepo.accounts[0].RateMultiplier = testPtrFloat64(1.6)
			owner := testBatchImageOwner()
			owner.GroupID = &group.ID

			response, err := svc.Submit(ctx, owner, validBatchImageSubmitRequest(), "")
			require.NoError(t, err)
			require.Len(t, provider.submits, 1)
			require.Equal(t, []string{response.ID}, queue.enqueued)
			job := repo.jobs[response.ID]
			require.NotNil(t, job)
			require.Equal(t, BatchImageJobStatusSubmitted, job.Status)
			require.Equal(t, 1, job.PricingSnapshotVersion)
			require.Equal(t, 2, job.ItemCount)
			require.InDelta(t, tt.wantActual, job.GroupRateMultiplier, 1e-12)
			require.InDelta(t, 1.6, job.AccountRateMultiplier, 1e-12)
			require.InDelta(t, 0.5, job.BatchDiscountMultiplier, 1e-12)
			require.NotNil(t, job.DisplayRateMultiplier)
			require.InDelta(t, tt.wantDisplay, *job.DisplayRateMultiplier, 1e-12, "the public snapshot includes the batch discount and excludes the account multiplier")
			// Base price .25 × account rate 1.6 × discount .5 = .2 per actual group rate.
			require.InDelta(t, 0.2*tt.wantActual, job.BillableUnitPrice, 1e-12)
			require.InDelta(t, 0.4*tt.wantActual, response.EstimatedCost, 1e-12)
			require.InDelta(t, 0.4*tt.wantActual, job.EstimatedCost, 1e-12)
			require.NotNil(t, job.HoldAmount)
			require.InDelta(t, 0.48*tt.wantActual, response.HoldAmount, 1e-12)
			require.InDelta(t, 0.48*tt.wantActual, *job.HoldAmount, 1e-12)
			billing := svc.BillingRepo.(*fakeBatchImageBillingRepo)
			if tt.wantActual == 0 {
				require.Empty(t, billing.reserves, "a free actual rate needs no balance hold")
			} else {
				require.Len(t, billing.reserves, 1)
				require.InDelta(t, 0.48*tt.wantActual, billing.reserves[0].HoldAmount, 1e-12, "the hold follows the actual rate")
			}

			// Settlement must keep the submitted prices after every live rate changes.
			group.RateMultiplier, group.HiddenRateMultiplier = 99, testPtrFloat64(88)
			group.ImageRateIndependent, group.ImageRateMultiplier = true, 77
			group.BatchImageDiscountMultiplier, group.BatchImageHoldMultiplier = 0.9, 0.95
			userRates.rates[group.ID] = testPtrFloat64(66)
			accountRepo.accounts[0].RateMultiplier = testPtrFloat64(55)
			job.Status, job.SuccessCount, job.FailCount = BatchImageJobStatusSettling, 1, 1
			usageRepo := &openAIRecordUsageLogRepoStub{inserted: true}
			settlement := &BatchImageSettlementService{
				Repo: repo, BillingRepo: billing, UsageLogRepo: usageRepo,
				Pricing: &fakeBatchImagePricingResolver{unitPrice: 999},
			}
			result, err := settlement.Settle(ctx, job.BatchID)
			require.NoError(t, err)
			require.InDelta(t, 0.2*tt.wantActual, result.ActualCost, 1e-12, "only the successful image is charged at the submitted rate")
			require.Equal(t, BatchImageJobStatusCompleted, job.Status)
			require.Len(t, billing.captures, 1)
			require.InDelta(t, 0.2*tt.wantActual, billing.captures[0].ActualAmount, 1e-12)
			require.InDelta(t, 0.48*tt.wantActual, billing.captures[0].HoldAmount, 1e-12)
			log := usageRepo.lastLog
			require.NotNil(t, log)
			require.Equal(t, 1, log.ImageCount)
			require.InDelta(t, 0.2*tt.wantActual, log.ActualCost, 1e-12)
			require.InDelta(t, tt.wantActual*0.5, log.RateMultiplier, 1e-12)
			require.NotNil(t, log.AccountRateMultiplier)
			require.InDelta(t, 1.6, *log.AccountRateMultiplier, 1e-12)
			require.NotNil(t, log.DisplayRateMultiplier)
			require.InDelta(t, tt.wantDisplay, *log.DisplayRateMultiplier, 1e-12)
			require.NotSame(t, job.DisplayRateMultiplier, log.DisplayRateMultiplier)
			*job.DisplayRateMultiplier = 44
			require.InDelta(t, tt.wantDisplay, *log.DisplayRateMultiplier, 1e-12, "usage owns a copy of the saved public snapshot")
		})
	}
}

func TestHiddenRateBatch_SettlementPreservesMissingLegacyDisplaySnapshot(t *testing.T) {
	repo := newFakeBatchImageRepository()
	job := testSettlingBatchImageJob("imgbatch_legacy_hidden_display")
	job.PricingSnapshotVersion = 1
	job.GroupRateMultiplier, job.BatchDiscountMultiplier, job.BillableUnitPrice = 2, 0.5, 0.4
	repo.jobs[job.BatchID] = job
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: true}
	svc := &BatchImageSettlementService{
		Repo: repo, BillingRepo: &fakeBatchImageBillingRepo{}, UsageLogRepo: usageRepo,
		Pricing: &fakeBatchImagePricingResolver{unitPrice: 999},
	}

	result, err := svc.Settle(context.Background(), job.BatchID)
	require.NoError(t, err)
	require.InDelta(t, 0.8, result.ActualCost, 1e-12)
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, 1.0, usageRepo.lastLog.RateMultiplier)
	require.Nil(t, usageRepo.lastLog.DisplayRateMultiplier, "a historical job cannot reconstruct a public snapshot from its actual rate")
}

// The shared fake predates this persisted field. Keep the adapter local to these
// tests so Submit must supply the snapshot before settlement can read it back.
type hiddenBatchSnapshotRepository struct {
	*fakeBatchImageRepository
}

func (r *hiddenBatchSnapshotRepository) CreateBatchImageJob(ctx context.Context, params CreateBatchImageJobParams) (*BatchImageJob, error) {
	job, err := r.fakeBatchImageRepository.CreateBatchImageJob(ctx, params)
	if err == nil && params.DisplayRateMultiplier != nil {
		display := *params.DisplayRateMultiplier
		job.DisplayRateMultiplier = &display
	}
	return job, err
}
