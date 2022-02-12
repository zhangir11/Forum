package oauth2

import (
	"bytes"
	"context"
	"forum/pkg/internal"
	"net/url"
	"strings"
)

// AuthCodeURL returns a URL to OAuth 2.0 provider's consent page
// that asks for permissions for the required scopes explicitly.
//
// State is a token to protect the user from CSRF attacks. You must
// always provide a non-empty string and validate that it matches the
// the state query parameter on your redirect callback.
// See http://tools.ietf.org/html/rfc6749#section-10.12 for more info.
func (c *Config) AuthCodeURL(state string) string {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.AuthURL)
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {c.ClientID},
	}
	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}
	if len(c.Scopes) > 0 {
		v.Set("scope", strings.Join(c.Scopes, " "))
	}
	if state != "" {
		// TODO(light): Docs say never to omit state; don't allow empty.
		v.Set("state", state)
	}

	if strings.Contains(c.Endpoint.AuthURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	return buf.String()
}

// Exchange converts an authorization code into a token.
//
// It is used after a resource provider redirects the user back
// to the Redirect URI (the URL obtained from AuthCodeURL).
//
// The provided context optionally controls which HTTP client is used. See the HTTPClient variable.
//
// The code will be in the *http.Request.FormValue("code"). Before
// calling Exchange, be sure to validate FormValue("state").
//
// Opts may include the PKCE verifier code if previously used in AuthCodeURL.
// See https://www.oauth.com/oauth2-servers/pkce/ for more info.
func (c *Config) Exchange(ctx context.Context, code string) (*Token, error) {
	v := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
	}

	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}

	return retrieveToken(ctx, c, v)
}

// retrieveToken takes a *Config and uses that to retrieve an *internal.Token.
// This token is then mapped from *internal.Token into an *oauth2.Token which is returned along
// with an error..
func retrieveToken(ctx context.Context, c *Config, v url.Values) (*Token, error) {
	tk, err := internal.RetrieveToken(ctx, c.ClientID, c.ClientSecret, c.Endpoint.TokenURL, v, internal.AuthStyle(c.Endpoint.AuthStyle))
	if err != nil {
		if rErr, ok := err.(*internal.RetrieveError); ok {
			return nil, (*internal.RetrieveError)(rErr)
		}
		return nil, err
	}
	return tokenFromInternal(tk), nil
}

// tokenFromInternal maps an *internal.Token struct into
// a *Token struct.
func tokenFromInternal(t *internal.Token) *Token {
	if t == nil {
		return nil
	}
	return &Token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
		raw:          t.Raw,
	}
}
