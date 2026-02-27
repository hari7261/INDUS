package cli

import "context"

type contextKey string

const (
	configKey contextKey = "config"
)

func WithConfig(ctx context.Context, cfg interface{}) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

func ConfigFromContext(ctx context.Context) (interface{}, bool) {
	val := ctx.Value(configKey)
	if val == nil {
		return nil, false
	}
	return val, true
}
