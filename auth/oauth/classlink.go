// auth/oauth/classlink.go
package oauth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// NewClassLinkOIDC discovers the OIDC metadata for a ClassLink tenant
// (sub-domain like “mydistrict”) and returns a gowebcore *Provider
// that works with middleware.AuthRoutes / RequireAuth.
func NewClassLinkOIDC(
	ctx context.Context,
	subdomain, clientID, clientSecret, redirect string,
) (*Provider, error) {

	issuer := fmt.Sprintf("https://%s.classlink.com", subdomain)
	wellKnown := issuer + "/.well-known/openid-configuration"

	op, err := oidc.NewProvider(ctx, wellKnown)
	if err != nil {
		return nil, err
	}

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
		UserInfo: func(c context.Context, tok *oauth2.Token) (map[string]any, error) {
			rawID, ok := tok.Extra("id_token").(string)
			if !ok {
				return nil, fmt.Errorf("no id_token in token response")
			}
			idTok, err := verifier.Verify(c, rawID)
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
