package postgres

import "context"

type ctxKeyQuerier struct{}

func ContextWithQuerier(ctx context.Context, q Querier) context.Context {
	ctx = context.WithValue(ctx, ctxKeyQuerier{}, q)
	return ctx
}

func QuerierFromContext(ctx context.Context, fallback Querier) Querier {
	if q, ok := ctx.Value(ctxKeyQuerier{}).(Querier); ok {
		return q
	}
	return fallback
}
