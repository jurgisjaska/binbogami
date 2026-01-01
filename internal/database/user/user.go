package user

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleDefault int = iota + 1
	RoleBilling
	RoleAdmin
)

type (
	User struct {
		Id          *uuid.UUID `json:"id"`
		Email       *string    `json:"email"`
		Name        *string    `json:"name"`
		Surname     *string    `json:"surname"`
		Salt        string     `json:"-"`
		Password    string     `json:"-"`
		Role        int        `json:"role"`
		CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
		UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
		ConfirmedAt *time.Time `db:"confirmed_at" json:"confirmedAt"`
		DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
	}

	Users []User
)
