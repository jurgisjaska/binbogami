package database

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/random"
)

type (
	User struct {
		Id        *uuid.UUID `json:"id"`
		Email     *string    `json:"email"`
		Name      *string    `json:"name"`
		Surname   *string    `json:"surname"`
		Salt      string     `json:"-"`
		Password  string     `json:"-"`
		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	Users []User

	UserRepository struct {
		database *sqlx.DB
	}
)

func (r *UserRepository) FindBy(column string, email string) (*User, error) {
	user := &User{}
	sql := fmt.Sprintf("SELECT * FROM users WHERE %s = ? AND deleted_at IS NULL", column)
	err := r.database.Get(user, sql, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Create(user *User) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	user.Id = &id
	user.CreatedAt = time.Now()

	hash := sha256.Sum256([]byte(
		fmt.Sprintf(
			"%s%s%s",
			user.Email,
			time.Now().Format(time.RFC3339Nano),
			random.String(16),
		),
	))

	user.Salt = string(hash[:])

	return nil
}

func CreateUser(d *sqlx.DB) *UserRepository {
	return &UserRepository{database: d}
}
