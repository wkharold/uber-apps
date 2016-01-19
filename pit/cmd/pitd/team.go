package main

import (
	"net/http"

	"github.com/uber-apps/pit/cmd/pitd/db"
	"golang.org/x/net/context"
)

type members []db.Member
type member db.Member

func addmember(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := loggerFromContext(ctx)

	logger.Log(DEBUG, "addmember: %s", "enter")

	writeError("addmember", w, logger, "NotImplemented", http.StatusNotImplemented, "addmember: unimplemented request")
}

func teamlist(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := loggerFromContext(ctx)

	logger.Log(DEBUG, "teamlist: %s", "enter")

	writeError("teamlist", w, logger, "NotImplemented", http.StatusNotImplemented, "teamlist: unimplemented request")
}
