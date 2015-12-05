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

	rows, err := db.Query("SELECT members.ID, members.Email FROM members")
	if err != nil {
		return []Member{}, err
	}

	return collectMembers(ctx, rows)
}

// FindByEmail retrieves the project team member with the given email address.
func (Members) FindByEmail(ctx context.Context, email string) (Member, error) {
	db := databaseFromContext(ctx)
	result := Member{}

	err := db.QueryRow("SELECT members.ID, members.Email FROM members WHERE members.Email == ?", email).Scan(&result.id, &result.email)
	if err != nil {
		return Member{}, err
	}

	return result, nil
}

// FindByID retrieves the project team member with the given id.
func (Members) FindByID(ctx context.Context, memberid int) (Member, error) {
	db := databaseFromContext(ctx)
	result := Member{}

	err := db.QueryRow("SELECT members.ID, members.Email FROM members WHERE members.ID == ?", memberid).Scan(&result.id, &result.email)
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
	return []Project{}, nil
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
