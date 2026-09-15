package audit

import "context"

// Context holds audit trail details for the request lifecycle.
type Context struct {
	IP        string
	UserAgent string
}

type contextKey struct{}

var auditContextKey = contextKey{}

// WithContext returns a new context with the audit context.
func WithContext(ctx context.Context, auditCtx *Context) context.Context {
	return context.WithValue(ctx, auditContextKey, auditCtx)
}

// FromContext retrieves the audit context from the context.
func FromContext(ctx context.Context) *Context {
	if v, ok := ctx.Value(auditContextKey).(*Context); ok {
		return v
	}
	return nil
}
