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
	pone = Project{
		id: 101, name: "project one", description: "first test project", owner: "owner@test.net",
	}
	ptwo = Project{
		id: 102, name: "project two", description: "second test project", owner: "owner@test.net",
	}
	pthree = Project{
		id: 103, name: "project three", description: "third test project", owner: "owner@test.io",
	}

	rowsfindertests = []rowsfindertest{
		{"FindAll from empty tables", Projects.FindAll, emptytables(), []Project{}},
		{"FindAll one project", Projects.FindAll, oneproject(), []Project{pone}},
		{"FindAll multiple projects", Projects.FindAll, manyprojects(), []Project{pone, ptwo, pthree}},
	}
	rowfindertests = []rowfindertest{}
)

func TestRowsFinders(t *testing.T) {
	for _, rft := range rowsfindertests {
		db := rft.ctx.Value("database").(*sql.DB)

		ps, err := rft.ffn(struct{}{}, rft.ctx)
		if err != nil {
			t.Errorf("%s: unexpected error [%+v]", rft.description, err)
			dropdb(db)
			continue
		}

		if !rowsequal(rft.expected, ps) {
			t.Errorf("%s: expected %+v, got %+v", rft.description, rft.expected, ps)
			dropdb(db)
			continue
		}

		dropdb(db)
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

func oneproject() context.Context {
	db := createdb("oneproject")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create a transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES (101, "project one", "first test project", 1001);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES (1001, "owner@test.net");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func manyprojects() context.Context {
	db := createdb("manyprojects")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create a transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES 
						  (101, "project one", "first test project", 1001),
						  (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES 
						  (1001, "owner@test.net"),
						  (1002, "owner@test.io");`); err != nil {
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

func rowsequal(a, b []Project) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil || b == nil:
		return false
	case len(a) != len(b):
		return false
	default:
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}
