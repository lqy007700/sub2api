package service

import (
	"context"
	"time"
)

// publicUsageRateMultiplier captures the rate shown to the user at pricing time.
// Billing amounts and the internal RateMultiplier always retain the actual rate.
// The resolver only receives public rates, so billing overrides never enter its cache.
func publicUsageRateMultiplier(ctx context.Context, apiKey *APIKey, userID int64, pricingAt time.Time, actual float64, image, video bool, resolve func(context.Context, int64, int64, float64) float64) *float64 {
	display := actual
	if apiKey != nil && apiKey.GroupID != nil && apiKey.Group != nil && apiKey.Group.HiddenRateMultiplier != nil {
		display = apiKey.Group.RateMultiplier
		if resolve != nil && userID > 0 {
			display = resolve(ctx, userID, *apiKey.GroupID, display)
		}
		switch {
		case video:
			display = resolveVideoRateMultiplier(apiKey, display)
		case image:
			display = resolveImageRateMultiplier(apiKey, display)
		default:
			display *= apiKey.Group.PeakMultiplierAt(pricingAt)
		}
	}
	return &display
}
