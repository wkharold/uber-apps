// There are no tests here, just some helper functions used by the various tests.

package db

import (
	"database/sql"
	"fmt"
)

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

func sameissues(a, b []Issue) bool {
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
func samemembers(a, b []Member) bool {
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

func sameprojects(a, b []Project) bool {
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
