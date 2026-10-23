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
	*restored.Group.HiddenRateMultiplier = 0.8
	require.Equal(t, 0.4, *entry.Snapshot.Group.HiddenRateMultiplier)

	apiKey.Group.HiddenRateMultiplier = nil
	cleared := svc.snapshotFromAPIKey(context.Background(), apiKey)
	restored, used, err = svc.applyAuthCacheEntry(apiKey.Key, &APIKeyAuthCacheEntry{Snapshot: cleared})
	require.NoError(t, err)
	require.True(t, used)
	require.Nil(t, restored.Group.HiddenRateMultiplier)
	require.Equal(t, 0.2, restored.Group.BillingRateMultiplier(0.2))

	// v24 did not carry hidden billing configuration and must always reload.
	cleared.Version = 24
	restored, used, err = svc.applyAuthCacheEntry(apiKey.Key, &APIKeyAuthCacheEntry{Snapshot: cleared})
	require.NoError(t, err)
	require.False(t, used)
	require.Nil(t, restored)
}
