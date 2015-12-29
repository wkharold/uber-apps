package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/uber-apps/pit/cmd/pitd/db"
	"github.com/uber-apps/pit/cmd/pitd/uber"

	"golang.org/x/net/context"
)

type links struct{}

type projects []db.Project

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
				Model:  "n={name}&d={description}",
				Data:   []uber.Data{},
			},
		},
	}, nil
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

	pl, err := db.FindAllProjects(ctx)
	if err != nil {
		rc, reason := http.StatusInternalServerError, "no projects in context"

		if logger.level == DEBUG {
			logger.logger.Printf("projectlist: exit with %d [%s]", http.StatusInternalServerError, reason)
		}

		w.WriteHeader(rc)
		w.Write(mkError("ServerError", "reason", reason))
	}

	ud, err := uber.Marshal(links(struct{}{}), projects(pl))
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

func (ps projects) MarshalUBER() (uber.Data, error) {
	return uber.Data{ID: "projects", Data: []uber.Data{}}, nil
}
