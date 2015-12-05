package db

import (
	"database/sql"
	"fmt"
	"testing"

	"golang.org/x/net/context"
)

type contributorstest struct {
	description string
	fn          func(context.Context) ([]Member, error)
	ctxfn       func() context.Context
	expected    []Member
	err         error
}

type findalltest struct {
	description string
	ffn         func(Projects, context.Context) ([]Project, error)
	ctxfn       func() context.Context
	expected    []Project
	err         error
}

type findbyownertest struct {
	description string
	ffn         func(Projects, context.Context, string) ([]Project, error)
	owner       string
	ctxfn       func() context.Context
	expected    []Project
	err         error
}

type findbyidtest struct {
	description string
	ffn         func(Projects, context.Context, int) (Project, error)
	id          int
	ctxfn       func() context.Context
	expected    Project
	err         error
}

type findbynametest struct {
	description string
	ffn         func(Projects, context.Context, string) (Project, error)
	name        string
	ctxfn       func() context.Context
	expected    Project
	err         error
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

	contributorstests = []contributorstest{
		{"Contributors no contributors", pone.Contributors, contributors, []Member{}, nil},
	}
	findalltests = []findalltest{
		{"FindAll from empty tables", Projects.FindAll, emptytables, []Project{}, nil},
		{"FindAll one project", Projects.FindAll, oneproject, []Project{pone}, nil},
		{"FindAll multiple projects", Projects.FindAll, manyprojects, []Project{pone, ptwo, pthree}, nil},
	}
	findbyownertests = []findbyownertest{
		{"FindByOwner from empty tables", Projects.FindByOwner, "owner@test.net", emptytables, []Project{}, nil},
		{"FindByOwner one project no match", Projects.FindByOwner, "owner@test.org", oneproject, []Project{}, nil},
		{"FindByOwner multiple projects no match", Projects.FindByOwner, "owner@test.com", manyprojects, []Project{}, nil},
		{"FindByOwner one project", Projects.FindByOwner, "owner@test.net", oneproject, []Project{pone}, nil},
		{"FindByOwner multiple projects one match", Projects.FindByOwner, "owner@test.io", manyprojects, []Project{pthree}, nil},
		{"FindByOwner multiple projects", Projects.FindByOwner, "owner@test.net", manyprojects, []Project{pone, ptwo}, nil},
	}
	findbyidtests = []findbyidtest{
		{"FindByID empty tables", Projects.FindByID, 42, emptytables, Project{}, sql.ErrNoRows},
		{"FindByID multiple projects none match", Projects.FindByID, 42, manyprojects, Project{}, sql.ErrNoRows},
		{"FindByID one project", Projects.FindByID, 101, oneproject, pone, nil},
		{"FindByID multiple projects", Projects.FindByID, 103, manyprojects, pthree, nil},
	}
	findbynametests = []findbynametest{
		{"FindByName empty tables", Projects.FindByName, "unknown", emptytables, Project{}, sql.ErrNoRows},
		{"FindByName multiple projects none match", Projects.FindByName, "unknown", manyprojects, Project{}, sql.ErrNoRows},
		{"FindByName one project", Projects.FindByName, "project one", oneproject, pone, nil},
		{"FindByName multiple projects", Projects.FindByName, "project two", manyprojects, ptwo, nil},
	}
)

func TestFindAll(t *testing.T) {
	for _, nt := range findalltests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		ps, err := nt.ffn(struct{}{}, ctx)
		if err != nil {
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
			dropdb(db)
			continue
		}

		if !rowsequal(nt.expected, ps) {
			t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, ps)
			dropdb(db)
			continue
		}

		dropdb(db)
	}
}

func TestFindByOwner(t *testing.T) {
	for _, nt := range findbyownertests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		ps, err := nt.ffn(struct{}{}, ctx, nt.owner)
		if err != nil {
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
			dropdb(db)
			continue
		}

		if !rowsequal(nt.expected, ps) {
			t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, ps)
			dropdb(db)
			continue
		}

		dropdb(db)
	}
}

func TestFindByID(t *testing.T) {
	for _, nt := range findbyidtests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		p, err := nt.ffn(struct{}{}, ctx, nt.id)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil && err == nt.err:
			break
		default:
			if p != nt.expected {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, p)
				break
			}
		}

		dropdb(db)
	}
}

func TestFindByName(t *testing.T) {
	for _, nt := range findbynametests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		p, err := nt.ffn(struct{}{}, ctx, nt.name)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil && err == nt.err:
			break
		default:
			if p != nt.expected {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, p)
				break
			}
		}

		dropdb(db)
	}
}

func contributors() context.Context {
	db := createdb("contributors")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

	return ctx
}

func emptytables() context.Context {
	db := createdb("emptytables")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

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
