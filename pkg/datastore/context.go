package datastore

import (
	"code.google.com/p/go.net/context"
)

const reqKey = "datastore"

type wrapper struct {
	context.Context
	ds Datastore
}

// NewContext returns a Context whose Value method returns the application's data
// storage objects.
func NewContext(parent context.Context, ds Datastore) context.Context {
	return &wrapper{parent, ds}
}

// Value returns the named key from the context.
func (c *wrapper) Value(key interface{}) interface{} {
	if key == reqKey {
		return c.ds
	}
	return c.Context.Value(key)
}

// FromContext returns the sql.DB associated with this context.
func FromContext(c context.Context) Datastore {
	return c.Value(reqKey).(Datastore)
}
