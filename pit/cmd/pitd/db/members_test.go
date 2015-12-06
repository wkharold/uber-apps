package db

import (
	"database/sql"
	"fmt"
	"testing"

	"golang.org/x/net/context"
)

type findAllMembersTest struct {
	description string
	fn          func(Members, context.Context) ([]Member, error)
	ctxfn       func() context.Context
	expected    []Member
	err         error
}

var (
	bob   = Member{id: 1003, email: "bob@members.com"}
	carol = Member{id: 1004, email: "carol@members.com"}
	ted   = Member{id: 1005, email: "ted@members.com"}
	alice = Member{id: 1006, email: "alice@members.com"}

	findAllMemberTests = []findAllMembersTest{
		{"FindAll no members", Members.FindAll, emptytables, []Member{}, nil},
		{"FindAll one member", Members.FindAll, onemember, []Member{bob}, nil},
		{"FindAll members", Members.FindAll, manymembers, []Member{carol, ted, alice}, nil},
	}
)

func TestFindAllMembers(t *testing.T) {
	for _, nt := range findAllMemberTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		ms, err := nt.fn(struct{}{}, ctx)
		if err != nil {
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
			dropdb(db)
			continue
		}

		if !samemembers(nt.expected, ms) {
			t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, ms)
			dropdb(db)
			continue
		}

		dropdb(db)
	}
}

func onemember() context.Context {
	db := createdb("onemember")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES (1003, "bob@members.com");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func manymembers() context.Context {
	db := createdb("manymembers")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES 
						  (1004, "carol@members.com"),
						  (1005, "ted@members.com"),
						  (1006, "alice@members.com");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}
