package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

// Clever endpoints
var cleverEndpoint = oauth2.Endpoint{
	AuthURL:  "https://clever.com/oauth/authorize",
	TokenURL: "https://clever.com/oauth/tokens",
}

// NewClever wires an OAuth2 provider for clever.com.
//
//	redirect must exactly match the “Redirect URI” set in the Clever dashboard.
func NewClever(clientID, clientSecret, redirect string) *Provider {
	return &Provider{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     cleverEndpoint,
			RedirectURL:  redirect,
			Scopes:       []string{"read:me"},
		},
		UserInfo: cleverUserInfo,
	}
}

// ------------------------------------------------------------------
// Private helpers
// ------------------------------------------------------------------

func cleverUserInfo(ctx context.Context, token *oauth2.Token) (map[string]any, error) {
	cl := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := cl.Get("https://api.clever.com/v3.0/me")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("clever userinfo: non-200 response")
	}
	var wrapper struct {
		Data struct {
			ID         string         `json:"id"`
			Type       string         `json:"type"` // district_admin, teacher, student, …
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	// Flatten into a generic map for the session.
	info := map[string]any{
		"id":   wrapper.Data.ID,
		"type": wrapper.Data.Type,
	}
	for k, v := range wrapper.Data.Attributes {
		info[k] = v
	}
	return info, nil
}
