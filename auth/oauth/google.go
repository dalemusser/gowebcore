package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
	googleCfg "golang.org/x/oauth2/google"
)

func NewGoogle(clientID, clientSecret, redirect string) *Provider {
	return &Provider{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     googleCfg.Endpoint,
			RedirectURL:  redirect,
			Scopes:       []string{ScopeEmail, "profile"},
		},
		UserInfo: googleUserInfo,
	}
}

func googleUserInfo(ctx context.Context, token *oauth2.Token) (map[string]any, error) {
	cl := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := cl.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("google userinfo: non-200")
	}
	var data map[string]any
	return data, json.NewDecoder(resp.Body).Decode(&data)
}
