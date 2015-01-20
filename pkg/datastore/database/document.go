package database

import (
	"log"
	"time"

	"github.com/gedex/simdoc/pkg/model"
	"github.com/russross/meddler"
)

type Documentstore struct {
	meddler.DB
}

type DocStatus bool

func init() {
	meddler.Register("doc_status", DocStatus(false))
}

func NewDocumentstore(db meddler.DB) *Documentstore {
	return &Documentstore{db}
}

func (db *Documentstore) GetDocumentById(docId int64) (*model.Document, error) {
	var doc = new(model.Document)
	var err = meddler.Load(db, docTable, doc, docId)

	return doc, err
}

func (db *Documentstore) GetAllDocuments() ([]*model.Document, error) {
	var docs []*model.Document
	var err = meddler.QueryAll(db, &docs, docListQuery)

	return docs, err
}

func (db *Documentstore) AddDocument(doc *model.Document) error {
	if doc.Created == 0 {
		doc.Created = time.Now().UTC().Unix()
	}
	doc.Updated = time.Now().UTC().Unix()

	return meddler.Save(db, docTable, doc)
}

func (db *Documentstore) UpdateDocument(doc *model.Document) error {
	doc.Updated = time.Now().UTC().Unix()

	return meddler.Save(db, docTable, doc)
}

func (db *Documentstore) DeleteDocument(docId int64) error {
	var _, err = db.Exec(docDeleteQuery, docId)

	return err
}

func (db *Documentstore) GetAllDocumentFiles(docId int64) ([]*model.DocumentFile, error) {
	var files []*model.DocumentFile
	var err = meddler.QueryAll(db, &files, docFilesListQuery, docId)

	return files, err
}

func (db *Documentstore) AddDocumentFile(f *model.DocumentFile) error {
	if f.Created == 0 {
		f.Created = time.Now().UTC().Unix()
	}
	f.Updated = time.Now().UTC().Unix()

	return meddler.Save(db, docFilesTable, f)
}

func (db *Documentstore) DeleteDocumentFile(fileId int64) error {
	var _, err = db.Exec(docFileDeleteQuery, fileId)

	return err
}

func (db *Documentstore) DeleteDocumentFiles(docId int64) error {
	var _, err = db.Exec(docFilesDeleteQuery, docId)

	return err
}

const docTable = "documents"

const docListQuery = `
SELECT * FROM documents
ORDER BY name
`

const docDeleteQuery = `
DELETE FROM documents
WHERE id=?
`

const docFilesTable = "document_files"

const docFilesListQuery = `
SELECT * FROM document_files
WHERE document_id=?
ORDER BY created
`

const docFileDeleteQuery = `
DELETE FROM document_files
WHERE id=?
`

const docFilesDeleteQuery = `
DELETE FROM document_files
WHERE document_id=?
`

func (ds DocStatus) PreRead(fieldAddr interface{}) (scanTarget interface{}, err error) {
	log.Printf("%+v\n", fieldAddr)
	return fieldAddr, nil
}

func (ds DocStatus) PostRead(fieldAddr, scanTarget interface{}) error {
	log.Printf("%+v\n", fieldAddr)
	log.Printf("%+v\n", scanTarget)
	return nil
}

func (ds DocStatus) PreWrite(field interface{}) (saveValue interface{}, err error) {
	log.Printf("%+v\n", field)
	log.Printf("%+v\n", saveValue)
	return field, nil
}
