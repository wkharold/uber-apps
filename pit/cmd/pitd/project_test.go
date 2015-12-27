package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/uber-apps/pit/cmd/pitd/httpctx"
	"github.com/uber-apps/pit/cmd/pitd/testdata"
	"golang.org/x/net/context"
)

const (
	GET  = "GET"
	POST = "POST"
)

type projecttest struct {
	description string
	hfn         httpctx.ContextHandlerFunc
	req         string
	method      string
	payload     string
	ctx         context.Context
	rc          int
	body        string
}

var ptes = []projecttest{
	{"empty project list", projectlist, "/projects", GET, "", noprojects(), 200, testdata.EmptyProjectList},
}

func TestProjects(t *testing.T) {
	for _, pt := range ptes {
		req, err := http.NewRequest(pt.method, pt.req, strings.NewReader(pt.payload))
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		pt.hfn(pt.ctx, w, req)

		if w.Code != pt.rc {
			t.Errorf("%s: Response Code mismatch: expected %d, got %d", pt.description, pt.rc, w.Code)
			continue
		}

		if len(pt.body) == 0 {
			continue
		}

		if equaljson(w.Body.Bytes(), []byte(pt.body)) == false {
			body := bytes.NewBuffer([]byte{})
			json.Compact(body, []byte(pt.body))
			t.Errorf("%s: Body mismatch:\nexpected %s\ngot      %s", pt.description, string(body.Bytes()), w.Body.String())
			continue
		}
	}
}

func noprojects() context.Context {
	ctx := context.WithValue(context.Background(), "projects", &projects{})
	ctx = context.WithValue(ctx, "logger", &leveledLogger{logger: log.New(os.Stdout, "pittest: ", log.LstdFlags), level: DEBUG})
	return ctx
}

func equaljson(p, q []byte) bool {
	cp := bytes.NewBuffer([]byte{})

	if err := json.Compact(cp, p); err != nil {
		log.Printf("unable to compact cp json for equaljson: %+v", err)
		return false
	}

	cq := bytes.NewBuffer([]byte{})

	if err := json.Compact(cq, q); err != nil {
		log.Printf("unable to compact cq json for equaljson: %+v", err)
		return false
	}

	if len(cp.Bytes()) != len(cq.Bytes()) {
		return false
	}

	cpb, cqb := cp.Bytes(), cq.Bytes()

	for i, b := range cpb {
		if b != cqb[i] {
			return false
		}
	}

	return true
}
