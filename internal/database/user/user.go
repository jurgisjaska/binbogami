package user

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleDefault int = iota + 1
	RoleReadonly
	RoleBilling
	RoleAdmin
	RoleManager
	RoleOwner
)

type (
	User struct {
		Id          uuid.UUID  `json:"id"`
		Email       string     `json:"email"`
		Name        string     `json:"name"`
		Surname     string     `json:"surname"`
		Position    *string    `json:"position"`
		Salt        string     `json:"-"`
		Password    string     `json:"-"`
		Role        int        `json:"role"`
		CreatedAt   time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
		ConfirmedAt *time.Time `db:"confirmed_at" json:"confirmed_at"`
		DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	Users []User
)
