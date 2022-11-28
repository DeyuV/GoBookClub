package users

type UserResponse struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Password  string
	Enabled   bool
	Role      string
}
