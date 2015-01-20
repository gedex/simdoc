package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/gedex/simdoc/pkg/model"
)

type Userstore interface {
	// GetUserById retrieves a specific user from the datastore for the given ID.
	GetUserById(id int64) (*model.User, error)

	// GetUserByLogin retrieves a user from the datastore for the specified login
	// (username) or email.
	GetUserByLogin(loginOrEmail string) (*model.User, error)

	GetUserByLoginAndPassword(loginOrEmail, password string) (*model.User, error)

	// GetAllUsers retrieves a list of all users from the datastore.
	GetAllUsers() ([]*model.User, error)

	// AddUser adds a user into the datastore.
	AddUser(user *model.User) error

	// UpdateUser update a user in the datastore.
	UpdateUser(user *model.User) error

	// DeleteUser deletes a user, for the given ID, in the datastore.
	DeleteUser(id int64) error
}

// GetUserById retrieves a specific user from the datastore for the given ID.
func GetUserById(c context.Context, id int64) (*model.User, error) {
	return FromContext(c).GetUserById(id)
}

// GetUserByLogin retrieves a user from the datastore for the specified login
// (username) or email.
func GetUserByLogin(c context.Context, loginOrEmail string) (*model.User, error) {
	return FromContext(c).GetUserByLogin(loginOrEmail)
}

func GetUserByLoginAndPassword(c context.Context, loginOrEmail, password string) (*model.User, error) {
	return FromContext(c).GetUserByLoginAndPassword(loginOrEmail, password)
}

// GetAllUsers retrieves a list of all users from the datastore.
func GetAllUsers(c context.Context) ([]*model.User, error) {
	return FromContext(c).GetAllUsers()
}

// AddUser adds a user into the datastore.
func AddUser(c context.Context, user *model.User) error {
	return FromContext(c).AddUser(user)
}

// UpdateUser update a user in the datastore.
func UpdateUser(c context.Context, user *model.User) error {
	return FromContext(c).UpdateUser(user)
}

// DeleteUser deletes a user, for the given ID, in the datastore.
func DeleteUser(c context.Context, id int64) error {
	return FromContext(c).DeleteUser(id)
}
