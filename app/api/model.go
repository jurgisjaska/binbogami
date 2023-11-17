package api

type (
	SigninModel struct {
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required" json:"password"`
	}

	SignupModel struct {
		Email            string `validate:"required" json:"email"`
		Password         string `validate:"required" json:"password"`
		RepeatedPassword string `validate:"required" json:"repeated_password"`
	}
)
