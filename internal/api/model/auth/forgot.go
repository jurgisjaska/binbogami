package auth

type (
	// ForgotRequest struct is a data structure representing the request for forgot password feature,
	ForgotRequest struct {
		Email string `validate:"required,email" json:"email"`

		User      interface{}
		Ip        string
		UserAgent string
	}
)
