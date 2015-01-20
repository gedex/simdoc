package model

type User struct {
	ID       int64  `meddler:"id,pk"    json:"-"`
	Login    string `meddler:"login"    validate:"login" json:"login"`
	Email    string `meddler:"email"    validate:"email" json:"email"`
	Password string `meddler:"password" validate:"min=6" json:"-"`
	Name     string `meddler:"name"     json:"name"`
	Role     string `meddler:"role"     validate:"role" json:"role"`
	Created  int64  `meddler:"created"  json:"created_at"`
	Updated  int64  `meddler:"updated"  json:"updated_at"`
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// IsAdmin checks whether user has admin role or not.
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
