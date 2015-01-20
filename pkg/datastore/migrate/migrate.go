package migrate

import (
	"github.com/BurntSushi/migration"
)

func Setup(tx migration.LimitedTx) error {
	var createTablesCmds = []string{
		userTable,
		documentTable,
		documentFilesTable,
		documentParticipantsTable,
	}

	for _, cmd := range createTablesCmds {
		_, err := tx.Exec(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

var userTable = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	login VARCHAR(255),
	password VARCHAR(255),
	email VARCHAR(255),
	role VARCHAR(255),
	name VARCHAR(255),
	created INTEGER,
	updated INTEGER,
	UNIQUE(login),
	UNIQUE(email)
)
`

var documentTable = `
CREATE TABLE IF NOT EXISTS documents (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	status VARCHAR(255),
	name VARCHAR(255),
	created_by INTEGER,
	created INTEGER,
	updated INTEGER
)
`
var documentFilesTable = `
CREATE TABLE IF NOT EXISTS document_files (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	document_id INTEGER,
	user_id INTEGER,
	name VARCHAR(255),
	path TEXT,
	url TEXT,
	meta TEXT,
	versions TEXT,
	created INTEGER,
	updated INTEGER,
	UNIQUE(name)
)
`

var documentParticipantsTable = `
CREATE TABLE IF NOT EXISTS document_participants (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	document_id INTEGER,
	user_id INTEGER
)
`
