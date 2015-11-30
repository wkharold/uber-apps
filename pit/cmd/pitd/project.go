package main

import (
	"net/http"

	"golang.org/x/net/context"
)

type projects []project

type project struct {
	Name        string
	Description string
}

func (ps projects) MarshalUber() (*udoc, error) {
	links := udata{
		ID: "links",
		Data: []udata{
			udata{
				ID:     "alps",
				Rel:    []string{"profile"},
				URL:    "/pit-alps.xml",
				Action: "read",
				Data:   []udata{},
			},
			udata{
				ID:     "list",
				Name:   "links",
				Rel:    []string{"collection"},
				URL:    "/projects/",
				Action: "read",
				Data:   []udata{},
			},
			udata{
				ID:        "search",
				Name:      "links",
				Rel:       []string{"search"},
				URL:       "/projects/search{?name}",
				Templated: true,
				Data:      []udata{},
			},
			udata{
				ID:     "new",
				Name:   "links",
				Rel:    []string{"add"},
				URL:    "/projects",
				Action: "append",
				Model:  "n={name}&d={description}",
				Data:   []udata{},
			},
		},
	}

	projlist := udata{ID: "projects", Data: []udata{}}

	return &udoc{ubody{Version: "1.0", Data: []udata{links, projlist}, Error: []udata{}}}, nil
}

func projectlist(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
