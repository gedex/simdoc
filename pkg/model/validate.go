package model

import (
	"errors"
	"regexp"

	validator "gopkg.in/validator.v2"
)

var (
	ErrorInvalidLogin          = errors.New("Invalid login format")
	ErrorInvalidEmail          = errors.New("Invalid email format")
	ErrorInvalidRole           = errors.New("Invalid user role")
	ErrorInvalidDocumentStatus = errors.New("Invalid document status")
)

func init() {
	validator.SetValidationFunc("login", validateLogin)
	validator.SetValidationFunc("email", validateEmail)
	validator.SetValidationFunc("role", validateRole)
	validator.SetValidationFunc("doc_status", validateDocumentStatus)
}

func Validate(v interface{}) error {
	switch vv := v.(type) {
	case *User:
		return vv.Validate()
	case *Document:
		return vv.Validate()
	default:
		return validator.ErrUnsupported
	}
}

func (u *User) Validate() error {
	return validator.Validate(u)
}

func (d *Document) Validate() error {
	return validator.Validate(d)
}

func validateLogin(v interface{}, param string) error {
	vv, ok := v.(string)
	if !ok {
		return ErrorInvalidLogin
	}
	if len(vv) < 5 {
		return ErrorInvalidLogin
	}
	if exp, _ := regexp.Compile("^[a-zA-Z_0-9]+$"); !exp.MatchString(vv) {
		return ErrorInvalidLogin
	}
	return nil
}

func validateEmail(v interface{}, param string) error {
	vv, ok := v.(string)
	if !ok {
		return ErrorInvalidEmail
	}
	if exp, _ := regexp.Compile("^.+\\@.+\\..+$"); !exp.MatchString(vv) {
		return ErrorInvalidEmail
	}
	return nil
}

func validateRole(v interface{}, param string) error {
	vv, ok := v.(string)
	if !ok {
		return ErrorInvalidRole
	}
	if vv != RoleUser && vv != RoleAdmin {
		return ErrorInvalidRole
	}
	return nil
}

func validateDocumentStatus(v interface{}, param string) error {
	switch v.(type) {
	case string:
		vv := v.(string)
		if vv != DocumentStatusDraft && vv != DocumentStatusPublished {
			return ErrorInvalidDocumentStatus
		}
	default:
		return ErrorInvalidDocumentStatus
	}

	return nil
}
