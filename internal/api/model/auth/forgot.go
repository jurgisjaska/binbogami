package auth

type (
	// ForgotPasswordRequest struct is a data structure representing the request for forgot password feature,
	ForgotPasswordRequest struct {
		Email string `validate:"required,email" json:"email"`
	}
)
