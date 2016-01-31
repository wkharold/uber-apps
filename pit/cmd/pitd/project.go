package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/uber-apps/pit/cmd/pitd/db"
	"github.com/uber-apps/pit/cmd/pitd/uber"

	"golang.org/x/net/context"
)

type links struct{}

type issues struct {
	il  []db.Issue
	pid int
}

type issue db.Issue

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
				ID:     "project-list",
				Name:   "links",
				Rel:    []string{"collection"},
				URL:    "/projects/",
				Action: "read",
				Data:   []uber.Data{},
			},
			uber.Data{
				ID:        "project-search",
				Name:      "links",
				Rel:       []string{"search"},
				URL:       "/projects/search{?name}",
				Templated: true,
				Action:    "read",
				Data:      []uber.Data{},
			},
			uber.Data{
				ID:     "project-create",
				Name:   "links",
				Rel:    []string{"add"},
				URL:    "/projects/",
				Action: "append",
				Model:  "n={name}&d={description}&o={owner}",
				Data:   []uber.Data{},
			},
			uber.Data{
				ID:     "team-members-list",
				Name:   "links",
				Rel:    []string{"collection"},
				URL:    "/team",
				Action: "read",
				Data:   []uber.Data{},
			},
			uber.Data{
				ID:     "team-member-create",
				Name:   "links",
				Rel:    []string{"add"},
				URL:    "/team",
				Action: "append",
				Model:  "m={email}",
				Data:   []uber.Data{},
			},
		},
	}, nil
}

