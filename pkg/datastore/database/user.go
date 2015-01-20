package database

import (
	"time"

	"github.com/gedex/simdoc/pkg/model"
	"github.com/russross/meddler"
)

type Userstore struct {
	meddler.DB
}

func NewUserstore(db meddler.DB) *Userstore {
	return &Userstore{db}
}

func (db *Userstore) GetUserById(id int64) (*model.User, error) {
	var usr = new(model.User)
	var err = meddler.Load(db, userTable, usr, id)

	return usr, err
}

func (db *Userstore) GetUserByLogin(loginOrEmail string) (*model.User, error) {
	var usr = new(model.User)
	var err = meddler.QueryRow(db, usr, userByLoginQuery, loginOrEmail, loginOrEmail)

	return usr, err
}

func (db *Userstore) GetUserByLoginAndPassword(loginOrEmail, password string) (*model.User, error) {
	var usr = new(model.User)
	var err = meddler.QueryRow(db, usr, userByLoginAndPasswordQuery, loginOrEmail, loginOrEmail, password)

	return usr, err
}

func (db *Userstore) GetAllUsers() ([]*model.User, error) {
	var users []*model.User
	var err = meddler.QueryAll(db, &users, userListQuery)

	return users, err
}

func (db *Userstore) AddUser(user *model.User) error {
	if user.Created == 0 {
		user.Created = time.Now().UTC().Unix()
	}
	user.Updated = time.Now().UTC().Unix()

	return meddler.Save(db, userTable, user)
}

func (db *Userstore) UpdateUser(user *model.User) error {
	user.Updated = time.Now().UTC().Unix()

	return meddler.Save(db, userTable, user)
}

func (db *Userstore) DeleteUser(id int64) error {
	var _, err = db.Exec(userDeleteQuery, id)

	return err
}

const userTable = "users"

const userByLoginQuery = `
SELECT * FROM users
WHERE login=? OR email=? LIMIT 1
`

const userByLoginAndPasswordQuery = `
SELECT * FROM users
WHERE
	(login=? OR email=?)
	AND
	password=?
LIMIT 1
`

const userListQuery = `
SELECT * FROM users
ORDER BY login ASC
`

const userDeleteQuery = `
DELETE FROM users
WHERE id=?
`
