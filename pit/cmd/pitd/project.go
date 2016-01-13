package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/uber-apps/pit/cmd/pitd/db"
	"github.com/uber-apps/pit/cmd/pitd/uber"

	"golang.org/x/net/context"
)

type links struct{}

type projects []db.Project
type project db.Project

func (ls links) MarshalUBER() (uber.Data, error) {
	return uber.Data{
		ID: "links",
		Data: []uber.Data{
			uber.Data{
				ID:     "alps",
				Rel:    []string{"profile"},
				URL:    "/pit-alps.xml",
				Action: "read",
				Data:   []uber.Data{},
			},
			uber.Data{
				ID:     "list",
				Name:   "links",
				Rel:    []string{"collection"},
				URL:    "/projects/",
				Action: "read",
				Data:   []uber.Data{},
			},
			uber.Data{
				ID:        "search",
				Name:      "links",
				Rel:       []string{"search"},
				URL:       "/projects/search{?name}",
				Templated: true,
				Data:      []uber.Data{},
			},
			uber.Data{
				ID:     "new",
				Name:   "links",
				Rel:    []string{"add"},
				URL:    "/projects/",
				Action: "append",
				Model:  "n={name}&d={description}&o={owner}",
				Data:   []uber.Data{},
			},
		},
	}, nil
}

func addproject(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := loggerFromContext(ctx)

	logger.Log(DEBUG, "addproject: %s", "enter")

	logger.Log(DEBUG, "addproject: %s", "exit")
	w.WriteHeader(http.StatusNotImplemented)
}

func getproject(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := loggerFromContext(ctx)

	logger.Log(DEBUG, "project: %s", "enter")

	vars := mux.Vars(req)
	id := vars["id"]

	pid, err := strconv.Atoi(id)
	if err != nil {
		logger.Log(DEBUG, "getproject: strconv.Atoi(%d) failed [%+v]", id, err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(mkError("ServerError", "reason", fmt.Sprintf("Project ID must be an integer not: [%s]", id)))

		logger.Log(DEBUG, "project: %s", "exit")
		return
	}

	p, err := db.FindProjectByID(ctx, pid)
	switch {
	case err == sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write(mkError("RequestFailed", "reason", fmt.Sprintf("No project exists with specified ID: [%d]", pid)))

		logger.Log(DEBUG, "project: %s", "exit")
		return
	case err != nil:
		logger.Log(DEBUG, "getproject: db.FindProjectByID(ctx, %d) failed [%+v]", pid, err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(mkError("ServerError", "reason", fmt.Sprintf("Project lookup error: [%+v]", err)))

		logger.Log(DEBUG, "project: %s", "exit")
		return
	default:
		ud, err := uber.Marshal(links(struct{}{}), project(p))
		if err != nil {
			logger.Log(DEBUG, "getproject: uber.Marshal(...uber.Marshaler) failed [%+v]", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(mkError("ServerError", "reason", fmt.Sprintf("Unable to marshal response [%+v]", err)))

			logger.Log(DEBUG, "project: %s", "exit")
			return
		}

		w.Write(ud)
	}

	logger.Log(DEBUG, "project: %s", "exit")
}

func (p project) MarshalUBER() (uber.Data, error) {
	dbp := db.Project(p)
	pdata := uber.Data{
		ID:   strconv.Itoa(dbp.ID()),
		Name: dbp.Name(),
		Rel:  []string{"self"},
		URL:  fmt.Sprintf("/project/%d", dbp.ID()),
		Data: []uber.Data{
			{
				Rel:    []string{"add"},
				URL:    fmt.Sprintf("/project/%d/issues", dbp.ID()),
				Action: "append",
				Model:  "n={name}&d={description}&p={priority}&r={reporter}",
			},
			{
				Rel:       []string{"search"},
				URL:       fmt.Sprintf("/project/%d/search{?name}", dbp.ID()),
				Templated: true,
			},
			{Name: "description", Value: dbp.Description()},
			{Name: "owner", Value: dbp.Owner()},
		},
	}

	return uber.Data{ID: "project", Data: []uber.Data{pdata}}, nil
}

func projectlist(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := loggerFromContext(ctx)

	logger.Log(DEBUG, "projectlist: %s", "enter")

	pl, err := db.FindAllProjects(ctx)
	if err != nil {
		rc, reason := http.StatusInternalServerError, "no projects in context"

		logger.Log(DEBUG, "projectlist: exit with %d [%s]", http.StatusInternalServerError, reason)

		w.WriteHeader(rc)
		w.Write(mkError("ServerError", "reason", reason))
	}

	ud, err := uber.Marshal(links(struct{}{}), projects(pl))
	if err != nil {
		rc, reason := http.StatusInternalServerError, fmt.Sprintf("Unable to marshal as UBER: %+v", err)

		logger.Log(DEBUG, "projectlist: exit with %d [%s]", rc, reason)

		w.WriteHeader(rc)
		w.Write(mkError("ServerError", "reason", reason))
	}

	logger.Log(DEBUG, "projectlist: exit with 200 [%s]", string(ud))

	w.Write(ud)
}

func (ps projects) MarshalUBER() (uber.Data, error) {
	summaries := []uber.Data{}

	for _, p := range ps {
		s := uber.Data{
			ID:   strconv.Itoa(p.ID()),
			Name: p.Name(),
			Rel:  []string{"self"},
			URL:  fmt.Sprintf("/project/%d", p.ID()),
			Data: []uber.Data{
				{
					Rel:    []string{"add"},
					URL:    fmt.Sprintf("/project/%d/issues", p.ID()),
					Action: "append",
					Model:  "n={name}&d={description}&p={priority}&r={reporter}",
				},
				{
					Rel:       []string{"search"},
					URL:       fmt.Sprintf("/project/%d/search{?name}", p.ID()),
					Templated: true,
				},
			},
		}
		summaries = append(summaries, s)
	}

	return uber.Data{ID: "projects", Data: summaries}, nil
}