func addissue(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := loggerFromContext(ctx)
	rc := http.StatusCreated // happy path return code

	logger.Log(DEBUG, "addissue: %s", "enter")

	vars := mux.Vars(req)
	id := vars["id"]

	pid, err := strconv.Atoi(id)
	if err != nil {
		logger.Log(DEBUG, "addissue: strconv.Atoi(%d) failed [%+v]", id, err)
		writeError("addissue", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Project ID must be an integer not: [%s]", id))
		return
	}

	p, err := db.FindProjectByID(ctx, pid)
	switch {
	case err == sql.ErrNoRows:
		writeError("addissue", w, logger, "RequestFailed", http.StatusNotFound, fmt.Sprintf("No project exists with specified ID: [%d]", pid))
		return
	case err != nil:
		logger.Log(DEBUG, "addissue: db.FindProjectByID(ctx, %d) failed [%+v]", pid, err)
		writeError("addissue", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Project lookup error: [%+v]", err))
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writeError("addissue", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Cannot read HTTP request body [%+v]", err))
		return
	}

	re := regexp.MustCompile("n=([([:word:][:space:])]+)&d=([([:word:][:space:])]+)&p=([[:digit:]])&r=(.+@.+)")
	sm := re.FindStringSubmatch(string(body))
	if sm == nil || len(sm) < 4 {
		writeError("addissue", w, logger, "ClientError", http.StatusBadRequest, fmt.Sprintf("Issue specification must be of the form: \"n={name}&d={description}&p={priority}&r={reporter}\" not [%s]", string(body)))
		return
	}

	priority, err := strconv.Atoi(sm[3])
	if err != nil {
		writeError("addissue", w, logger, "ClientError", http.StatusBadRequest, fmt.Sprintf("Priority must be a single digit integer [%+v]", err))
		return
	}

	name, desc, reporter := sm[1], sm[2], sm[4]

	_, err = p.OpenIssue(ctx, name, desc, reporter, priority)
	switch {
	case err == db.ErrIssueExists:
		writeError("addissue", w, logger, "IssueExists", http.StatusConflict, fmt.Sprintf("Cannot create duplicate issues [%s]", sm[1]))
		return
	case err == db.ErrNoSuchMember:
		writeError("addissue", w, logger, "NoSuchMember", http.StatusBadRequest, fmt.Sprintf("Issue reporter is not a project team member [%s]", reporter))
		return
	case err != nil:
		writeError("addissue", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Cannot open issue [%+v]", err))
		return
	default:
		w.WriteHeader(rc)
	}

	logger.Log(DEBUG, "addissue: exit with %d", rc)
}

func issuelist(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := loggerFromContext(ctx)
	rc := http.StatusOK

	logger.Log(DEBUG, "issuelist: %s", "enter")

	vars := mux.Vars(req)
	id := vars["id"]

	pid, err := strconv.Atoi(id)
	if err != nil {
		logger.Log(DEBUG, "issuelist: strconv.Atoi(%d) failed [%+v]", id, err)
		writeError("issuelist", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Project ID must be an integer not: [%s]", id))
		return
	}

	_, err = db.FindProjectByID(ctx, pid)
	switch {
	case err == sql.ErrNoRows:
		writeError("issuelist", w, logger, "RequestFailed", http.StatusNotFound, fmt.Sprintf("No project exists with specified ID: [%d]", pid))
		return
	case err != nil:
		logger.Log(DEBUG, "issuelist: db.FindProjectByID(ctx, %d) failed [%+v]", pid, err)
		writeError("issuelist", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Project lookup error: [%+v]", err))
		return
	}

	il, err := db.FindIssuesByProject(ctx, pid)
	switch {
	case err == sql.ErrNoRows:
		return
	case err != nil:
		logger.Log(DEBUG, "issuelist: db.FindIssuesByProject(ctx, %d) failed [%+v]", pid, err)
		writeError("issuelist", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Project lookup error: [%+v]", err))
		return
	default:
		ud, err := uber.Marshal(issues{il: il, pid: pid})
		if err != nil {
			writeError("issuelist", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to marshal as UBER: %+v", err))
			return
		}

		logger.Log(DEBUG, "issuelist: exit with %d [%s]", rc, string(ud))

		w.Write(ud)
	}
}

func (is issues) MarshalUBER() (uber.Data, error) {
	ic := []uber.Data{}

	for _, i := range is.il {
		ud := uber.Data{
			ID:     strconv.Itoa(i.ID()),
			Name:   i.Name(),
			Rel:    []string{"self"},
			URL:    fmt.Sprintf("/project/%d/issue/%d", is.pid, i.ID()),
			Action: "read",
			Data: []uber.Data{
				{Rel: []string{"close"}, URL: fmt.Sprintf("/project/%d/issue/close", is.pid), Action: "append", Model: fmt.Sprintf("i=%d", i.ID())},
				{Rel: []string{"return"}, URL: fmt.Sprintf("/project/%d/issue/return", is.pid), Action: "append", Model: fmt.Sprintf("i=%d", i.ID())},
				{Rel: []string{"assign"}, URL: fmt.Sprintf("/project/%d/issue/%d/assign", is.pid, i.ID()), Action: "append", Model: "m={member}"},
				{Name: "description", Value: i.Description()},
				{Name: "priority", Value: strconv.Itoa(i.Priority())},
				{Name: "status", Value: i.Status()},
				{Name: "reporter", Value: i.Reporter()},
			},
		}
		ic = append(ic, ud)
	}

	return uber.Data{
		ID:     "issues",
		Rel:    []string{"self"},
		URL:    fmt.Sprintf("/project/%d/issues", is.pid),
		Action: "read",
		Data:   ic}, nil
}

func addproject(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := loggerFromContext(ctx)

	logger.Log(DEBUG, "addproject: %s", "enter")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writeError("addproject", w, logger, "SeverError", http.StatusInternalServerError, fmt.Sprintf("Cannot read HTTP request body [%+v]", err))
		return
	}

	re := regexp.MustCompile("n=([([:word:][:space:])]+)&d=([([:word:][:space:])]+)&o=(.+@.+)")
	sm := re.FindStringSubmatch(string(body))
	if sm == nil || len(sm) < 4 {
		writeError("addproject", w, logger, "ClientError", http.StatusBadRequest, fmt.Sprintf("Project specification must be of the form: \"n={name}&d={description}&o={owner}\" not [%s]", string(body)))
		return
	}

	_, err = db.NewProject(ctx, sm[1], sm[2], sm[3])
	if err != nil {
		writeError("addproject", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to create new project [%+v]", err))
		return
	}

	w.WriteHeader(http.StatusCreated)

	logger.Log(DEBUG, "addproject: exit with %d", http.StatusCreated)
}

func findproject(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := loggerFromContext(ctx)

	logger.Log(DEBUG, "findproject: %s", "enter")

	pnm, ok := req.URL.Query()["name"]
	if !ok {
		writeError("findproject", w, logger, "ClientError", http.StatusBadRequest, fmt.Sprintf("Request must include a \"name\" query parameter"))
		return
	}

	p, err := db.FindProjectByName(ctx, pnm[0])
	switch {
	case err == sql.ErrNoRows:
		writeError("findproject", w, logger, "NoSuchProject", http.StatusNotFound, fmt.Sprintf("Cannot find a project named: %s", pnm[0]))
		return
	case err != nil:
		writeError("findproject", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to lookup project %s: [%v]", pnm[0], err))
		return
	default:
		ud, err := uber.Marshal(links(struct{}{}), project(p))
		if err != nil {
			logger.Log(INFO, "findproject: uber.Marshal(...uber.Marshaler) failed [%+v", err)
			writeError("findproject", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to marshal response [%+v]", err))
			return
		}

		w.Write(ud)
	}

	logger.Log(DEBUG, "findproject: exiting with %d", http.StatusOK)
}

func getproject(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := loggerFromContext(ctx)

	logger.Log(DEBUG, "getproject: %s", "enter")

	vars := mux.Vars(req)
	id := vars["id"]

	pid, err := strconv.Atoi(id)
	if err != nil {
		logger.Log(DEBUG, "getproject: strconv.Atoi(%d) failed [%+v]", id, err)
		writeError("getproject", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Project ID must be an integer not: [%s]", id))
		return
	}

	p, err := db.FindProjectByID(ctx, pid)
	switch {
	case err == sql.ErrNoRows:
		writeError("getproject", w, logger, "RequestFailed", http.StatusNotFound, fmt.Sprintf("No project exists with specified ID: [%d]", pid))
		return
	case err != nil:
		logger.Log(DEBUG, "getproject: db.FindProjectByID(ctx, %d) failed [%+v]", pid, err)
		writeError("getproject", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Project lookup error: [%+v]", err))
		return
	default:
		ud, err := uber.Marshal(links(struct{}{}), project(p))
		if err != nil {
			logger.Log(DEBUG, "getproject: uber.Marshal(...uber.Marshaler) failed [%+v]", err)
			writeError("getproject", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to marshal response [%+v]", err))
			return
		}

		w.Write(ud)
	}

	logger.Log(DEBUG, "getproject: exit with %d", http.StatusOK)
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
		writeError("projectlist", w, logger, "ServerError", http.StatusInternalServerError, "no projects in context")
		return
	}

	ud, err := uber.Marshal(links(struct{}{}), projects(pl))
	if err != nil {
		writeError("projectlist", w, logger, "ServerError", http.StatusInternalServerError, fmt.Sprintf("Unable to marshal as UBER: %+v", err))
		return
	}

	logger.Log(DEBUG, "projectlist: exit with %d [%s]", http.StatusOK, string(ud))

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
