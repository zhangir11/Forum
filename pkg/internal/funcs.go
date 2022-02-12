package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"forum/pkg/ctxhttp"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

func RetrieveToken(ctx context.Context, clientID, clientSecret, tokenURL string, v url.Values, authStyle AuthStyle) (*Token, error) {
	needsAuthStyleProbe := authStyle == 0
	if needsAuthStyleProbe {
		if style, ok := lookupAuthStyle(tokenURL); ok {
			authStyle = style
			needsAuthStyleProbe = false
		} else {
			authStyle = AuthStyleInHeader // the first way we'll try
		}
	}
	req, err := newTokenRequest(tokenURL, clientID, clientSecret, v, authStyle)
	if err != nil {
		return nil, err
	}
	token, err := doTokenRoundTrip(ctx, req)
	if err != nil && needsAuthStyleProbe {
		// If we get an error, assume the server wants the
		// clientID & clientSecret in a different form.
		// See https://code.google.com/p/goauth2/issues/detail?id=31 for background.
		// In summary:
		// - Reddit only accepts client secret in the Authorization header
		// - Dropbox accepts either it in URL param or Auth header, but not both.
		// - Google only accepts URL param (not spec compliant?), not Auth header
		// - Stripe only accepts client secret in Auth header with Bearer method, not Basic
		//
		// We used to maintain a big table in this code of all the sites and which way
		// they went, but maintaining it didn't scale & got annoying.
		// So just try both ways.
		authStyle = AuthStyleInParams // the second way we'll try
		req, _ = newTokenRequest(tokenURL, clientID, clientSecret, v, authStyle)
		token, err = doTokenRoundTrip(ctx, req)
	}
	if needsAuthStyleProbe && err == nil {
		setAuthStyle(tokenURL, authStyle)
	}
	// Don't overwrite `RefreshToken` with an empty value
	// if this was a token refreshing request.
	if token != nil && token.RefreshToken == "" {
		token.RefreshToken = v.Get("refresh_token")
	}
	return token, err
}

// lookupAuthStyle reports which auth style we last used with tokenURL
// when calling RetrieveToken and whether we have ever done so.
func lookupAuthStyle(tokenURL string) (style AuthStyle, ok bool) {
	authStyleCache.Lock()
	defer authStyleCache.Unlock()
	style, ok = authStyleCache.m[tokenURL]
	return
}

// authStyleCache is the set of tokenURLs we've successfully used via
// RetrieveToken and which style auth we ended up using.
// It's called a cache, but it doesn't (yet?) shrink. It's expected that
// the set of OAuth2 servers a program contacts over time is fixed and
// small.
var authStyleCache struct {
	sync.Mutex
	m map[string]AuthStyle // keyed by tokenURL
}

// newTokenRequest returns a new *http.Request to retrieve a new token
// from tokenURL using the provided clientID, clientSecret, and POST
// body parameters.
//
// inParams is whether the clientID & clientSecret should be encoded
// as the POST body. An 'inParams' value of true means to send it in
// the POST body (along with any values in v); false means to send it
// in the Authorization header.
func newTokenRequest(tokenURL, clientID, clientSecret string, v url.Values, authStyle AuthStyle) (*http.Request, error) {
	if authStyle == AuthStyleInParams {
		v = cloneURLValues(v)
		if clientID != "" {
			v.Set("client_id", clientID)
		}
		if clientSecret != "" {
			v.Set("client_secret", clientSecret)
		}
	}
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if authStyle == AuthStyleInHeader {
		req.SetBasicAuth(url.QueryEscape(clientID), url.QueryEscape(clientSecret))
	}
	return req, nil
}

func cloneURLValues(v url.Values) url.Values {
	v2 := make(url.Values, len(v))
	for k, vv := range v {
		v2[k] = append([]string(nil), vv...)
	}
	return v2
}

func doTokenRoundTrip(ctx context.Context, req *http.Request) (*Token, error) {
	r, err := ctxhttp.Do(ctx, ContextClient(ctx), req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	r.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}
	if code := r.StatusCode; code < 200 || code > 299 {
		return nil, &RetrieveError{
			Response: r,
			Body:     body,
		}
	}

	var token *Token
	content, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch content {
	case "application/x-www-form-urlencoded", "text/plain":
		vals, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		token = &Token{
			AccessToken:  vals.Get("access_token"),
			TokenType:    vals.Get("token_type"),
			RefreshToken: vals.Get("refresh_token"),
			Raw:          vals,
		}
		e := vals.Get("expires_in")
		expires, _ := strconv.Atoi(e)
		if expires != 0 {
			token.Expiry = time.Now().Add(time.Duration(expires) * time.Second)
		}
	default:
		var tj tokenJSON
		if err = json.Unmarshal(body, &tj); err != nil {
			return nil, err
		}
		token = &Token{
			AccessToken:  tj.AccessToken,
			TokenType:    tj.TokenType,
			RefreshToken: tj.RefreshToken,
			Expiry:       tj.expiry(),
			Raw:          make(map[string]interface{}),
		}
		json.Unmarshal(body, &token.Raw) // no error checks for optional fields
	}
	if token.AccessToken == "" {
		return nil, errors.New("oauth2: server response missing access_token")
	}
	return token, nil
}

// setAuthStyle adds an entry to authStyleCache, documented above.
func setAuthStyle(tokenURL string, v AuthStyle) {
	authStyleCache.Lock()
	defer authStyleCache.Unlock()
	if authStyleCache.m == nil {
		authStyleCache.m = make(map[string]AuthStyle)
	}
	authStyleCache.m[tokenURL] = v
}

func ContextClient(ctx context.Context) *http.Client {
	if ctx != nil {
		if hc, ok := ctx.Value(HTTPClient).(*http.Client); ok {
			return hc
		}
	}
	if appengineClientHook != nil {
		return appengineClientHook(ctx)
	}
	return http.DefaultClient
}

func (e *tokenJSON) expiry() (t time.Time) {
	if v := e.ExpiresIn; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}
	return
}

var appengineClientHook func(context.Context) *http.Client
