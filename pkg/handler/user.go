package handler

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gedex/simdoc/pkg/datastore"
	"github.com/gedex/simdoc/pkg/model"
	"github.com/gedex/simdoc/pkg/util"

	"github.com/goji/context"
	"github.com/zenazn/goji/web"
)

// GetUserByLogin accepts a request to retrieve user information, for given provided
// login or email, from the datastore and returns in JSON format.
//
// GET /api/users/:login
//
func GetUserByLogin(c web.C, w http.ResponseWriter, r *http.Request) {
	var ctx = context.FromC(c)
	var loginOrEmail = c.URLParams["login"]

	user, err := datastore.GetUserByLogin(ctx, loginOrEmail)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// GetAllUsers accepts a request to retrieve all users from the datastore and
// returns in JSON format.
//
// GET /api/users
//
func GetAllUsers(c web.C, w http.ResponseWriter, r *http.Request) {
	var ctx = context.FromC(c)

	users, err := datastore.GetAllUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if users == nil {
		w.Write([]byte(`[]`))
	} else {
		json.NewEncoder(w).Encode(users)
	}
}

// AddUser accepts a request to add new user into the datastore.
//
// POST /api/users
//
func AddUser(c web.C, w http.ResponseWriter, r *http.Request) {
	usr, err := parseSubmittedUser(r)
	if err != nil {
		respWithError(w, http.StatusBadRequest, ErrorInvalidJSONRequest)
		return
	}

	// Checks if login exists in datastore. No need to check the email, as both
	// login and email are unique.

	// Validate the model.
	if ve := model.Validate(usr); ve != nil {
		respWithError(w, http.StatusBadRequest, ErrorValidationFailed, getValidationErrors("users", ve)...)
		return
	}

	// Sets password.
	if usr.Password == "" {
		usr.Password = util.GetRandomString(8)
	}
	// Encrypt the password.
	usr.Password = encryptPassword(c, usr.Password)

	// @todo send notification to email informing about email being used in
	// this app. Email contains a link to verify the user.
	err = datastore.AddUser(context.FromC(c), usr)
	if err != nil {
		// @todo refactor this by checking the given login first from the datastore.
		// @todo also validate posted email and login.
		if isDuplicateLogin(err) {
			respWithError(w, http.StatusBadRequest, ErrorValidationFailed, newFieldError("users", "login or email already exists", ErrorFieldAlreadyExists))
		} else {
			log.Printf("%+v\n", err)
			respWithError(w, http.StatusBadRequest, ErrorValidationFailed)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usr)
}

//
// POST /api/user/login
//
func UserLogin(c web.C, w http.ResponseWriter, r *http.Request) {
	var loginInfo = new(struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	})

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(loginInfo); err != nil {
		respWithError(w, http.StatusBadRequest, ErrorInvalidJSONRequest)
		return
	}

	var fe []*fieldError
	if loginInfo.Login == "" {
		fe = append(fe, newFieldError("user", "login", ErrorFieldMissing))
	}
	if loginInfo.Password == "" {
		fe = append(fe, newFieldError("user", "password", ErrorFieldMissing))
	}
	if len(fe) > 0 {
		respWithError(w, http.StatusBadRequest, ErrorValidationFailed, fe...)
		return
	}

	var ctx = context.FromC(c)
	var salt = c.Env["passwdSalt"].(string)

	// Encrypt the password.
	encPass := util.PBKDF2([]byte(loginInfo.Password), []byte(salt), 10000, 50, sha256.New)
	loginInfo.Password = fmt.Sprintf("%x", encPass)

	usr, err := datastore.GetUserByLoginAndPassword(ctx, loginInfo.Login, loginInfo.Password)
	if err != nil || usr == nil {
		respWithError(w, http.StatusBadRequest, ErrorBadCredentials)
		return
	}

	tokenStr, err := util.GenerateToken(ctx, r, usr)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, ErrorInternalServerError)
		return
	}

	resp := struct {
		*model.User
		Token string `json:"token"`
	}{
		usr,
		tokenStr,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func GetCurrentUser(c web.C, w http.ResponseWriter, r *http.Request) {
	var user = ToUser(c)
	if user == nil {
		respWithError(w, http.StatusUnauthorized, ErrorRequireAuthentication)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateCurrentUser
//
// PATCH /api/user
// PUT   /api/user
//
func UpdateCurrentUser(c web.C, w http.ResponseWriter, r *http.Request) {
	var cusr = ToUser(c)
	if cusr == nil {
		respWithError(w, http.StatusUnauthorized, ErrorRequireAuthentication)
		return
	}

	usr, err := parseSubmittedUser(r)
	if err != nil {
		respWithError(w, http.StatusBadRequest, ErrorInvalidJSONRequest)
		return
	}

	usr.ID = cusr.ID

	// Login can not be changed.
	if usr.Login != "" && usr.Login != cusr.Login {
		respWithError(w, http.StatusBadRequest, ErrorValidationFailed, newFieldError("user", "login", ErrorFieldImmutable))
		return
	}

	// Use current email if not supplied.
	if usr.Email == "" {
		usr.Email = cusr.Email
	}

	usr.Login = cusr.Login
	usr.Password = cusr.Password // Update password should be handled separately.
	usr.Created = cusr.Created   // Created is immutable

	// If current user is not an admin, prevent updating role.
	if !cusr.IsAdmin() {
		usr.Role = cusr.Role
	}

	// Validate the model.
	if ve := model.Validate(usr); ve != nil {
		log.Printf("%+v\n", usr)
		respWithError(w, http.StatusBadRequest, ErrorValidationFailed, getValidationErrors("user", ve)...)
		return
	}

	// @todo send notification to email informing updated profile information.
	err = datastore.UpdateUser(context.FromC(c), usr)
	if err != nil {
		respWithError(w, http.StatusBadRequest, ErrorValidationFailed)
		return
	}

	json.NewEncoder(w).Encode(usr)
}

func GetCurrentUserDocuments(c web.C, w http.ResponseWriter, r *http.Request) {
	var user = ToUser(c)
	if user == nil {
		respWithError(w, http.StatusUnauthorized, ErrorRequireAuthentication)
		return
	}

	json.NewEncoder(w)
}

// parseSubmittedUser returns User through POST or PUT.
func parseSubmittedUser(r *http.Request) (*model.User, error) {
	decoder := json.NewDecoder(r.Body)

	usrWrapper := new(struct {
		model.User

		Password string `meddler:"password" json:"password"`
	})
	err := decoder.Decode(usrWrapper)
	if err != nil {
		return nil, err
	}

	usr := &usrWrapper.User
	if usrWrapper.Password != "" {
		usr.Password = usrWrapper.Password
	}

	return usr, nil
}

func isDuplicateLogin(err error) bool {
	if !isDuplicateEntryError(err) {
		return false
	}
	if strings.Contains(err.Error(), "for key 'login'") || strings.Contains(err.Error(), "for key 'email'") {
		return true
	}
	return false
}

func isDuplicateEntryError(err error) bool {
	if strings.Contains(err.Error(), "Duplicate entry") {
		return true
	}
	return false
}

func encryptPassword(c web.C, pass string) string {
	var salt = c.Env["passwdSalt"].(string)

	return fmt.Sprintf("%x", util.PBKDF2([]byte(pass), []byte(salt), 10000, 50, sha256.New))
}
