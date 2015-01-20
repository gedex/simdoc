package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gedex/simdoc/pkg/datastore"
	"github.com/gedex/simdoc/pkg/model"
	"github.com/gedex/simdoc/pkg/util"

	"github.com/goji/context"
	"github.com/zenazn/goji/web"
)

// DevSeedUser accepts a request to seed admin and regular users.
//
// GET /api/dev/seed_users
//
func DevSeedUsers(c web.C, w http.ResponseWriter, r *http.Request) {

	stub := [][]string{
		[]string{"admin01", "admin01@local.host", "password", model.RoleAdmin},
		[]string{"user01", "user01@local.host", "password", model.RoleUser},
	}

	for _, u := range stub {
		_, err := seedUser(c, u...)
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
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stub)
}

func seedUser(c web.C, fields ...string) (*model.User, error) {
	usr := &model.User{
		Login:    fields[0],
		Email:    fields[1],
		Password: fields[2],
		Role:     fields[3],
	}

	// Sets password.
	if usr.Password == "" {
		usr.Password = util.GetRandomString(8)
	}
	// Encrypt the password.
	usr.Password = encryptPassword(c, usr.Password)

	// Add new user into the datastore.
	if err := datastore.AddUser(context.FromC(c), usr); err != nil {
		return nil, err
	}

	return usr, nil
}

// DevSeedDocuments accepts a request to seed documents.
//
// GET /api/dev/seed_documents
//
func DevSeedDocuments(c web.C, w http.ResponseWriter, r *http.Request) {
	// First check if we've seeded users.
	user, err := datastore.GetUserByLogin(context.FromC(c), "user01")
	if err != nil {
		respWithError(w, http.StatusNotAcceptable, errors.New("No user01 found. Please seed via /api/dev/seed_users."))
		return
	}

	stub := []*model.Document{
		&model.Document{Name: "Document-1", Status: model.DocumentStatusDraft, CreatedBy: user.ID},
		&model.Document{Name: "Document-2", Status: model.DocumentStatusPublished, CreatedBy: user.ID},
	}

	docs := make([]*model.Document, len(stub))
	for _, d := range stub {
		doc, err := seedDoc(c, d)
		if err != nil {
			log.Println(err)
			respWithError(w, http.StatusBadRequest, ErrorValidationFailed)
			return
		}
		docs = append(docs, doc)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stub)
}

func seedDoc(c web.C, doc *model.Document) (*model.Document, error) {
	// Add Document into the datastore.
	if err := datastore.AddDocument(context.FromC(c), doc); err != nil {
		return nil, err
	}

	return doc, nil
}
