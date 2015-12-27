// Package httpctx provides a set of interfaces and adapters that add a context parameter to http requests
package httpctx

import (
	"net/http"

	"golang.org/x/net/context"
)

// ContextHandler defines the ServeHTTPWithContext method. Types that implement ContextHandler
// can be registered, via a ContextAdapter, to serve a particular path or subtree in an HTTP server.
type ContextHandler interface {
	ServeHTTPWithContext(context.Context, http.ResponseWriter, *http.Request)
}

// ContextHandlerFunc is an adapter to allow the use of ordinary functions as, context aware, HTTP
// handlers. If f is a function with the appropriate signature, ContextHandlerFunc(f) is a ContextHandler
// tha calls f.
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// ServeHTTPWithContext calls h(ctx, w, req).
func (h ContextHandlerFunc) ServeHTTPWithContext(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	h(ctx, w, req)
}

// ContextAdapter associates a Context and a ContextHandler. Because it implements the http.Handler interface
// ContextAdapter instances can be registered to serve a particular path or subtree in an HTTP server.
type ContextAdapter struct {
	Ctx     context.Context
	Handler ContextHandler
}

// ServeHTTP calls the handler's ServeHTTPWithContext method with the associated Context.
func (ca ContextAdapter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ca.Handler.ServeHTTPWithContext(ca.Ctx, w, req)
}
