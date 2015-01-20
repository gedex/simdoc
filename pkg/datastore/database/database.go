package database

import (
	"database/sql"

	"github.com/BurntSushi/migration"
	"github.com/gedex/simdoc/pkg/datastore"
	"github.com/gedex/simdoc/pkg/datastore/migrate"
	_ "github.com/go-sql-driver/mysql"
	"github.com/russross/meddler"
)

const (
	// @todo use me
	ROWS_LIMIT_DEFAULT = 20
	ROWS_LIMIT_MAX     = 100

	// @todo use me
	ASC  = "ASC"
	DESC = "DESC"
)

// @todo use me
type ORDER_BY int

// @todo use me
const (
	USER_ORDER_BY_ID ORDER_BY = iota
	USER_ORDER_BY_LOGIN
	USER_ORDER_BY_EMAIL
	USER_ORDER_BY_NAME
	USER_ORDER_BY_CREATED

	DOC_ORDER_BY_ID
	DOC_ORDER_BY_NAME
	DOC_ORDER_BY_CREATED
)

func MustConnect(dsn string) *sql.DB {
	meddler.Default = meddler.MySQL

	migration.DefaultGetVersion = migrate.GetVersion
	migration.DefaultSetVersion = migrate.SetVersion

	var migrations = []migration.Migrator{
		migrate.Setup,
	}

	db, err := migration.Open("mysql", dsn, migrations)
	if err != nil {
		panic(err)
	}
	return db
}

func NewDatastore(db *sql.DB) datastore.Datastore {
	return struct {
		*Userstore
		*Documentstore
	}{
		NewUserstore(db),
		NewDocumentstore(db),
	}
}
