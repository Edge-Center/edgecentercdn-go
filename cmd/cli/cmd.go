package cli

import (
	"context"
	"fmt"
	"os"
)

type cobraContextKey string

const (
	EC_CDN_API_KEY = "EC_CDN_API_KEY"
	EC_CDN_API_URL = "EC_CDN_API_URL"
)

func Run(args []string) error {
	rootCmd.SetArgs(args)

	defer func() {
		if x := recover(); x != nil {
			fmt.Println(x)
			os.Exit(1)
		}
	}()

	ctx := context.Background()

	for _, key := range []string{EC_CDN_API_KEY, EC_CDN_API_URL} {
		if value, exists := os.LookupEnv(key); exists {
			ctx = context.WithValue(ctx, cobraContextKey(key), value)
		}
	}

	return rootCmd.ExecuteContext(ctx)
}
