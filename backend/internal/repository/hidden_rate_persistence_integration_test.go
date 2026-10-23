//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func (s *GroupRepoSuite) TestHiddenRateMultiplierRoundTripAndAuthProjection() {
	hidden := 3.5
	group := &service.Group{Name: "hidden-rate-roundtrip", Platform: service.PlatformOpenAI, RateMultiplier: 1.25, HiddenRateMultiplier: &hidden, Status: service.StatusActive, SubscriptionType: service.SubscriptionTypeStandard}
	s.Require().NoError(s.repo.Create(s.ctx, group))
	got, err := s.repo.GetByID(s.ctx, group.ID)
	s.Require().NoError(err)
	s.Require().Equal(hidden, *got.HiddenRateMultiplier)
	s.Require().Equal(1.25, got.RateMultiplier)

	user := mustCreateUser(s.T(), s.tx.Client(), &service.User{Email: "hidden-auth@example.com"})
	apiKey := mustCreateApiKey(s.T(), s.tx.Client(), &service.APIKey{UserID: user.ID, GroupID: &group.ID, Key: "sk-hidden-auth", Name: "hidden", Status: service.StatusActive})
	keys := newAPIKeyRepositoryWithSQL(s.tx.Client(), s.tx)
	auth, err := keys.GetByKeyForAuth(s.ctx, apiKey.Key)
	s.Require().NoError(err)
	s.Require().NotNil(auth.Group)
	s.Require().Equal(hidden, *auth.Group.HiddenRateMultiplier)

	next := 4.5
	got.HiddenRateMultiplier = &next
	s.Require().NoError(s.repo.Update(s.ctx, got))
	got, err = s.repo.GetByID(s.ctx, group.ID)
	s.Require().NoError(err)
	s.Require().Equal(next, *got.HiddenRateMultiplier)
	got.HiddenRateMultiplier = nil
	s.Require().NoError(s.repo.Update(s.ctx, got))
	got, err = s.repo.GetByID(s.ctx, group.ID)
	s.Require().NoError(err)
	s.Require().Nil(got.HiddenRateMultiplier)
	s.Require().Equal(1.25, got.RateMultiplier)
}

func (s *UsageLogRepoSuite) TestDisplayRateSnapshotRoundTripAndHistoricalNull() {
	user := mustCreateUser(s.T(), s.client, &service.User{Email: "display-rate@example.com"})
	apiKey := mustCreateApiKey(s.T(), s.client, &service.APIKey{UserID: user.ID, Key: "sk-display-rate", Name: "display"})
	account := mustCreateAccount(s.T(), s.client, &service.Account{Name: "display-rate-account"})
	display := 1.25
	for _, snapshot := range []*float64{nil, &display} {
		log := &service.UsageLog{UserID: user.ID, APIKeyID: apiKey.ID, AccountID: account.ID, Model: "gpt-5", RateMultiplier: 3.5, DisplayRateMultiplier: snapshot, InputTokens: 1}
		_, err := s.repo.Create(s.ctx, log)
		s.Require().NoError(err)
		got, err := s.repo.GetByID(s.ctx, log.ID)
		s.Require().NoError(err)
		s.Require().Equal(3.5, got.RateMultiplier)
		s.Require().Equal(snapshot, got.DisplayRateMultiplier)
	}
}

func TestBatchImageJobDisplaySnapshotRoundTrip(t *testing.T) {
	ctx := context.Background()
	repo := newBatchImageRepositoryWithSQL(testTx(t))
	display := 0.625
	for _, snapshot := range []*float64{nil, &display} {
		job, err := repo.CreateBatchImageJob(ctx, service.CreateBatchImageJobParams{
			UserID: 1001, Provider: service.BatchImageProviderGeminiAPI, Model: "gemini-2.5-flash-image", ItemCount: 1,
			GroupRateMultiplier: 3.5, DisplayRateMultiplier: snapshot,
		})
		require.NoError(t, err)
		require.Equal(t, snapshot, job.DisplayRateMultiplier)
		got, err := repo.GetBatchImageJobByBatchID(ctx, job.BatchID)
		require.NoError(t, err)
		require.Equal(t, 3.5, got.GroupRateMultiplier)
		require.Equal(t, snapshot, got.DisplayRateMultiplier)
	}
}
