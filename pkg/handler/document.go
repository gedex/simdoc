package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gedex/simdoc/pkg/datastore"
	"github.com/gedex/simdoc/pkg/middleware"
	"github.com/gedex/simdoc/pkg/model"
	"github.com/gedex/simdoc/pkg/util/upload"
	"github.com/gedex/simdoc/pkg/util/upload/processor"

	"code.google.com/p/go-uuid/uuid"

	"github.com/goji/context"
	"github.com/zenazn/goji/web"
)

// GetAllDocuments accepts a request to retrieve all docuemnts from the datastore and
// returns in JSON format.
//
// GET /api/documents
//
func GetAllDocuments(c web.C, w http.ResponseWriter, r *http.Request) {
	var ctx = context.FromC(c)

	docs, err := datastore.GetAllDocuments(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if docs == nil {
		w.Write([]byte(`[]`))
	} else {
		json.NewEncoder(w).Encode(docs)
	}
}

// GetDocumentById accepts a request to retrieve a document, for given document
// ID docId, from the datastore and returns in JSON format.
//
// GET /api/documents/:docId
//
func GetDocumentById(c web.C, w http.ResponseWriter, r *http.Request) {
	// @todo remove me once DocumentToContextInjector is being used.
	if ok := docToContext(&c, w); !ok {
		return
	}

	var doc = ToDocument(c)
	if doc == nil {
		respWithError(w, http.StatusNotFound, ErrorNotFound)
		return
	}

	json.NewEncoder(w).Encode(doc)
}

// AddDocument accepts a request to add new document into the datastore.
//
// POST /api/documents
//
func AddDocument(c web.C, w http.ResponseWriter, r *http.Request) {
	doc, err := parseSubmittedDoc(r)
	if err != nil {
		respWithError(w, http.StatusBadRequest, ErrorInvalidJSONRequest)
		return
	}

	var usr = ToUser(c)
	if usr == nil {
		respWithError(w, http.StatusUnauthorized, ErrorRequireAuthentication)
		return
	}

	// Sets creator to current user.
	doc.CreatedBy = usr.ID

	// Sets default status if not provided.
	if doc.Status == "" {
		doc.Status = model.DefaultDocumentStatus
	}

	// Validate the model.
	if ve := model.Validate(doc); ve != nil {
		respWithError(w, http.StatusBadRequest, ErrorValidationFailed, getValidationErrors("documents", ve)...)
		return
	}

	err = datastore.AddDocument(context.FromC(c), doc)
	if err != nil {
		respWithError(w, http.StatusBadRequest, ErrorValidationFailed)
		return
	}

	// @todo send notification to registered listeners.

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(doc)
}

// DeleteDocument accepts a request to delete a document specified by docId in
// the URL.
//
// DELETE /api/documents/:docId
//
func DeleteDocument(c web.C, w http.ResponseWriter, r *http.Request) {
	// @todo remove me once DocumentToContextInjector is being used.
	if ok := docToContext(&c, w); !ok {
		return
	}

	var doc = ToDocument(c)
	if doc == nil {
		respWithError(w, http.StatusNotFound, ErrorNotFound)
		return
	}

	var usr = ToUser(c)
	if usr == nil {
		respWithError(w, http.StatusUnauthorized, ErrorRequireAuthentication)
		return
	}

	var err error

	// Check if current user has priviledge to delete the document.
	if doc.CreatedBy == usr.ID || usr.IsAdmin() {
		err = datastore.DeleteDocument(context.FromC(c), doc.ID)
	} else {
		respWithError(w, http.StatusForbidden, ErrorForbidden)
		return
	}

	if err != nil {
		respWithError(w, http.StatusInternalServerError, ErrorInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// @todo implement me
func GetDocumentFiles(c web.C, w http.ResponseWriter, r *http.Request) {
	respWithError(w, http.StatusNotImplemented, ErrorNotImplemented)
}

// @todo refactor me!
//
// POST /api/documents/:docId/files
//
func AddDocumentFile(c web.C, w http.ResponseWriter, r *http.Request) {
	fsRoot := c.Env["fsRoot"].(string)
	prefix := c.Env["filesPrefix"].(string)

	files, err := upload.FromHttp(r, fsRoot)
	switch {
	case err != nil && err == upload.ErrorIncomplete:
		w.WriteHeader(http.StatusOK)

		chunkResp := make(map[string]interface{}, 1)
		chunkResp["file"] = struct {
			Sid  string `json:"sid"`
			Size int64  `json:"size"`
		}{files[0].Sid, files[0].Size}
		chunkResp["status"] = "ok"

		json.NewEncoder(w).Encode(chunkResp)
		return
	case err != nil && err != upload.ErrorIncomplete:
		w.WriteHeader(http.StatusBadRequest)

		// If files is nil, error happens in early stage (during the upload).
		if files == nil {
			respWithError(w, http.StatusBadRequest, err)
		} else {
			fe := getUploadFilesErrors(files)

			wrapFe := make(map[string]interface{}, 1)
			wrapFe["files"] = fe

			json.NewEncoder(w).Encode(wrapFe)
		}
		return
	}

	// Callback function that's called after each processor.Process.
	afterFn := func(out *upload.File, err error) (*upload.File, error) {
		if err != nil || out == nil {
			return out, err
		}
		if out.Filepath == "" {
			return out, err
		}
		if rel, relErr := filepath.Rel(fsRoot, out.Filepath); relErr == nil {
			out.URL = filepath.Join(prefix, rel)
		}
		return out, err
	}

	resp := make([]*upload.FileResult, 0)
	for _, f := range files {
		var fr *upload.FileResult

		bpath, err := upload.CreateDir(fsRoot, f.Type)
		if err != nil {
			f.Error = err
			fr = &upload.FileResult{f, nil}
		} else {
			fr = upload.ProcessFile(
				f,
				upload.AfterProcessFn(afterFn),
				processor.Mover("default", upload.SourceOriginal, filepath.Join(bpath.Abs(), generateFilename(f, "default"))),
				processor.Resizer("thumbnail-120x90", "default", filepath.Join(bpath.Abs(), generateFilename(f, "thumbnail-120x90")), 120, 90),
			)
		}

		resp = append(resp, fr)
	}

	// Wrap resp so that JS uploader can consumes the response.
	wrapResp := struct {
		Files []*upload.FileResult `json:"files"`
	}{resp}

	// @todo write to DB.

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wrapResp)
}

// GetDocumentFilesSid
//
// GET /api/documents/:docId/files/sid
//
func GetDocumentFilesSid(c web.C, w http.ResponseWriter, r *http.Request) {
	sid := uuid.New()

	resp := struct {
		Sid string `json:"sid"`
	}{sid}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func DeleteDocumentFile(c web.C, w http.ResponseWriter, r *http.Request) {
	respWithError(w, http.StatusNotImplemented, ErrorNotImplemented)
}

// parseSubmittedDoc returns Document through POST or PUT.
func parseSubmittedDoc(r *http.Request) (*model.Document, error) {
	decoder := json.NewDecoder(r.Body)

	doc := new(model.Document)
	err := decoder.Decode(doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// @todo remove me once DocumentToContextInjector is being used.
func docToContext(c *web.C, w http.ResponseWriter) bool {
	docIdStr := c.URLParams["docId"]
	docId, _ := strconv.ParseInt(docIdStr, 10, 64)

	if docId > 0 {
		return checkDoc(c, w, docId)
	}

	respWithError(w, http.StatusNotFound, ErrorNotFound)

	return false
}

// @todo remove me once DocumentToContextInjector is being used.
func checkDoc(c *web.C, w http.ResponseWriter, docId int64) bool {
	ctx := context.FromC(*c)
	user := middleware.ToUser(c)

	// @todo check document participants?

	doc, err := datastore.GetDocumentById(ctx, docId)
	switch {
	case err != nil && user == nil:
		respWithError(w, http.StatusUnauthorized, ErrorRequireAuthentication)
		return false
	case err != nil && user != nil:
		respWithError(w, http.StatusNotFound, ErrorNotFound)
		return false
	}

	middleware.DocToC(c, doc)

	return true
}

func generateFilename(f *upload.File, suffix string) string {
	t := time.Now()
	ts := int64(t.Hour()*3600 + t.Minute()*60 + t.Second())

	if suffix != "" {
		suffix = "-" + suffix
	}

	return fmt.Sprintf("%x_%s%s%s", md5.Sum([]byte(f.Sid)), strconv.FormatInt(ts, 36), suffix, filepath.Ext(f.Name))
}
