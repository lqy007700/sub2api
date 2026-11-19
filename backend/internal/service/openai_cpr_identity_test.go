package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestCPRUpstreamRebuildsIdentityForEachTarget(t *testing.T) {
	ctx := WithAuthenticatedCPRIdentity(context.Background(), 42, 101, "团队/A +")
	account := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"base_url": "https://cpr.example/v1",
			"api_key":  "first-cpr",
		},
		Extra: map[string]any{AccountExtraCPRTrustedIdentity: true, accountExtraCPRSendTokenName: true},
	}
	upstream := &httpUpstreamRecorder{}
	gateway := &OpenAIGatewayService{httpUpstream: upstream}
	req := httptest.NewRequest(http.MethodPost, "https://cpr.example/v1/responses", nil).WithContext(ctx)
	req.Header["x-CpR-ToKeN-Id"] = []string{"555", "556"}
	req.Header["authorization"] = []string{"Bearer terminal-key"}
	req.Header.Set("X-API-Key", "terminal-key")
	_, err := gateway.doOpenAIUpstream(req, "", account)
	require.NoError(t, err)
	require.Equal(t, []string{"101"}, upstream.lastReq.Header.Values(cprTokenIDHeader))
	require.NotContains(t, upstream.lastReq.Header, "x-CpR-ToKeN-Id")
	require.NotContains(t, upstream.lastReq.Header, "authorization")
	require.Empty(t, upstream.lastReq.Header.Get("X-API-Key"))
	require.Equal(t, "42", upstream.lastReq.Header.Get(cprUserIDHeader))
	require.Equal(t, "%E5%9B%A2%E9%98%9F%2FA%20+", upstream.lastReq.Header.Get(cprTokenNameHeader))
	require.Equal(t, "Bearer first-cpr", upstream.lastReq.Header.Get("Authorization"))
	require.True(t, HTTPUpstreamRedirectsDisabled(upstream.lastReq.Context()))

	// Reuse the headers from the first attempt to verify that stale values are discarded.
	account.Credentials["api_key"] = "second-cpr"
	delete(account.Extra, accountExtraCPRSendTokenName)
	_, err = gateway.doOpenAIUpstream(req, "", account)
	require.NoError(t, err)
	require.Equal(t, "Bearer second-cpr", upstream.lastReq.Header.Get("Authorization"))
	require.Empty(t, upstream.lastReq.Header.Get(cprTokenNameHeader))

	req.Header.Set("Authorization", "Bearer regular-upstream")
	_, err = gateway.doOpenAIUpstream(req, "", &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKey})
	require.NoError(t, err)
	require.Equal(t, "Bearer regular-upstream", upstream.lastReq.Header.Get("Authorization"))
	require.Empty(t, upstream.lastReq.Header.Get(cprUserIDHeader))
	require.Empty(t, upstream.lastReq.Header.Get(cprTokenIDHeader))
	require.Empty(t, upstream.lastReq.Header.Get(cprTokenNameHeader))
	require.False(t, HTTPUpstreamRedirectsDisabled(upstream.lastReq.Context()))
}

func TestCPRAccountProbeHasOnlyUpstreamAuthentication(t *testing.T) {
	account := &Account{
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Credentials: map[string]any{"base_url": "https://cpr.example/v1", "api_key": "shared-cpr"},
		Extra:       map[string]any{AccountExtraCPRTrustedIdentity: true, accountExtraCPRSendTokenName: true},
	}
	upstream := &httpUpstreamRecorder{}
	gateway := &OpenAIGatewayService{httpUpstream: upstream}
	req := httptest.NewRequest(http.MethodPost, "https://cpr.example/v1/responses", nil)
	_, err := gateway.doOpenAIUpstream(req, "", account)
	require.ErrorContains(t, err, "requires an authenticated user and API key")
	require.Empty(t, upstream.requests)

	ctx := WithAuthenticatedCPRIdentity(req.Context(), 42, 101, "ignored during probe")
	req = req.WithContext(ctx)
	req.Header.Set(cprTokenIDHeader, "forged")
	testService := &AccountTestService{httpUpstream: upstream}
	_, err = testService.doOpenAIAccountTestUpstream(req, "", account, false)
	require.NoError(t, err)
	require.Equal(t, http.Header{"Authorization": {"Bearer shared-cpr"}}, upstream.lastReq.Header)
	require.True(t, HTTPUpstreamRedirectsDisabled(upstream.lastReq.Context()))

	wrongTarget := httptest.NewRequest(http.MethodPost, "https://other.example/v1/responses", nil).WithContext(ctx)
	_, err = gateway.doOpenAIUpstream(wrongTarget, "", account)
	require.ErrorContains(t, err, "does not target the configured CPR origin")
	require.Len(t, upstream.requests, 1)
}

func TestCPRAccountRequiresExplicitBaseURL(t *testing.T) {
	account := &Account{
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Credentials: map[string]any{"api_key": "shared-cpr"},
		Extra:       map[string]any{AccountExtraCPRTrustedIdentity: true},
	}
	require.ErrorContains(t, validateCPRIdentityConfig(account), "custom CPR Base URL")
	account.Credentials["base_url"] = "https://api.openai.com/v1"
	require.ErrorContains(t, validateCPRIdentityConfig(account), "custom CPR Base URL")
	account.Credentials["base_url"] = "https://cpr.example/v1"
	require.NoError(t, validateCPRIdentityConfig(account))
}

func TestCPRAccountAllowsOnlyHTTPUpstreamTransport(t *testing.T) {
	cfg := &config.Config{}
	cfg.Gateway.OpenAIWS.Enabled = true
	cfg.Gateway.OpenAIWS.APIKeyEnabled = true
	cfg.Gateway.OpenAIWS.ResponsesWebsocketsV2 = true
	cfg.Gateway.OpenAIWS.ModeRouterV2Enabled = true
	resolver := NewOpenAIWSProtocolResolver(cfg)
	gateway := &OpenAIGatewayService{cfg: cfg, openaiWSResolver: resolver}
	account := &Account{
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Concurrency: 1,
		Extra: map[string]any{
			AccountExtraCPRTrustedIdentity:               true,
			"openai_apikey_responses_websockets_v2_mode": OpenAIWSIngressModeCtxPool,
		},
	}
	require.Equal(t, OpenAIUpstreamTransportHTTPSSE, resolver.Resolve(account).Transport)
	require.False(t, gateway.isOpenAIAccountTransportCompatible(account, OpenAIUpstreamTransportResponsesWebsocketV2Ingress))
	account.Extra["openai_apikey_responses_websockets_v2_mode"] = OpenAIWSIngressModeHTTPBridge
	require.True(t, gateway.isOpenAIAccountTransportCompatible(account, OpenAIUpstreamTransportResponsesWebsocketV2Ingress))
}
