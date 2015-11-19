// Package main provides a simple UBER hypermedia driven project/issue tracking service.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	ctx     context.Context
	handler ContextHandler
}

// ServeHTTP calls the handler's ServeHTTPWithContext method with the associated Context.
func (ca ContextAdapter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ca.handler.ServeHTTPWithContext(ca.ctx, w, req)
}

// udata represents the individual data elements of an Uber hypermedia document.
type udata struct {
	ID         string   `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Rel        []string `json:"rel,omitempty"`
	Label      string   `json:"label,omitempty"`
	URL        string   `json:"url,omitempty"`
	Template   bool     `json:"template,omitempty"`
	Action     string   `json:"action,omitempty"`
	Transclude bool     `json:"transclude,omitempty"`
	Model      string   `json:"model,omitempty"`
	Sending    string   `json:"sending,omitempty"`
	Accepting  []string `json:"accepting,omitempty"`
	Value      string   `json:"value,omitempty"`
	Data       []udata  `json:"data,omitempty"`
}

// ubody is the body of an Uber hypermedia document.
type ubody struct {
	Version string  `json:"version"`
	Data    []udata `json:"data,omitempty"`
	Error   []udata `json:"error,omitempty"`
}

// udoc represents an Uber hypermedia document.
type udoc struct {
	Uber ubody `json:"uber"`
}

var (
	logging = 0
	pitctx  = context.Background()
)

func init() {
	pitctx = context.WithValue(pitctx, "logger", log.New(os.Stdout, "pitd: ", log.LstdFlags))
	http.Handle("/", handlers.CompressHandler(handlers.LoggingHandler(os.Stdout, router())))
}

func main() {
	http.ListenAndServe(":3006", nil)
}

func router() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/_loglevel", http.Handler(ContextAdapter{ctx: pitctx, handler: ContextHandlerFunc(loglevel)})).Methods("POST")
	return r
}

// logglevel sets the desired level of logging
func loglevel(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(mkError("ServerError", "reason", fmt.Sprintf("Cannot read HTTP request body [%+v]", err)))
		return
	}

	re := regexp.MustCompile("level=([[:digit:]])")
	sm := re.FindStringSubmatch(string(body))
	if sm == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(mkError("ClientError", "reason", fmt.Sprintf("Expecting level={digit}, got %s", string(body))))
		return
	}

	logging, err = strconv.Atoi(sm[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(mkError("ServerError", "reason", fmt.Sprintf("Unable to convert log level [%+v]", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// mkError creates an Uber hypermedia document that represents an error.
func mkError(name, rel, value string) []byte {
	bs, err := json.Marshal(udoc{ubody{Version: "1.0", Error: []udata{udata{Name: name, Rel: []string{rel}, Value: value}}}})
	if err != nil {
		panic(err)
	}
	return bs
}
