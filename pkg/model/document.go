package model

var (
	DocumentStatusDraft     = "draft"
	DocumentStatusPublished = "published"
	DefaultDocumentStatus   = DocumentStatusDraft
)

// Document represents document. It has one-to-many relationship with DocumentFile
// and one-to-many DocumentParticipant.
type Document struct {
	ID        int64  `meddler:"id,pk"       json:"id"`
	Status    string `meddler:"status"      validate:"doc_status" json:"status"`
	Name      string `meddler:"name"        validate:"nonzero" json:"name"`
	CreatedBy int64  `meddler:"created_by"  json:"created_by"`
	Created   int64  `meddler:"created"     json:"created_at"`
	Updated   int64  `meddler:"updated"     json:"updated_at"`
}

// DocumentParticipant represents participant in a document.
type DocumentParticipant struct {
	ID         int64 `meddler:"id,pk" json:"id"`
	DocumentID int64 `meddler:"document_id"`
	UserID     int64 `meddler:"user_id"`
}

// DocumentFile represents attached file in a document.
type DocumentFile struct {
	ID         int64                           `meddler:"id,pk"         json:"id"`
	DocumentID int64                           `meddler:"document_id"   validate:"nonzero" json:"document_id"` // Associated document
	Name       string                          `meddler:"name"          json:"name"`                           // Filename of file being uploaded
	Filepath   string                          `meddler:"filepath"      json:"-"`                              // Location of this file in document store
	URL        string                          `meddler:"url"           json:"url"`                            // URL, without domain
	Meta       *DocumentFileMeta               `meddler:"meta,json"     json:"meta"`                           // meta of uploaded file
	Versions   map[string]*DocumentFileVersion `meddler:"versions,json" json:"versions"`                       // Key is processor name, for instance "thumbnail-150x90"
	Created    int64                           `meddler:"created"       json:"created_at"`
	Updated    int64                           `meddler:"updated"       json:"updated_at"`
}

type DocumentFileMeta struct {
	Type string // Base mime type, for instance "image" for "image/jpeg"
	Mime string // For instance "image/jpeg"
	Size int64  // File size in bytes
}

type DocumentFileVersion struct {
	Filepath string            `json:"filepath"` // Location of this file version in document store
	URL      string            `json:"url"`      // URL, without domain
	Meta     *DocumentFileMeta `json:"meta"`     // Meta of this file version
}
