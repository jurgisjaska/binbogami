package api

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `validate:"required"`
	Age  int    `validate:"gte=1,lte=150"`
}

func TestValidator_Validate(t *testing.T) {

	tt := []struct {
		name   string
		in     interface{}
		hasErr bool
	}{
		{
			name: "Valid",
			in: &TestStruct{
				Name: "Example Name",
				Age:  25,
			},
			hasErr: false,
		},
		{
			name: "Empty Name",
			in: &TestStruct{
				Name: "",
				Age:  50,
			},
			hasErr: true,
		},
		{
			name: "Age Out of Range",
			in: &TestStruct{
				Name: "Example Name",
				Age:  151,
			},
			hasErr: true,
		},
		{
			name: "Negative Age",
			in: &TestStruct{
				Name: "Example Name",
				Age:  -1,
			},
			hasErr: true,
		},
	}

	v := &Validator{
		Validator: validator.New(),
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.hasErr, v.Validate(tc.in) != nil)
		})
	}
}
