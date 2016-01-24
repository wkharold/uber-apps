package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writeError("addmember", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Cannot read HTTP request body [%+v]", err))
		return
	}

	re := regexp.MustCompile("m=(.+@.+)")
	sm := re.FindStringSubmatch(string(body))
	if sm == nil || len(sm) < 2 {
		writeError("addmember", w, logger, "ClientError", http.StatusBadRequest, fmt.Sprintf("Member specification must be of the form: \"m={email}\" not [%s]", string(body)))
		return
	}

	_, err = db.NewMember(ctx, sm[1])
	switch {
	case err == db.ErrMemberExists:
		writeError("addmember", w, logger, "DuplicateMember", http.StatusConflict, fmt.Sprintf("Member exists [%s]", sm[1]))
		return
	case err != nil:
		writeError("addmember", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to create new project [%+v]", err))
		return
	default:
		w.WriteHeader(http.StatusCreated)
	}

	logger.Log(DEBUG, "addmember: exis with %d", http.StatusCreated)
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

// MarshalUBER generates the UBER representation of a list of project team members.
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
