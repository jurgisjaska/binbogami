package api

type (
	Signin struct {
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required" json:"password"`
	}

	Signup struct {
		Email            string `validate:"required,email" json:"email"`
		Password         string `validate:"required,gte=8" json:"password"`
		RepeatedPassword string `validate:"required,gte=8" json:"repeated_password"`
		Name             string `validate:"required,gte=3" json:"name"`
		Surname          string `validate:"required,gte=3" json:"surname"`
	}
)
