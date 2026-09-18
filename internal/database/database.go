package database

import (
	"fmt"

	"github.com/google/uuid"
)

type (
	Filter struct {
		Field string
		Id    uuid.UUID
	}
)

func (f *Filter) ToSql() string {
	return fmt.Sprintf("%s = '%s'", f.Field, f.Id)
}
