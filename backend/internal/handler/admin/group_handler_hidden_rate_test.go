package admin

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateGroupRequestHiddenRateMultiplierTriState(t *testing.T) {
	for _, tt := range []struct {
		name, body string
		set        bool
		value      *float64
	}{
		{name: "omitted", body: `{}`},
		{name: "null", body: `{"hidden_rate_multiplier":null}`, set: true},
		{name: "number", body: `{"hidden_rate_multiplier":2.75}`, set: true, value: func() *float64 { v := 2.75; return &v }()},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var req UpdateGroupRequest
			require.NoError(t, json.Unmarshal([]byte(tt.body), &req))
			require.Equal(t, tt.set, req.HiddenRateMultiplier.set)
			require.Equal(t, tt.value, req.HiddenRateMultiplier.value)
		})
	}
	for _, body := range []string{
		`{"hidden_rate_multiplier":"2.75"}`, `{"hidden_rate_multiplier":true}`, `{"hidden_rate_multiplier":{}}`,
		`{"hidden_rate_multiplier":NaN}`, `{"hidden_rate_multiplier":Infinity}`, `{"hidden_rate_multiplier":-Infinity}`,
		`{"hidden_rate_multiplier":1e309}`, `{"hidden_rate_multiplier":-1e309}`, `{"hidden_rate_multiplier":"NaN"}`,
	} {
		var req UpdateGroupRequest
		require.Error(t, json.Unmarshal([]byte(body), &req))
		var createReq CreateGroupRequest
		require.Error(t, json.Unmarshal([]byte(body), &createReq))
	}
}

func TestCreateGroupRequestHiddenRateMultiplier(t *testing.T) {
	var req CreateGroupRequest
	require.NoError(t, json.Unmarshal([]byte(`{"hidden_rate_multiplier":2.75}`), &req))
	require.NotNil(t, req.HiddenRateMultiplier)
	require.Equal(t, 2.75, *req.HiddenRateMultiplier)
}
