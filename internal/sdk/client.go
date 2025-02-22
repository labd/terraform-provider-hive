package sdk

import (
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

// HiveRoundTripper wraps an underlying RoundTripper.
type HiveRoundTripper struct {
	rt        http.RoundTripper
	authToken string
	userAgent string
}

// RoundTrip adds the Authorization and User-Agent headers to every request.
func (hrt *HiveRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+hrt.authToken)
	req.Header.Set("User-Agent", hrt.userAgent)
	return hrt.rt.RoundTrip(req)
}

// HiveClient encapsulates an HTTP client, a GraphQL endpoint, and an API token.
type HiveClient struct {
	client *graphql.Client
}

// NewHiveClient creates a new HiveClient instance.
func NewHiveClient(client *http.Client, endpoint, token string) *HiveClient {
	client.Transport = &HiveRoundTripper{
		rt:        client.Transport,
		authToken: token,
		userAgent: "terraform-provider-hive/0.0.1",
	}

	gqlClient := graphql.NewClient(endpoint, client)

	return &HiveClient{
		client: &gqlClient,
	}
}
