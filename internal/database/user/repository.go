package user

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	database *sqlx.DB
}

// Find retrieves a User from the database by its UUID.
func (r *Repository) Find(id uuid.UUID) (*User, error) {
	user := &User{}
	err := r.database.Get(user, "SELECT * FROM users WHERE id = ? AND deleted_at IS NULL", id.String())
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindActiveByEmail retrieves a user from the database based on their email address, ensuring they are active
// (not deleted and confirmed).
func (r *Repository) FindActiveByEmail(e string) (*User, error) {
	return r.findByEmail(e, "SELECT * FROM users WHERE email = ? AND deleted_at IS NULL AND confirmed_at IS NOT NULL")
}

// FindByEmail retrieves a user from the database based on their email address.
func (r *Repository) FindByEmail(e string) (*User, error) {
	return r.findByEmail(e, "SELECT * FROM users WHERE email = ?")
}

// findByEmail is a private helper method to retrieve a user by their email address using a custom query.
func (r *Repository) findByEmail(e string, q string) (*User, error) {
	user := &User{}
	err := r.database.Get(user, q, e)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) FindMany(filter string) (*Users, error) {
	users := &Users{}

	if len(filter) == 0 {
		query := "SELECT * FROM users WHERE deleted_at IS NULL"
		err := r.database.Select(users, query)
		if err != nil {
			return nil, err
		}

		return users, nil
	}

	// @todo this is a horrible way to search for things
	query := `
		SELECT * FROM users 
		WHERE (email LIKE ? OR CONCAT(users.name, ' ', users.surname) LIKE ?) AND deleted_at IS NULL
	 `
	filter = fmt.Sprintf("%%%s%%", filter)

	err := r.database.Select(users, query, filter, filter)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Create inserts a new user into the database.
func (r *Repository) Create(u *User) error {
	query := `
		INSERT INTO users (id, email, name, surname, salt, role, password, created_at, confirmed_at)
		VALUES (:id, :email, :name, :surname, :salt, :role, :password, :created_at, :confirmed_at) 
	`

	_, err := r.database.NamedExec(query, u)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdatePassword(u *User) error {
	_, err := r.database.NamedExec(`UPDATE users SET password = :password WHERE id = :id`, u)
	if err != nil {
		return err
	}

	return nil
}

// CreateUser creates a new instance of the Repository with the specified SQL database connection.
func CreateUser(d *sqlx.DB) *Repository {
	return &Repository{database: d}
}
