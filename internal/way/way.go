/*
 * wars of warp - an implementation of warpwar
 *
 * Copyright (c) 2022 Michael D Henderson
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package way

import (
	"context"
	"net/http"
	"strings"
)

// wayContextKey is the context key type for storing
// parameters in context.Context.
type wayContextKey string

// Router routes HTTP requests.
type Router struct {
	routes []*route
	// NotFound is the http.Handler to call when no routes
	// match. By default uses http.NotFoundHandler().
	NotFound http.Handler
}

// NewRouter makes a new Router.
func NewRouter() *Router {
	return &Router{
		NotFound: http.NotFoundHandler(),
	}
}

func (r *Router) pathSegments(p string) []string {
	return strings.Split(strings.Trim(p, "/"), "/")
}

// Handle adds a handler with the specified method and pattern.
// Method can be any HTTP method string or "*" to match all methods.
// Pattern can contain path segments such as: /item/:id which is
// accessible via the Param function.
// If pattern ends with trailing /, it acts as a prefix.
func (r *Router) Handle(method, pattern string, handler http.Handler) {
	route := &route{
		method:  strings.ToLower(method),
		segs:    r.pathSegments(pattern),
		handler: handler,
		prefix:  strings.HasSuffix(pattern, "/") || strings.HasSuffix(pattern, "..."),
	}
	r.routes = append(r.routes, route)
}

// HandleFunc is the http.HandlerFunc alternative to http.Handle.
func (r *Router) HandleFunc(method, pattern string, fn http.HandlerFunc) {
	r.Handle(method, pattern, fn)
}

// ServeHTTP routes the incoming http.Request based on method and path
// extracting path parameters as it goes.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := strings.ToLower(req.Method)
	segs := r.pathSegments(req.URL.Path)
	for _, route := range r.routes {
		if route.method != method && route.method != "*" {
			continue
		}
		if ctx, ok := route.match(req.Context(), r, segs); ok {
			route.handler.ServeHTTP(w, req.WithContext(ctx))
			return
		}
	}
	r.NotFound.ServeHTTP(w, req)
}

// Param gets the path parameter from the specified Context.
// Returns an empty string if the parameter was not found.
func Param(ctx context.Context, param string) string {
	vStr, ok := ctx.Value(wayContextKey(param)).(string)
	if !ok {
		return ""
	}
	return vStr
}

type route struct {
	method  string
	segs    []string
	handler http.Handler
	prefix  bool
}

func (r *route) match(ctx context.Context, router *Router, segs []string) (context.Context, bool) {
	if len(segs) > len(r.segs) && !r.prefix {
		return nil, false
	}
	for i, seg := range r.segs {
		if i > len(segs)-1 {
			return nil, false
		}
		isParam := false
		if strings.HasPrefix(seg, ":") {
			isParam = true
			seg = strings.TrimPrefix(seg, ":")
		}
		if !isParam { // verbatim check
			if strings.HasSuffix(seg, "...") {
				if strings.HasPrefix(segs[i], seg[:len(seg)-3]) {
					return ctx, true
				}
			}
			if seg != segs[i] {
				return nil, false
			}
		}
		if isParam {
			ctx = context.WithValue(ctx, wayContextKey(seg), segs[i])
		}
	}
	return ctx, true
}
