package router

import (
	"github.com/gedex/simdoc/pkg/handler"
	"github.com/gedex/simdoc/pkg/middleware"

	"github.com/zenazn/goji/web"
)

func New() *web.Mux {
	mux := web.New()

	// Public endpoints.
	mux.Post("/api/user/login", handler.UserLogin)

	// @todo pass UserAuthorizer middleware on following endpoints.
	// @todo authenticated user should be able to get list of all users.
	// @todo however adding user to users must have 'admin' role.
	mux.Get("/api/users/:login", handler.GetUserByLogin)
	mux.Get("/api/users", handler.GetAllUsers)
	mux.Post("/api/users", handler.AddUser)

	// Authenticated user endpoints.
	user := web.New()
	user.Use(middleware.UserAuthorizer)
	user.Get("/api/user", handler.GetCurrentUser)
	user.Patch("/api/user", handler.UpdateCurrentUser)
	user.Put("/api/user", handler.UpdateCurrentUser)
	user.Get("/api/user/documents", handler.GetCurrentUserDocuments)
	mux.Handle("/api/user", user)
	mux.Handle("/api/user/documents", user)

	// Document endpoints.
	doc := web.New()
	doc.Use(middleware.UserAuthorizer)

	// @todo uncommented until middleware can captures parsed URL Params.
	// https://github.com/zenazn/goji/issues/76
	//
	// doc.Use(middleware.DocumentToContextInjector)

	doc.Get("/api/documents", handler.GetAllDocuments)
	doc.Post("/api/documents", handler.AddDocument)

	doc.Get("/api/documents/:docId", handler.GetDocumentById)
	doc.Delete("/api/documents/:docId", handler.DeleteDocument)

	// Document files.
	doc.Get("/api/documents/:docId/files", handler.GetDocumentFiles)
	doc.Post("/api/documents/:docId/files", handler.AddDocumentFile)
	doc.Get("/api/documents/:docId/files/sid", handler.GetDocumentFilesSid)
	doc.Delete("/api/documents/:docId/files/:fileid", handler.DeleteDocumentFile)

	mux.Handle("/api/documents", doc)
	mux.Handle("/api/documents/*", doc)

	// Dev endpoints. Provide helper handlers during development.
	dev := web.New()
	dev.Use(middleware.DevEnv)
	dev.Get("/api/dev/seed_users", handler.DevSeedUsers)
	dev.Get("/api/dev/seed_documents", handler.DevSeedDocuments)
	mux.Handle("/api/dev/*", dev)

	return mux
}
