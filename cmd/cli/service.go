package cli

import (
	"context"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter/provider"
	"github.com/spf13/cobra"
	"net/http"
)

func NewServiceCommandCobra(command *cobra.Command) (*edgecentercdn_go.Service, error) {
	ctx := command.Context()

	apiUrl, ok := apiUrlFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is not set", EC_CDN_API_URL)
	}

	apiKey, ok := apiKeyFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is not set", EC_CDN_API_KEY)
	}

	var opts []provider.ClientOption

	if apiKey != "" {
		opts = append(opts, provider.WithSignerFunc(createSignerFunc(apiKey)))
	}

	client := provider.NewClient(apiUrl, opts...)

	service := edgecentercdn_go.NewService(client)

	return service, nil
}

func apiKeyFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(cobraContextKey(EC_CDN_API_KEY)).(string)
	return value, ok
}

func apiUrlFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(cobraContextKey(EC_CDN_API_URL)).(string)
	return value, ok
}

func createSignerFunc(apiKey string) func(req *http.Request) error {
	return func(req *http.Request) error {
		for k, v := range provider.AuthenticatedHeaders(apiKey) {
			req.Header.Set(k, v)
		}
		return nil
	}
}
