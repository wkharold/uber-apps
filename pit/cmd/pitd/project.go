package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
)

type projects []project

type project struct {
	Name        string
	Description string
}

func projectlist(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	logger := ctx.Value("logger").(*leveledLogger)
	if logger == nil {
		devnull, _ := os.OpenFile("/dev/null", os.O_WRONLY, os.ModePerm)
		logger = &leveledLogger{logger: log.New(devnull, "nulllogger", log.LstdFlags), level: INFO}
	}

	if logger.level == DEBUG {
		logger.logger.Println("projectlist: enter")
	}

	pl := ctx.Value("projects").(*projects)
	if pl == nil {
		rc, reason := http.StatusInternalServerError, "no projects in context"

		if logger.level == DEBUG {
			logger.logger.Printf("projectlist: exit with %d [%s]", http.StatusInternalServerError, reason)
		}

		w.WriteHeader(rc)
		w.Write(mkError("ServerError", "reason", reason))
	}

	ud, err := pl.MarshalUber()
	if err != nil {
		rc, reason := http.StatusInternalServerError, fmt.Sprintf("Unable to marshal as UBER: %+v", err)

		if logger.level == DEBUG {
			logger.logger.Printf("projectlist: exit with %d [%s]", rc, reason)
		}

		w.WriteHeader(rc)
		w.Write(mkError("ServerError", "reason", reason))
	}

	if logger.level == DEBUG {
		logger.logger.Printf("projectlist: exit with 200 [%s]", string(ud))
	}

	w.Write(ud)
}

func (ps projects) MarshalUber() ([]byte, error) {
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
				URL:    "/projects/",
				Action: "append",
				Model:  "n={name}&d={description}",
				Data:   []udata{},
			},
		},
	}

	projlist := udata{ID: "projects", Data: []udata{}}

	return json.Marshal(udoc{ubody{Version: "1.0", Data: []udata{links, projlist}, Error: []udata{}}})
}
