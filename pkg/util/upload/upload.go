package upload

import (
	"errors"
	"net/http"
)

type File struct {
	Name     string `json:"name"`
	Sid      string `json:"sid,omitempty"`
	Mime     string `json:"mime,omitempty"`
	Type     string `json:"type,omitempty"`
	Filepath string `json:"-"`
	URL      string `json:"url,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Error    error  `json:"error,omitempty"`
}

var ErrorIncomplete = errors.New("Incomplete")

type Uploader interface {
	Upload(dirPath string) ([]*File, error)
}

func FromHttp(r *http.Request, dirPath string) ([]*File, error) {
	u := newHttpUploader(r)

	return u.Upload(dirPath)
}
