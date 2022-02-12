package google

import "forum/pkg/oauth2"

// Endpoint is Google's OAuth 2.0 default endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://accounts.google.com/o/oauth2/auth",
	TokenURL:  "https://oauth2.googleapis.com/token",
	AuthStyle: oauth2.AuthStyleInParams,
}
