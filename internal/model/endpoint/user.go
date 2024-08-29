package endpoint

type UserSignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
