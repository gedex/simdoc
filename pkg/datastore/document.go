package datastore

import (
	"code.google.com/p/go.net/context"
	"github.com/gedex/simdoc/pkg/model"
)

type Documentstore interface {
	// GetDocumentById retrieves a document from the datastore for the given docId.
	GetDocumentById(docId int64) (*model.Document, error)

	// GetAllDocuments retrieves a list of all documents from the datastore.
	GetAllDocuments() ([]*model.Document, error)

	// AddDocuments adds a document into the datastore.
	AddDocument(doc *model.Document) error

	// UpdateDocument updates a document in the datastore.
	UpdateDocument(doc *model.Document) error

	// DeleteDocument deletes a document, for the given docId, in the datastore.
	DeleteDocument(docId int64) error

	// GetAllDocumentFiles retrieves a list of all files of a document, for the
	// given docId, from the datastore.
	GetAllDocumentFiles(docId int64) ([]*model.DocumentFile, error)

	// AddDocumentFile adds a file to a document in the datastore.
	AddDocumentFile(f *model.DocumentFile) error

	// DeleteDocumentFile deletes a file, for the given fileId, in the datastore.
	DeleteDocumentFile(fileId int64) error

	// DeleteDocumentFIles delete all files in a document, for the given docId,
	// in the datastore.
	DeleteDocumentFiles(docId int64) error
}

// GetDocumentById retrieves a document from the datastore for the given docId.
func GetDocumentById(c context.Context, docId int64) (*model.Document, error) {
	return FromContext(c).GetDocumentById(docId)
}

// GetAllDocuments retrieves a list of all documents from the datastore.
func GetAllDocuments(c context.Context) ([]*model.Document, error) {
	return FromContext(c).GetAllDocuments()
}

// AddDocuments adds a document into the datastore.
func AddDocument(c context.Context, doc *model.Document) error {
	return FromContext(c).AddDocument(doc)
}

// UpdateDocument updates a document in the datastore.
func UpdateDocument(c context.Context, doc *model.Document) error {
	return FromContext(c).UpdateDocument(doc)
}

// DeleteDocument deletes a document, for the given docId, in the datastore.
func DeleteDocument(c context.Context, docId int64) error {
	return FromContext(c).DeleteDocument(docId)
}

// GetAllDocumentFiles retrieves a list of all files of a document, for the
// given docId, from the datastore.
func GetAllDocumentFiles(c context.Context, docId int64) ([]*model.DocumentFile, error) {
	return FromContext(c).GetAllDocumentFiles(docId)
}

// AddDocumentFile adds a file to a document, for the given docId, in the datastore.
func AddDocumentFile(c context.Context, f *model.DocumentFile) error {
	return FromContext(c).AddDocumentFile(f)
}

// DeleteDocumentFile deletes a file, for the given fileId, in the datastore.
func DeleteDocumentFile(c context.Context, fileId int64) error {
	return FromContext(c).DeleteDocumentFile(fileId)
}

// DeleteDocumentFIles delete all files in a document, for the given docId,
// in the datastore.
func DeleteDocumentFiles(c context.Context, docId int64) error {
	return FromContext(c).DeleteDocumentFiles(docId)
}
