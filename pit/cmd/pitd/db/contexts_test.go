// The db contexts used by the various tests are here.

package db

import (
	"fmt"

	"golang.org/x/net/context"
)

func emptytables() context.Context {
	db := createdb("emptytables")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

	return ctx
}

func alltheissues() context.Context {
	db := createdb("alltheissues")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES 
						  (101, "project one", "first test project", 1001),
						  (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002),
						  (104, "project four", "fourth test project", 1001),
						  (105, "project five", "fifth test project", 1002),
						  (106, "project six", "sixth test project", 1001);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES 
						  (1001, "owner@test.net"),
						  (1002, "owner@test.io"),
						  (1003, "bob@members.com"),
						  (1004, "carol@members.com"),
						  (1005, "ted@members.com"),
						  (1006, "alice@members.com"),
						  (1007, "fred@testrock.org"),
						  (1008, "barney@testrock.org"),
						  (1009, "wilma@testrock.org");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO contributors VALUES
						  (102, 1006),
						  (103, 1004),
						  (103, 1005),
						  (103, 1006);`); err != nil {
		panic(fmt.Sprintf("cannot setup contributors table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO issues VALUES
	                      (2001, "issueone", "issue one", 1, "OPEN", 101, 1007),
						  (2002, "issuetwo", "issue two", 2, "OPEN", 102, 1008),
						  (2003, "issuethree", "issue three", 2, "CLOSED", 102, 1008),
						  (2004, "issuefour", "issue four", 3, "RETURNED", 104, 1007),
						  (2005, "issuefive", "issue five", 2, "OPEN", 103, 1007),
						  (2006, "issuesix", "issue six", 4, "CLOSED", 106, 1009);`); err != nil {
		panic(fmt.Sprintf("cannot setup issues table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO assignments VALUES
	                      (1003, 2002),
						  (1004, 2005),
						  (1004, 2006),
						  (1005, 2006),
						  (1006, 2006);`); err != nil {
		panic(fmt.Sprintf("cannot setup assignments table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO watchers VALUES
	                      (1006, 2005),
	                      (1003, 2002),
						  (1003, 2003),
						  (1003, 2006),
	                      (1004, 2003),
						  (1005, 2003);`); err != nil {
		panic(fmt.Sprintf("cannot setup watchers table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func noissues() context.Context {
	db := createdb("noissues")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES 
						  (101, "project one", "first test project", 1001),
						  (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002),
						  (104, "project four", "fourth test project", 1001),
					      (105, "project five", "fifth test project", 1002),
						  (106, "project six", "sixth test project", 1001);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES 
						  (1001, "owner@test.net"),
						  (1002, "owner@test.io"),
						  (1003, "bob@members.com"),
						  (1004, "carol@members.com"),
						  (1005, "ted@members.com"),
						  (1006, "alice@members.com"),
						  (1007, "fred@testrock.org"),
						  (1008, "barney@testrock.org");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO contributors VALUES
						  (102, 1006),
						  (103, 1004),
						  (103, 1005),
						  (103, 1006);`); err != nil {
		panic(fmt.Sprintf("cannot setup contributors table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func oneissue() context.Context {
	db := createdb("alltheissues")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES 
						  (101, "project one", "first test project", 1001),
						  (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002),
						  (104, "project four", "fourth test project", 1001),
					      (105, "project five", "fifth test project", 1002),
						  (106, "project six", "sixth test project", 1001);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES 
						  (1001, "owner@test.net"),
						  (1002, "owner@test.io"),
						  (1003, "bob@members.com"),
						  (1004, "carol@members.com"),
						  (1005, "ted@members.com"),
						  (1006, "alice@members.com"),
						  (1007, "fred@testrock.org"),
						  (1008, "barney@testrock.org");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO contributors VALUES
						  (102, 1006),
						  (103, 1004),
						  (103, 1005),
						  (103, 1006);`); err != nil {
		panic(fmt.Sprintf("cannot setup contributors table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO issues VALUES
	                      (2001, "issueone", "issue one", 1, "OPEN", 101, 1007);`); err != nil {
		panic(fmt.Sprintf("cannot setup issues table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func contributions() context.Context {
	db := createdb("contributions")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES
						  (101, "project one", "first test project", 1001),
						  (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES
						  (1001, "owner@test.net"),
						  (1002, "owner@test.io"),
						  (1003, "bob@members.com"),
						  (1004, "carol@members.com"),
						  (1005, "ted@members.com"),
						  (1006, "alice@members.com");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO contributors VALUES
						  (101, 1006),
						  (101, 1004),
						  (103, 1006);`); err != nil {
		panic(fmt.Sprintf("cannot setup contributors table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func onemember() context.Context {
	db := createdb("onemember")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

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
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

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

func contributors() context.Context {
	db := createdb("contributors")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES 
						  (101, "project one", "first test project", 1001),
						  (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES 
						  (1001, "owner@test.net"),
						  (1002, "owner@test.io"),
						  (1003, "bob@members.com"),
						  (1004, "carol@members.com"),
						  (1005, "ted@members.com"),
						  (1006, "alice@members.com");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO contributors VALUES
						  (102, 1006),
						  (103, 1004),
						  (103, 1005),
						  (103, 1006);`); err != nil {
		panic(fmt.Sprintf("cannot setup contributors table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func manyprojects() context.Context {
	db := createdb("manyprojects")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

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

func noprojects() context.Context {
	db := createdb("oneproject")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create a transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES (1001, "owner@test.net");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func oneproject() context.Context {
	db := createdb("oneproject")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

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

func projectissues() context.Context {
	db := createdb("projectissues")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

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

	if _, err := tx.Exec(`INSERT INTO issues VALUES
	                      (2002, "issuetwo", "issue two", 2, "OPEN", 102, 1008),
						  (2011, "issueone", "issue one", 2, "OPEN", 103, 1008);`); err != nil {
		panic(fmt.Sprintf("cannot setup issues table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES 
						  (1001, "owner@test.net"),
						  (1002, "owner@test.io"),
						  (1003, "bob@members.com"),
						  (1004, "carol@members.com"),
						  (1005, "ted@members.com"),
						  (1006, "alice@members.com"),
						  (1007, "fred@testrock.org"),
						  (1008, "barney@testrock.org"),
						  (1009, "wilma@testrock.org");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}
