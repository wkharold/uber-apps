package db

import (
	"database/sql"
	"fmt"
	"testing"

	"golang.org/x/net/context"
)

type rowsfindertest struct {
	description string
	ffn         func(Projects, context.Context) ([]Project, error)
	ctx         context.Context
	expected    []Project
}

type rowfindertest struct {
	description string
	ffn         func(Projects, ctx context.Context) (Project, error)
	ctx         context.Context
	expected    Project
}

var (
	rowsfindertests = []rowsfindertest{
		{"FindAll from empty tables", Projects.FindAll, emptytables(), []Project{}},
	}
	rowfindertests = []rowfindertest{}
)

func TestRowsFinders(t *testing.T) {
	for _, rft := range rowsfindertests {
		ps, err := rft.ffn(struct{}{}, rft.ctx)
		if err != nil {
			t.Errorf("%s: unexpected error [%+v]", rft.description, err)
			continue
		}

		if !rowsequal(rft.expected, ps) {
			t.Errorf("%s: expected %+v, got %+v", rft.description, rft.expected, ps)
			continue
		}
	}
}

func TestRowFinders(t *testing.T) {
	for _, _ = range rowfindertests {
	}
}

func emptytables() context.Context {
	db := createdb("emptytables")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

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

func rowsequal(a, b []Project) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
