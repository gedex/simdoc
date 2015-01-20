package upload

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gedex/simdoc/pkg/util/mimetype"
)

type HttpUploader struct {
	*meta
	*body
	*http.Request
}

type meta struct {
	filename  string
	mediaType string
	boundary  string
	cr        *contentRange
	sid       string // @todo move me to middleware?
}

type contentRange struct {
	start int64
	end   int64
	size  int64
}

type body struct {
	content io.ReadCloser
	closed  bool
	reader  *contentReader
}

type contentReader struct {
	r        io.Reader
	mr       *multipart.Reader
	filename string
}

func newHttpUploader(r *http.Request) Uploader {
	return &HttpUploader{&meta{}, &body{}, r}
}

func (u *HttpUploader) Upload(dirPath string) ([]*File, error) {
	if err := u.parseRequest(); err != nil {
		return nil, err
	}
	defer u.content.Close()

	files, err := u.saveFiles(dirPath)
	if err == ErrorIncomplete {
		return files, err
	}
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (u *HttpUploader) parseRequest() error {
	if err := u.parseContentType(); err != nil {
		return err
	}
	if err := u.parseContentRange(); err != nil {
		return err
	}
	if err := u.parseContentDisposition(); err != nil {
		return err
	}
	if err := u.parseSessionID(); err != nil {
		return err
	}
	if err := u.parseBody(); err != nil {
		return err
	}
	return nil
}

func (u *HttpUploader) parseContentType() error {
	ct := u.Header.Get("Content-Type")
	if ct == "" {
		u.mediaType = "application/octet-stream"
		return nil
	}

	mt, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return err
	}

	if mt == "multipart/form-data" {
		boundary, ok := params["boundary"]
		if !ok {
			return errors.New("meta: boundary not defined")
		}

		u.mediaType = mt
		u.boundary = boundary
	} else {
		u.mediaType = "application/octet-stream"
	}

	return nil
}

func (u *HttpUploader) parseContentRange() error {
	cr := u.Header.Get("Content-Range")
	if cr == "" {
		return nil
	}

	var start, end, size int64

	_, err := fmt.Sscanf(cr, "bytes %d-%d/%d", &start, &end, &size)
	if err != nil {
		return err
	}

	u.cr = &contentRange{start, end, size}

	return nil
}

func (u *HttpUploader) parseContentDisposition() error {
	cd := u.Header.Get("Content-Disposition")
	if cd == "" {
		return nil
	}

	_, params, err := mime.ParseMediaType(cd)
	if err != nil {
		return err
	}

	filename, ok := params["filename"]
	if !ok {
		return errors.New("meta: filename in Content-Disposition is not defined")
	}
	u.filename = filename

	return nil
}

// @todo move this into middleware?
func (u *HttpUploader) parseSessionID() error {
	sid := u.URL.Query().Get("sid")
	if sid == "" {
		return errors.New("meta: missing sid in query param")
	}
	u.sid = sid

	return nil
}

func (u *HttpUploader) parseBody() error {
	xfile := u.Header.Get("X-File")
	if xfile == "" {
		u.body.content = u.Body
		return nil
	}

	fh, err := os.Open(xfile)
	if err != nil {
		return err
	}
	u.body.content = fh

	return nil
}

func (u *HttpUploader) saveFiles(dirPath string) ([]*File, error) {
	files := make([]*File, 0)
	for {
		f, err := u.saveFile(dirPath)
		if err == io.EOF {
			break
		}

		if err == ErrorIncomplete {
			files = append(files, f)
			return files, err
		}

		if err != nil {
			return nil, err
		}

		files = append(files, f)
	}

	return files, nil
}

func (u *HttpUploader) saveFile(dirPath string) (*File, error) {
	r, err := u.getReader()
	if err != nil {
		return nil, err
	}

	f, err := r.writeTo(dirPath, u.meta)
	if err != nil {
		return nil, err
	}

	if u.cr != nil && f.Size != u.cr.size {
		return f, ErrorIncomplete
	}

	f.Mime, err = mimetype.FromFilepath(f.Filepath)
	if err != nil {
		return nil, err
	}
	f.Type = mimetype.Base(f.Mime)

	return f, nil
}

func (u *HttpUploader) getReader() (*contentReader, error) {
	if u.reader == nil {
		u.reader = &contentReader{}
	}
	if u.mediaType == "multipart/form-data" {
		if u.reader.mr == nil {
			u.reader.mr = multipart.NewReader(u.content, u.boundary)
		}
		for {
			part, err := u.reader.mr.NextPart()
			if err != nil {
				return nil, err
			}
			if part.FormName() == "files[]" {
				u.reader.r = part
				u.reader.filename = part.FileName()

				return u.reader, nil
			}
		}
	}

	if u.body.closed {
		return nil, io.EOF
	}
	u.body.closed = true

	u.reader.r = u.body.content
	u.reader.filename = u.filename

	return u.reader, nil
}

func (r *contentReader) writeTo(dirPath string, m *meta) (*File, error) {
	var f *os.File
	var err error
	if m.cr == nil {
		f, err = ioutil.TempFile(os.TempDir(), "simdoc")
		_, err = io.Copy(f, r.r)
	} else {
		f, err = r.getFileChunk(dirPath, m)
		_, err = io.CopyN(f, r.r, m.cr.end-m.cr.start+1)
	}
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	file := &File{
		Name:     r.filename,
		Sid:      m.sid,
		Filepath: f.Name(),
		Size:     fi.Size(),
	}

	return file, nil
}

func (r *contentReader) getFileChunk(dirPath string, m *meta) (*os.File, error) {
	path := filepath.Join(dirPath, "chunks")
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	hasher := md5.New()
	hasher.Write([]byte(m.sid + m.filename))

	fname := hex.EncodeToString(hasher.Sum(nil))
	fpath := filepath.Join(path, fname)

	file, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return nil, err
	}

	if _, err = file.Seek(m.cr.start, 0); err != nil {
		return nil, err
	}

	return file, nil
}
