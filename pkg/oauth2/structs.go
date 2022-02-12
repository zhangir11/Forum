package oauth2

import "time"

type Config struct {
	// ClientID is the application's ID.
	ClientID string

	// ClientSecret is the application's secret.
	ClientSecret string

	// Endpoint contains the resource server's token endpoint
	// URLs. These are constants specific to each server and are
	// often available via site-specific packages, such as
	// google.Endpoint or github.Endpoint.
	Endpoint Endpoint

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string

	// Scope specifies optional requested permissions.
	Scopes []string
}

// Endpoint represents an OAuth 2.0 provider's authorization and token
// endpoint URLs.
type Endpoint struct {
	AuthURL  string
	TokenURL string

	// AuthStyle optionally specifies how the endpoint wants the
	// client ID & client secret sent. The zero value means to
	// auto-detect.
	AuthStyle AuthStyle
}

// AuthStyle represents how requests for tokens are authenticated
// to the server.
type AuthStyle int

// Token represents the credentials used to authorize
// the requests to access protected resources on the OAuth 2.0
// provider's backend.
//
// Most users of this package should not access fields of Token
// directly. They're exported mostly for use by related packages
// implementing derivative OAuth2 flows.
type Token struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string `json:"access_token"`

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string `json:"token_type,omitempty"`

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string `json:"refresh_token,omitempty"`

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time `json:"expiry,omitempty"`

	// raw optionally contains extra metadata from the server
	// when updating a token.
	raw interface{}
}
