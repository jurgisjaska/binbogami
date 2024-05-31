package auth

type (
	ResetPasswordRequest struct {
		Password         string `validate:"required,gte=8" json:"password"`
		RepeatedPassword string `validate:"required,gte=8,eqfield=Password" json:"repeatedPassword"`
	}
)
