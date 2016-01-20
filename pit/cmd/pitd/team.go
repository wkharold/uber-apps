package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/uber-apps/pit/cmd/pitd/db"
	"github.com/uber-apps/pit/cmd/pitd/uber"
	"golang.org/x/net/context"
)

type members []db.Member
type member db.Member

func addmember(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := loggerFromContext(ctx)

	logger.Log(DEBUG, "addmember: %s", "enter")

	ml, err := db.FindAllMembers(ctx)
	switch {
	case err == sql.ErrNoRows:
		ud, err := uber.Marshal(links(struct{}{}))
		if err != nil {
			writeError("addmember", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to marshal as UBER: %+v", err))
			return
		}

		w.Write(ud)
		return
	case err != nil:
		writeError("addmember", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to retrieve members: %+v", err))
		return
	default:
		ud, err := uber.Marshal(links(struct{}{}), members(ml))
		if err != nil {
			writeError("addmember", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to marshal as UBER: %+v", err))
			return
		}

		w.Write(ud)
		return
	}
}

func (m member) MarshalUBER() (uber.Data, error) {

	mdata := uber.Data{}

	return mdata, nil
}

func teamlist(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := loggerFromContext(ctx)

	logger.Log(DEBUG, "teamlist: %s", "enter")

	ml, err := db.FindAllMembers(ctx)
	switch {
	case err == sql.ErrNoRows:
		ud, err := uber.Marshal(links(struct{}{}))
		if err != nil {
			writeError("teamlist", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to marshal as UBER: %+v", err))
			return
		}

		rc := http.StatusOK

		w.WriteHeader(rc)
		w.Write(ud)

		logger.Log(DEBUG, "teamlist: exit with %d: [%s]", rc, string(ud))
		return
	case err != nil:
		writeError("teamlist", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to retrieve members: %+v", err))
		return
	default:
		ud, err := uber.Marshal(links(struct{}{}), members(ml))
		if err != nil {
			writeError("teamlist", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to marshal as UBER: %+v", err))
			return
		}

		rc := http.StatusOK

		w.WriteHeader(rc)
		w.Write(ud)

		logger.Log(DEBUG, "teamlist: exit with %d: [%s]", rc, string(ud))
		return
	}
}

func (ms members) MarshalUBER() (uber.Data, error) {
	md := []uber.Data{}

	for _, m := range ms {
		d := uber.Data{
			ID:  strconv.Itoa(m.ID()),
			Rel: []string{"self"},
			URL: fmt.Sprintf("/team/%d", m.ID()),
			Data: []uber.Data{
				{Name: "email", Value: m.Email()},
			},
		}

		md = append(md, d)
	}

	return uber.Data{ID: "members", Data: md}, nil
}
