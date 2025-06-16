package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
	ghcfg "golang.org/x/oauth2/github"
)

// NewGitHub returns an OAuth provider for github.com.
//
// redirect must match the “Authorization callback URL” you configure in
// GitHub → Settings → Developer settings → OAuth Apps.
func NewGitHub(clientID, clientSecret, redirect string) *Provider {
	return &Provider{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     ghcfg.Endpoint,
			RedirectURL:  redirect,
			Scopes:       []string{"read:user", "user:email"},
		},
		UserInfo: githubUserInfo,
	}
}

// ---------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------

func githubUserInfo(ctx context.Context, token *oauth2.Token) (map[string]any, error) {
	cl := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	// 1) primary profile
	profileResp, err := cl.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer profileResp.Body.Close()
	if profileResp.StatusCode != http.StatusOK {
		return nil, errors.New("github userinfo: profile request failed")
	}
	var user map[string]any
	if err := json.NewDecoder(profileResp.Body).Decode(&user); err != nil {
		return nil, err
	}

	// 2) primary e-mail
	emailResp, err := cl.Get("https://api.github.com/user/emails")
	if err != nil {
		return nil, err
	}
	defer emailResp.Body.Close()
	if emailResp.StatusCode == http.StatusOK {
		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}
		if err := json.NewDecoder(emailResp.Body).Decode(&emails); err == nil {
			for _, e := range emails {
				if e.Primary && e.Verified {
					user["email"] = e.Email
					break
				}
			}
		}
	}
	return user, nil
}
