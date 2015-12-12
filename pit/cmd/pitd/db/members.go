package db

import (
	"database/sql"
	"errors"

	"golang.org/x/net/context"
)

var (
	ErrMemberExists = errors.New("Member already exists")
)

// Member is a member of a project team.
type Member struct {
	id    int
	email string
}

// Members is the collection of all the project team members known to the PIT system.
type Members struct{}

// NewMember creates a new member associated with the specified email address.
func NewMember(ctx context.Context, email string) (Member, error) {
	db := databaseFromContext(ctx)
	ids := idsChanFromContext(ctx)
	id := <-ids

	tx, err := db.Begin()
	if err != nil {
		return Member{}, err
	}

	var memail string
	var mid int

	err = tx.QueryRow("SELECT ID, Email FROM members WHERE Email == $1", email).Scan(&mid, &memail)
	if err != sql.ErrNoRows {
		tx.Rollback()
		return Member{}, ErrMemberExists
	}

	_, err = tx.Exec("INSERT INTO members VALUES ($1, $2)", id, email)
	if err != nil {
		tx.Rollback()
		return Member{}, err
	}

	tx.Commit()

	return Member{id, email}, nil
}

// FindAll retreives a list of all the project team members in the repository.
func (Members) FindAll(ctx context.Context) ([]Member, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT ID, Email FROM members ORDER BY ID;")
	if err != nil {
		return []Member{}, err
	}

	return collectMembers(ctx, rows)
}

// FindByEmail retrieves the project team member with the given email address.
func (Members) FindByEmail(ctx context.Context, email string) (Member, error) {
	db := databaseFromContext(ctx)
	result := Member{}

	err := db.QueryRow("SELECT ID, Email FROM members WHERE Email == $1;", email).Scan(&result.id, &result.email)
	if err != nil {
		return Member{}, err
	}

	return result, nil
}

// FindByID retrieves the project team member with the given id.
func (Members) FindByID(ctx context.Context, memberid int) (Member, error) {
	db := databaseFromContext(ctx)
	result := Member{}

	err := db.QueryRow("SELECT ID, Email FROM members WHERE ID == $1", memberid).Scan(&result.id, &result.email)
	if err != nil {
		return Member{}, err
	}

	return result, nil
}

// Assignments retrieves a list of all the issues to which the proejct team member has been assigned.
func (m Member) Assignments(ctx context.Context) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query(`
	SELECT I.IID, I.Name, I.Description, I.Priority, I.Status, I.Project, I.Reporter
	FROM (SELECT issues.ID AS IID, issues.Name AS Name, issues.Description AS Description, issues.Priority AS Priority, issues.Status AS Status, issues.Project AS Project, members.Email AS Reporter
	      FROM issues, members
		  WHERE issues.Reporter == members.ID
		  ORDER BY IID) AS I
	FULL JOIN assignments ON (I.IID == assignments.IID)
    WHERE assignments.MID == $1
	ORDER BY I.IID
	`, m.id)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// ContributesTo retrieves a list of all the projects to which the project team member contributes.
func (m Member) ContributesTo(ctx context.Context) ([]Project, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query(`
	SELECT P.PID, P.Name, P.Description, P.Email
	FROM (SELECT projects.ID AS PID, projects.Name AS Name, projects.Description AS Description, members.Email AS Email
		  FROM projects, members
		  WHERE projects.Owner == members.ID
		  ORDER BY PID) AS P
	FULL JOIN contributors ON (P.PID == contributors.PID)
	WHERE contributors.MID == $1
	ORDER BY P.PID
	`, m.id)
	if err != nil {
		return []Project{}, err
	}

	return collectProjects(ctx, rows)
}

// Watching retrieves a list of all the issue the project team member is watching.
func (m Member) Watching(ctx context.Context) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query(`
	SELECT I.IID, I.Name, I.Description, I.Priority, I.Status, I.Project, I.Reporter
	FROM (SELECT issues.ID AS IID, issues.Name AS Name, issues.Description AS Description, issues.Priority AS Priority, issues.Status AS Status, issues.Project AS Project, members.Email AS Reporter
	      FROM issues, members
		  WHERE issues.Reporter == members.ID
		  ORDER BY IID) AS I
	FULL JOIN watchers ON (I.IID == watchers.IID)
    WHERE watchers.MID == $1
	ORDER BY I.IID
    `, m.id)
	if err != nil {
		return []Issue{}, nil
	}

	return collectIssues(ctx, rows)
}

func collectMembers(ctx context.Context, rows *sql.Rows) ([]Member, error) {
	result := []Member{}

	for rows.Next() {
		member := Member{}

		if err := rows.Scan(&member.id, &member.email); err != nil {
			return []Member{}, err
		}

		result = append(result, member)
	}

	return result, nil
}
