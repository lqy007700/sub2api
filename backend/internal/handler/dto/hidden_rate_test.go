package dto

import (
	"encoding/json"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestHiddenGroupRateAdminOnly(t *testing.T) {
	hidden := 0.375
	group := &service.Group{ID: 10, RateMultiplier: 2, HiddenRateMultiplier: &hidden}
	for _, userGroup := range []*Group{GroupFromService(group), GroupFromServiceShallow(group)} {
		payload, err := json.Marshal(userGroup)
		require.NoError(t, err)
		require.NotContains(t, string(payload), "hidden_rate_multiplier")
		require.Contains(t, string(payload), `"rate_multiplier":2`)
	}
	admin := GroupFromServiceAdmin(group)
	require.Equal(t, &hidden, admin.HiddenRateMultiplier)
	payload, err := json.Marshal(admin)
	require.NoError(t, err)
	require.Contains(t, string(payload), `"hidden_rate_multiplier":0.375`)
	group.HiddenRateMultiplier = nil
	payload, err = json.Marshal(GroupFromServiceAdmin(group))
	require.NoError(t, err)
	require.Contains(t, string(payload), `"hidden_rate_multiplier":null`)
}

func TestUsageDisplaySnapshotKeepsActualCostsAndAdminRate(t *testing.T) {
	for _, display := range []float64{2, 0} {
		log := &service.UsageLog{RateMultiplier: 0.375, DisplayRateMultiplier: &display, ActualCost: 0.00375, TotalCost: 0.01,
			Group: &service.Group{RateMultiplier: 99}}
		user := UsageLogFromService(log)
		admin := UsageLogFromServiceAdmin(log)
		require.Equal(t, display, user.RateMultiplier, "use the request snapshot even when current rates differ")
		require.Equal(t, 0.375, admin.RateMultiplier)
		require.Equal(t, log.ActualCost, user.ActualCost)
		require.Equal(t, log.TotalCost, user.TotalCost)
		payload, err := json.Marshal(user)
		require.NoError(t, err)
		require.NotContains(t, string(payload), "display_rate_multiplier")
		require.NotContains(t, string(payload), "hidden_rate_multiplier")
	}
	legacy := &service.UsageLog{RateMultiplier: 1.7}
	require.Equal(t, 1.7, UsageLogFromService(legacy).RateMultiplier)
	require.Equal(t, 1.7, UsageLogFromServiceAdmin(legacy).RateMultiplier)
}
