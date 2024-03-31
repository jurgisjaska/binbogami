package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	User struct {
		Id        *uuid.UUID `json:"id"`
		Email     *string    `json:"email"`
		Name      *string    `json:"name"`
		Surname   *string    `json:"surname"`
		Salt      string     `json:"-"`
		Password  string     `json:"-"`
		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
	}

	Users []User

	// Repository struct represents a repository for managing user data in a database.
	Repository struct {
		database *sqlx.DB
	}
)

func (r *Repository) By(column string, value interface{}) (*User, error) {
	user := &User{}
	sql := fmt.Sprintf("SELECT * FROM users WHERE %s = ? AND deleted_at IS NULL", column)
	err := r.database.Get(user, sql, value)
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

func (r *Repository) Create(user *User) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	user.Id = &id
	user.CreatedAt = time.Now()

	_, err = r.database.NamedExec(`
		INSERT INTO users (id, email, name, surname, salt, password, created_at)
		VALUES (:id, :email, :name, :surname, :salt, :password, :created_at) 
	`, user)

	if err != nil {
		return err
	}

	return nil
}

// CreateUser creates a new instance of the Repository with the specified SQL database connection.
func CreateUser(d *sqlx.DB) *Repository {
	return &Repository{database: d}
}
