package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/gedex/simdoc/pkg/datastore"
	"github.com/gedex/simdoc/pkg/datastore/database"
	"github.com/gedex/simdoc/pkg/handler"
	"github.com/gedex/simdoc/pkg/middleware"
	"github.com/gedex/simdoc/pkg/router"

	"code.google.com/p/go.net/context"
	webcontext "github.com/goji/context"
	"github.com/zenazn/goji/web"
	gojiMiddleware "github.com/zenazn/goji/web/middleware"
)

var (
	httpServerPort = flag.String("port", ":8080", "HTTP server port")
	dsn            = flag.String("dsn", "root:root@tcp(192.168.42.43:3306)/simdoc", "DSN")

	// Salt for password.
	passwdSalt = flag.String("pass_salt", "s1md0c!#", "Salt for password. Default to 's1md0c!#'")

	// JWT secret.
	jwtSecret = flag.String("jwt_secret", "s1md0c-jwt-secret", "JWT Secret")

	// Env (dev, staging or prod). Default to 'prod'.
	env = flag.String("env", "prod", "Env name. Use 'dev' for additional handler during development. Default to 'prod'")

	// URL Prefix to access uploaded files
	filesPrefix = flag.String("files_prefix", "/files/", "URL Prefix to access uploaded files. Default to '/files/'")

	// fsRoot is a root path to store files in file system.
	fsRoot = flag.String("fs_root", "/tmp/simdoc/files", "Filestore root. Default to '/tmp/simdoc/files'")

	// DB as Datastore
	db *sql.DB
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: simdoc [flags]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Usage = usage
	flag.Parse()

	// DB.
	db = database.MustConnect(*dsn)

	// Static resources for SPA.
	// @todo

	// API Routers.
	r := router.New()

	// Middleware for logging.
	r.Use(gojiMiddleware.Logger)

	// Middleware for OPTIONS request.
	r.Use(middleware.OptionsRequestHeaderInjector)

	// Middleware that inject services into context.
	r.Use(ContextMiddleware)

	// Middleware that inject header fields.
	r.Use(middleware.HeaderInjector)

	// Middleware that may injects user information into context.
	r.Use(middleware.UserToContextInjector)

	// Handles all /api/* requests with API routers.
	http.Handle("/api/", r)

	// Handle GET uploaded files requests with static file server.
	http.Handle(*filesPrefix, handler.NewFileServer(*fsRoot, *filesPrefix))

	// Starts HTTP server.
	// @todo supports HTTPS.
	panic(http.ListenAndServe(*httpServerPort, nil))
}

// ContextMiddleware creates a new go.net/context and injects into the current
// goji context.
func ContextMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ctx = context.Background()
		ctx = datastore.NewContext(ctx, database.NewDatastore(db))

		webcontext.Set(c, ctx)
		c.Env["passwdSalt"] = *passwdSalt
		c.Env["jwtSecret"] = *jwtSecret
		c.Env["env"] = *env
		c.Env["fsRoot"] = *fsRoot
		c.Env["filesPrefix"] = *filesPrefix

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
