package middleware

import (
	"github.com/gedex/simdoc/pkg/model"
	"github.com/zenazn/goji/web"
)

// UserToC sets the User in the current web context.
func UserToC(c *web.C, user *model.User) {
	c.Env["user"] = user
}

// ToUser returns the User from the current request context.
func ToUser(c *web.C) *model.User {
	var v = c.Env["user"]

	u, ok := v.(*model.User)
	if !ok {
		return nil
	}
	return u
}

// DocToC sets the Document in the current web context.
func DocToC(c *web.C, doc *model.Document) {
	c.Env["document"] = doc
}

// ToDoc returns the Document from the current request context.
func ToDoc(c *web.C) *model.Document {
	var v = c.Env["document"]

	d, ok := v.(*model.Document)
	if !ok {
		return nil
	}
	return d
}
