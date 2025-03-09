package auth

type (
	// SigninRequest struct is a data structure representing the user sign-in information.
	SigninRequest struct {
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required" json:"password"`
	}

	// SigninResponse struct represents the successful sign-in response.
	SigninResponse struct {
		Token        string      `json:"token"`
		User         interface{} `json:"user"`
		Organization interface{} `json:"organization"`
		Member       bool        `json:"member"`
	}
)
