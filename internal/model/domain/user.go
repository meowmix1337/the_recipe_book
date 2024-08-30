package domain

type UserSignup struct {
	Email    string
	Password string
}

type User struct {
	ID    uint
	UUID  string
	Email string
}
