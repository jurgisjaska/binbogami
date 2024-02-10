package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/database/book"
)

const (
	MemberRoleDefault = iota + 1 // standard member of the organization
	MemberRoleBilling            // member with rights to manage billing information
	MemberRoleAdmin              // organization administrator, can invite other member
	MemberRoleOwner              // organization owner
)

type (
	Member struct {
		Id             int        `json:"id"`
		Role           int        `json:"role"`
		OrganizationId *uuid.UUID `db:"organization_id" json:"organization_id"`
		UserId         *uuid.UUID `db:"user_id" json:"user_id"`

		CreatedBy *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	MemberRepository struct {
		database *sqlx.DB
	}
)

func (r *MemberRepository) Create(org *uuid.UUID, user *uuid.UUID, role int, createdBy *uuid.UUID) (*Member, error) {
	member := &Member{
		OrganizationId: org,
		UserId:         user,
		Role:           role,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
	}

	_, err := r.database.NamedExec(`
			INSERT INTO members (id, organization_id, user_id, role, created_by, created_at)
			VALUES (NULL, :organization_id, :user_id, :role, :created_by, :created_at)
		`, member)

	if err != nil {
		return nil, err
	}

	return member, nil
}

func (r *MemberRepository) Find(org *uuid.UUID, user *uuid.UUID) (*Member, error) {
	query := `
		SELECT members.* 
		FROM members 
		JOIN organizations AS o ON members.organization_id = o.id
		JOIN users AS u ON members.user_id = u.id
		WHERE 
		    members.organization_id = ? AND members.user_id = ? 
		    AND members.deleted_at IS NULL
			AND u.deleted_at IS NULL AND o.deleted_at IS NULL
	`

	member := &Member{}
	if err := r.database.Get(member, query, org, user); err != nil {
		return nil, err
	}

	return member, nil
}

// ByBook finds the member record using book that belongs to the organization and the user.
func (r *MemberRepository) ByBook(book *book.Book, user *uuid.UUID) (*Member, error) {
	query := `
		SELECT members.* 
		FROM members 
		JOIN organizations AS o ON members.organization_id = o.id
		JOIN users AS u ON members.user_id = u.id
		JOIN books AS b ON o.id = b.organization_id
		WHERE 
		    b.id = ? AND members.user_id = ? 
		    AND members.deleted_at IS NULL
			AND u.deleted_at IS NULL AND o.deleted_at IS NULL AND b.deleted_at IS NULL
	`

	member := &Member{}
	if err := r.database.Get(member, query, book.Id, user); err != nil {
		return nil, err
	}

	return member, nil
}

func CreateMember(d *sqlx.DB) *MemberRepository {
	return &MemberRepository{database: d}
}
