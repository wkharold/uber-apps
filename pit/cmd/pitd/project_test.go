package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
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
	ctxfn       func() context.Context
	rc          int
	body        string
}

var ptes = []projecttest{
	{"empty project list", projectlist, "/projects", GET, "", noprojects, 200, testdata.EmptyProjectList},
	{"single project list", projectlist, "/projects", GET, "", oneproject, 200, testdata.OneProjectList},
	{"multi project list", projectlist, "/projects", GET, "", multiproject, 200, testdata.MultiProjectList},
	{"get unknown project", getproject, "/project/001", GET, "", multiproject, http.StatusNotFound, testdata.UnknownProjectError},
	{"get the only project", getproject, "/project/101", GET, "", oneproject, 200, testdata.Project101},
	{"get a project", getproject, "/project/102", GET, "", multiproject, 200, testdata.Project102},
}

func TestProjects(t *testing.T) {
	for _, pt := range ptes {
		ctx := pt.ctxfn()

		req, err := http.NewRequest(pt.method, pt.req, strings.NewReader(pt.payload))
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		router(ctx).ServeHTTP(w, req)

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

		db := ctx.Value("database").(*sql.DB)
		dropdb(db)
	}
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

func multiproject() context.Context {
	db := createdb("multiproject")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))
	ctx = context.WithValue(ctx, "logger", &leveledLogger{logger: log.New(os.Stdout, "pittest: ", log.LstdFlags), level: INFO})

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create a transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES 
	                      (101, "project one", "first test project", 1001),
	                      (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002);;`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES (1001, "owner@test.net"), (1002, "owner@test.io");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func noprojects() context.Context {
	db := createdb("noprojects")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))
	ctx = context.WithValue(ctx, "logger", &leveledLogger{logger: log.New(os.Stdout, "pittest: ", log.LstdFlags), level: INFO})
	return ctx
}

func oneproject() context.Context {
	db := createdb("oneproject")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))
	ctx = context.WithValue(ctx, "logger", &leveledLogger{logger: log.New(os.Stdout, "pittest: ", log.LstdFlags), level: INFO})

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create a transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES (101, "project one", "first test project", 1001);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES (1001, "owner@test.net"), (1002, "owner@test.io");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func createdb(dbname string) *sql.DB {
	db, err := sql.Open("ql", fmt.Sprintf("memory://%s.db", dbname))
	if err != nil {
		panic(fmt.Sprintf("cannot create database instance: [%+v]", err))
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("database ping failed: [%+v]", err))
	}

	if err = mkTables(db); err != nil {
		panic(fmt.Sprintf("table creation failed: [%+v]", err))
	}

	return db
}

func dropdb(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create a transaction to drop the database: [%+v]", err))
	}

	if _, err := tx.Exec("DROP TABLE projects"); err != nil {
		panic(fmt.Sprintf("cannot drop the projects table: [%+v]", err))
	}

	if _, err := tx.Exec("DROP TABLE issues"); err != nil {
		panic(fmt.Sprintf("cannot drop the issues table: [%+v]", err))
	}

	if _, err := tx.Exec("DROP TABLE members"); err != nil {
		panic(fmt.Sprintf("cannot drop the members table: [%+v]", err))
	}

	if _, err := tx.Exec("DROP TABLE contributors"); err != nil {
		panic(fmt.Sprintf("cannot drop the contributors table: [%+v]", err))
	}

	if _, err := tx.Exec("DROP TABLE assignments"); err != nil {
		panic(fmt.Sprintf("cannot drop the assignments table: [%+v]", err))
	}

	if _, err := tx.Exec("DROP TABLE watchers"); err != nil {
		panic(fmt.Sprintf("cannont drop the watchers table: [%+v]", err))
	}

	tx.Commit()
}

func mkTables(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Entity tables: projects, issues, members
	if _, err = tx.Exec("CREATE TABLE IF NOT EXISTS projects (ID int, Name string, Description string, Owner int);"); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("CREATE TABLE IF NOT EXISTS issues (ID int, Name string,  Description string, Priority int, Status string, Project int, Reporter int);"); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("CREATE TABLE IF NOT EXISTS members (ID int, Email string);"); err != nil {
		tx.Rollback()
		return err
	}

	// Association tables: contributors, assignments, watchers
	if _, err = tx.Exec("CREATE TABLE IF NOT EXISTS contributors (PID int, MID int);"); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("CREATE TABLE IF NOT EXISTS assignments (MID int, IID int);"); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("CREATE TABLE IF NOT EXISTS watchers (MID int, IID int);"); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
