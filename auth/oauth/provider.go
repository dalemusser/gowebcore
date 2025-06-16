package oauth

import (
	"context"

	"golang.org/x/oauth2"
)

// Provider wraps an OAuth2 config and a function to fetch user info.
type Provider struct {
	Config   *oauth2.Config
	UserInfo func(ctx context.Context, token *oauth2.Token) (map[string]any, error)
}

// AuthURL returns the redirect URL for the caller to visit.
func (p *Provider) AuthURL(state string) string {
	return p.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Exchange turns ?code= into a token.
func (p *Provider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.Config.Exchange(ctx, code)
}

// Default scopes
const (
	ScopeEmail  = "email"
	ScopeOpenID = "openid"
)
