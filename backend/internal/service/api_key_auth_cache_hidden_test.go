package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIKeyAuthSnapshotHiddenRateRoundtripAndClear(t *testing.T) {
	svc := &APIKeyService{}
	apiKey := profitAuthTestAPIKey()
	hidden := 0.4
	apiKey.Group.HiddenRateMultiplier = &hidden
	apiKey.Group.ImageCountBilling = true
	apiKey.Group.ForceOpenAIUltrafast = true
	snapshot := svc.snapshotFromAPIKey(context.Background(), apiKey)
	require.NotNil(t, snapshot.Group.HiddenRateMultiplier)
	// Snapshot ownership is independent of subsequent in-memory edits.
	hidden = 0.9
	require.Equal(t, 0.4, *snapshot.Group.HiddenRateMultiplier)
	payload, err := json.Marshal(&APIKeyAuthCacheEntry{Snapshot: snapshot})
	require.NoError(t, err)
	var entry APIKeyAuthCacheEntry
	require.NoError(t, json.Unmarshal(payload, &entry))
	restored, used, err := svc.applyAuthCacheEntry(apiKey.Key, &entry)
	require.NoError(t, err)
	require.True(t, used)
	require.Equal(t, 0.4, restored.Group.BillingRateMultiplier(0.2))
	require.Equal(t, apiKey.Group.RateMultiplier, restored.Group.RateMultiplier)
	require.True(t, restored.Group.ImageCountBilling)
	require.True(t, restored.Group.ForceOpenAIUltrafast)
	*restored.Group.HiddenRateMultiplier = 0.8
	require.Equal(t, 0.4, *entry.Snapshot.Group.HiddenRateMultiplier)

	apiKey.Group.HiddenRateMultiplier = nil
	cleared := svc.snapshotFromAPIKey(context.Background(), apiKey)
	restored, used, err = svc.applyAuthCacheEntry(apiKey.Key, &APIKeyAuthCacheEntry{Snapshot: cleared})
	require.NoError(t, err)
	require.True(t, used)
	require.Nil(t, restored.Group.HiddenRateMultiplier)
	require.Equal(t, 0.2, restored.Group.BillingRateMultiplier(0.2))

	// Fork and upstream snapshots predate the combined billing and Ultra Fast settings.
	for _, legacyVersion := range []int{24, 25, 26, 27, 28} {
		cleared.Version = legacyVersion
		restored, used, err = svc.applyAuthCacheEntry(apiKey.Key, &APIKeyAuthCacheEntry{Snapshot: cleared})
		require.NoError(t, err)
		require.False(t, used)
		require.Nil(t, restored)
	}
}
