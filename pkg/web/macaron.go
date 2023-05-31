//go:build go1.3
// +build go1.3

// Copyright 2014 The Macaron Authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package macaron is a high productive and modular web framework in Go.
package web

import (
	"context"
	"fmt"
	_ "unsafe"

	"net/http"
	"strings"

	"github.com/beakeyz/dadjoke-gen/pkg/api/custom_context"
)

const _VERSION = "1.3.4.0805"

const (
	DEV  = "development"
	PROD = "production"
)

var (
	// Env is the environment that Macaron is executing in.
	// The MACARON_ENV is read on initialization to set this variable.
	Env = DEV
)

type (
	handlerStd       = func(http.ResponseWriter, *http.Request)
	handlerStdCtx    = func(http.ResponseWriter, *http.Request, *Context)
	handlerStdReqCtx = func(http.ResponseWriter, *http.Request, *ReqContext)
	handlerReqCtx    = func(*ReqContext)
	handlerReqCtxRes = func(*ReqContext) Response
	handlerCtx       = func(*Context)
)

func wrap_handler(h Handler) http.HandlerFunc {
	switch handle := h.(type) {
	case handlerStd:
		return handle
	case handlerStdCtx:
		return func(w http.ResponseWriter, r *http.Request) {
			handle(w, r, FromContext(r.Context()))
		}
	case handlerStdReqCtx:
		return func(w http.ResponseWriter, r *http.Request) {
			handle(w, r, getReqCtx(r.Context()))
		}
	case handlerReqCtx:
		return func(w http.ResponseWriter, r *http.Request) {
			handle(getReqCtx(r.Context()))
		}
	case handlerReqCtxRes:
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := getReqCtx(r.Context())
			res := handle(ctx)
			if res != nil {
				res.WriteTo(ctx)
			}
		}
	case handlerCtx:
		return func(w http.ResponseWriter, r *http.Request) {
			handle(FromContext(r.Context()))
		}
	}

	panic(fmt.Sprintf("unexpected handler type: %T", h))
}

func getReqCtx(ctx context.Context) *ReqContext {
	reqCtx, ok := ctx.Value(custom_context.Key{}).(*ReqContext)
	if !ok {
		panic("no *ReqContext found")
	}
	return reqCtx
}

func Version() string {
	return _VERSION
}

// Handler can be any callable function.
// Macaron attempts to inject services into the handler's argument list,
// and panics if an argument could not be fulfilled via dependency injection.
type Handler interface{}

// validateAndWrapHandler makes sure a handler is a callable function, it panics if not.
// When the handler is also potential to be any built-in inject.FastInvoker,
// it wraps the handler automatically to have some performance gain.
func validateAndWrapHandler(h Handler) http.Handler {
	return wrap_handler(h)
}

// validateAndWrapHandlers preforms validation and wrapping for each input handler.
// It accepts an optional wrapper function to perform custom wrapping on handlers.
func validateAndWrapHandlers(handlers []Handler) []http.Handler {
	wrappedHandlers := make([]http.Handler, len(handlers))
	for i, h := range handlers {
		wrappedHandlers[i] = validateAndWrapHandler(h)
	}

	return wrappedHandlers
}

// Macaron represents the top level web application.
// Injector methods can be invoked to map services on a global level.
type Macaron struct {
	handlers []http.Handler

	urlPrefix string // For suburl support.
	*Router
}

// New creates a bare bones Macaron instance.
// Use this method if you want to have full control over the middleware that is used.
func New() *Macaron {
	m := &Macaron{Router: NewRouter()}
	m.Router.m = m
	m.NotFound(http.NotFound)
	return m
}

// BeforeHandler represents a handler executes at beginning of every request.
// Macaron stops future process when it returns true.
type BeforeHandler func(rw http.ResponseWriter, req *http.Request) bool

// macaronContextKey is used to store/fetch web.Context inside context.Context
type macaronContextKey struct{}

// FromContext returns the macaron context stored in a context.Context, if any.
func FromContext(c context.Context) *Context {
	if mc, ok := c.Value(macaronContextKey{}).(*Context); ok {
		return mc
	}
	return nil
}

type paramsKey struct{}

// Params returns the named route parameters for the current request, if any.
func Params(r *http.Request) map[string]string {
	if rv := r.Context().Value(paramsKey{}); rv != nil {
		return rv.(map[string]string)
	}
	return map[string]string{}
}

// SetURLParams sets the named URL parameters for the given request. This should only be used for testing purposes.
func SetURLParams(r *http.Request, vars map[string]string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), paramsKey{}, vars))
}

// UseMiddleware is a traditional approach to writing middleware in Go.
// A middleware is a function that has a reference to the next handler in the chain
// and returns the actual middleware handler, that may do its job and optionally
// call next.
// Due to how Macaron handles/injects requests and responses we patch the web.Context
// to use the new ResponseWriter and http.Request here. The caller may only call
// `next.ServeHTTP(rw, req)` to pass a modified response writer and/or a request to the
// further middlewares in the chain.
func (m *Macaron) UseMiddleware(middleware func(http.Handler) http.Handler) {
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		c := FromContext(req.Context())
		c.Req = req
		if mrw, ok := rw.(*responseWriter); ok {
			c.Resp = mrw
		} else {
			c.Resp = NewResponseWriter(req.Method, rw)
		}
		c.Next()
	})
	m.handlers = append(m.handlers, middleware(next))
}

// Use adds a middleware Handler to the stack,
// and panics if the handler is not a callable func.
// Middleware Handlers are invoked in the order that they are added.
func (m *Macaron) Use(handler Handler) {
	h := validateAndWrapHandler(handler)
	m.handlers = append(m.handlers, h)
}

func (m *Macaron) createContext(rw http.ResponseWriter, req *http.Request) *Context {
	c := &Context{
		handlers: m.handlers,
		index:    0,
		Router:   m.Router,
		Resp:     NewResponseWriter(req.Method, rw),
	}

	c.Req = req.WithContext(context.WithValue(req.Context(), macaronContextKey{}, c))
	return c
}

// ServeHTTP is the HTTP Entry point for a Macaron instance.
// Useful if you want to control your own HTTP server.
// Be aware that none of middleware will run without registering any router.
func (m *Macaron) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.URL.Path = strings.TrimPrefix(req.URL.Path, m.urlPrefix)
	m.Router.ServeHTTP(rw, req)
}

// SetURLPrefix sets URL prefix of router layer, so that it support suburl.
func (m *Macaron) SetURLPrefix(prefix string) {
	m.urlPrefix = prefix
}
