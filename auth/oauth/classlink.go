package oauth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDCProvider struct {
	Provider *oidc.Provider
	Verifier *oidc.IDTokenVerifier
	Config   *oauth2.Config
}

// NewClassLinkOIDC discovers OIDC endpoints from the tenant’s sub-domain
// (e.g. “mydistrict”) and returns an OAuth provider compatible with the
// existing middleware.
func NewClassLinkOIDC(ctx context.Context, subdomain, clientID, clientSecret, redirect string) (*Provider, error) {
	iss := fmt.Sprintf("https://%s.classlink.com", subdomain)
	discovery := iss + "/.well-known/openid-configuration"

	op, err := oidc.NewProvider(ctx, discovery)
	if err != nil {
		return nil, err
	}

	// Configure OAuth2
	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     op.Endpoint(),
		RedirectURL:  redirect,
		Scopes:       []string{ScopeOpenID, ScopeEmail, "profile"},
	}

	verifier := op.Verifier(&oidc.Config{ClientID: clientID})

	return &Provider{
		Config: cfg,
		UserInfo: func(c context.Context, tk *oauth2.Token) (map[string]any, error) {
			// Parse & verify ID Token
			rawIDToken, ok := tk.Extra("id_token").(string)
			if !ok {
				return nil, fmt.Errorf("no id_token in token response")
			}
			idTok, err := verifier.Verify(c, rawIDToken)
			if err != nil {
				return nil, err
			}
			var claims map[string]any
			if err := idTok.Claims(&claims); err != nil {
				return nil, err
			}
			return claims, nil
		},
	}, nil
}
