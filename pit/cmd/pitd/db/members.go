package db

import (
	"database/sql"

	"golang.org/x/net/context"
)

// Member is a member of a project team.
type Member struct {
	id    int
	email string
}

// Members is the collection of all the project team members known to the PIT system.
type Members struct{}

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
	return []Issue{}, nil
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
	return []Issue{}, nil
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