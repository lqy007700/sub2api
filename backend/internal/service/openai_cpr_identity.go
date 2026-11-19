package service

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/util/urlvalidator"
)

const (
	AccountExtraCPRTrustedIdentity = "cpr_trusted_identity"
	accountExtraCPRSendTokenName   = "cpr_send_token_name"
	cprUserIDHeader                = "X-CPR-User-Id"
	cprTokenIDHeader               = "X-CPR-Token-Id"
	cprTokenNameHeader             = "X-CPR-Token-Name"
)

type authenticatedCPRIdentity struct {
	userID    int64
	tokenID   int64
	tokenName string
}

type cprIdentityContextKey struct{}

// WithAuthenticatedCPRIdentity records gateway API key identity after successful authentication.
func WithAuthenticatedCPRIdentity(ctx context.Context, userID, tokenID int64, tokenName string) context.Context {
	return context.WithValue(ctx, cprIdentityContextKey{}, authenticatedCPRIdentity{userID, tokenID, tokenName})
}

func (a *Account) usesCPRTrustedIdentity() bool {
	return a != nil && a.IsOpenAIApiKey() && a.Extra[AccountExtraCPRTrustedIdentity] == true
}

func cprBaseURL(account *Account) (*url.URL, error) {
	// The request builders also enforce the installation's outbound URL policy.
	base, err := urlvalidator.ValidateURLFormat(account.GetCredential("base_url"), true)
	if err == nil {
		parsed, parseErr := url.Parse(base)
		if parseErr == nil && parsed.User == nil && !parsed.ForceQuery && parsed.RawQuery == "" &&
			parsed.Fragment == "" && !strings.EqualFold(strings.TrimRight(parsed.Hostname(), "."), "api.openai.com") {
			return parsed, nil
		}
	}
	return nil, infraerrors.BadRequest("INVALID_CPR_BASE_URL", "CPR identity forwarding requires a custom CPR Base URL without user info, query or fragment")
}

func validateCPRIdentityConfig(account *Account) error {
	if account.Extra[AccountExtraCPRTrustedIdentity] != true {
		return nil
	}
	if !account.IsOpenAIApiKey() {
		return infraerrors.BadRequest("INVALID_CPR_ACCOUNT_TYPE", "CPR identity forwarding requires an OpenAI API key account")
	}
	if strings.TrimSpace(account.GetOpenAIApiKey()) == "" {
		return infraerrors.BadRequest("CPR_API_KEY_REQUIRED", "CPR identity forwarding requires a CPR API key")
	}
	_, err := cprBaseURL(account)
	return err
}

func isCPRIdentityHeader(name string) bool {
	return strings.EqualFold(name, cprUserIDHeader) ||
		strings.EqualFold(name, cprTokenIDHeader) ||
		strings.EqualFold(name, cprTokenNameHeader)
}

func prepareCPRUpstreamRequest(req *http.Request, account *Account, forwardIdentity bool) (*http.Request, error) {
	for name := range req.Header {
		if isCPRIdentityHeader(name) {
			delete(req.Header, name)
		}
	}
	if !account.usesCPRTrustedIdentity() {
		return req, nil
	}
	base, err := cprBaseURL(account)
	if err != nil {
		return nil, err
	}
	if req.URL == nil || !strings.EqualFold(req.URL.Scheme, base.Scheme) ||
		!strings.EqualFold(req.URL.Host, base.Host) ||
		(req.Host != "" && !strings.EqualFold(req.Host, base.Host)) {
		return nil, infraerrors.BadRequest("CPR_TARGET_MISMATCH", "Upstream request does not target the configured CPR origin")
	}
	key := strings.TrimSpace(account.GetOpenAIApiKey())
	if key == "" {
		return nil, infraerrors.BadRequest("CPR_API_KEY_REQUIRED", "CPR identity forwarding requires a CPR API key")
	}
	req = req.Clone(WithHTTPUpstreamRedirectsDisabled(req.Context()))
	if forwardIdentity {
		identity, ok := req.Context().Value(cprIdentityContextKey{}).(authenticatedCPRIdentity)
		if !ok || identity.userID <= 0 || identity.tokenID <= 0 {
			return nil, infraerrors.BadRequest("CPR_IDENTITY_MISSING", "CPR identity forwarding requires an authenticated user and API key")
		}
		req.Header.Set(cprUserIDHeader, strconv.FormatInt(identity.userID, 10))
		req.Header.Set(cprTokenIDHeader, strconv.FormatInt(identity.tokenID, 10))
		if account.Extra[accountExtraCPRSendTokenName] == true && strings.TrimSpace(identity.tokenName) != "" {
			req.Header.Set(cprTokenNameHeader, url.PathEscape(identity.tokenName))
		}
	}
	for name := range req.Header {
		if strings.EqualFold(name, "Authorization") || strings.EqualFold(name, "X-API-Key") ||
			strings.EqualFold(name, "X-Goog-API-Key") {
			delete(req.Header, name)
		}
	}
	req.Header.Set("Authorization", "Bearer "+key)
	return req, nil
}
