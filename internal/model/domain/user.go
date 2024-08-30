package domain

type UserSignup struct {
	Email    string
	Password string
}

type User struct {
	ID    uint
	Email string
}
