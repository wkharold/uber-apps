package db

import (
	"database/sql"
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
	findAllMemberTests = []findAllMembersTest{}
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
