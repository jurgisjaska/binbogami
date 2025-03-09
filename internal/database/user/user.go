package user

import (
	"time"

	"github.com/google/uuid"
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
)
