package auth

// ForgotRequest struct is a data structure representing the request for forgot password feature,
type ForgotRequest struct {
	Email string `validate:"required,email" json:"email"`

	User      interface{}
	Ip        string
	UserAgent string
}
