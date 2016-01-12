// Package main provides a simple UBER hypermedia driven project/issue tracking service.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/uber-apps/pit/cmd/pitd/httpctx"
	"github.com/uber-apps/pit/cmd/pitd/uber"
	"golang.org/x/net/context"
)

const (
	DEBUG = 0
	INFO  = 1
)

// leveledLogger combines a std log.Logger with a level at which to log. Possible levels
// are DEBUG, things that matter to developers, and INFO, things that matter to users.
type leveledLogger struct {
	logger *log.Logger
	level  int
}

var (
	pitctx = context.Background()
)

func init() {
	pitctx = context.WithValue(pitctx, "logger", &leveledLogger{logger: log.New(os.Stdout, "pitd: ", log.LstdFlags), level: INFO})
	http.Handle("/", handlers.CompressHandler(handlers.LoggingHandler(os.Stdout, router(pitctx))))
}

func main() {
	http.ListenAndServe(":3006", nil)
}

func router(ctx context.Context) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/_loglevel", http.Handler(httpctx.ContextAdapter{Ctx: ctx, Handler: httpctx.ContextHandlerFunc(loglevel)})).Methods("POST")
	r.Handle("/projects", http.Handler(httpctx.ContextAdapter{Ctx: ctx, Handler: httpctx.ContextHandlerFunc(projectlist)})).Methods("GET")
	r.Handle("/project/{id}", http.Handler(httpctx.ContextAdapter{Ctx: ctx, Handler: httpctx.ContextHandlerFunc(getproject)})).Methods("GET")
	return r
}

// loglevel sets the desired level of logging
func loglevel(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(mkError("ServerError", "reason", fmt.Sprintf("Cannot read HTTP request body [%+v]", err)))
		return
	}

	re := regexp.MustCompile("level=([[:alpha:]]+)")
	sm := re.FindStringSubmatch(string(body))
	if sm == nil || len(sm) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(mkError("ClientError", "reason", fmt.Sprintf("Expecting level={digit}, got %s", string(body))))
		return
	}

	switch strings.ToLower(sm[1]) {
	case "debug":
		ll := ctx.Value("logger").(*leveledLogger)
		if ll != nil {
			ll.level = DEBUG
		}
	case "info":
		ll := ctx.Value("logger").(*leveledLogger)
		if ll != nil {
			ll.level = INFO
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(mkError("ServerError", "reason", fmt.Sprintf("Unable to convert log level [%+v]", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// mkError creates an Uber hypermedia document that represents an error.
func mkError(name, rel, value string) []byte {
	ubererr, err := uber.MarshalError(uber.Data{Name: name, Rel: []string{rel}, Value: value})
	if err != nil {
		panic(err)
	}
	return ubererr
}

// loggerFromContext retrieves a leveledLogger from the given context; if no leveledLogger is present a null logger
// which writes to /dev/null is returned
func loggerFromContext(ctx context.Context) *leveledLogger {
	logger, ok := ctx.Value("logger").(*leveledLogger)
	if !ok {
		devnull, _ := os.OpenFile("/dev/null", os.O_WRONLY, os.ModePerm)
		logger = &leveledLogger{logger: log.New(devnull, "nulllogger", log.LstdFlags), level: INFO}
	}

	return logger
}

// Log generates a log message if the specified level less than or equal to the level in force at the time of the call.
func (ll leveledLogger) Log(level int, msg string, args ...interface{}) {
	if ll.level <= level {
		ll.logger.Printf(msg, args)
	}
}
