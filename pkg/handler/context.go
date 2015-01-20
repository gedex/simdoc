package handler

import (
	"github.com/gedex/simdoc/pkg/model"
	"github.com/zenazn/goji/web"
)

// ToUser returns the User from the current request context. If the user does not
// exists a nil is returned.
func ToUser(c web.C) *model.User {
	var v = c.Env["user"]
	if v == nil {
		return nil
	}
	if usr, ok := v.(*model.User); ok {
		return usr
	}
	return nil
}

// ToDocument returns the Document from the current request context. If the Document
// does not exists a nil value is returned.
func ToDocument(c web.C) *model.Document {
	var v = c.Env["document"]

	if v == nil {
		return nil
	}
	if doc, ok := v.(*model.Document); ok {
		return doc
	}
	return nil
}
