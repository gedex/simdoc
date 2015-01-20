package handler

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type fileServer struct {
	rootPath  string
	urlPrefix string
}

func NewFileServer(rootPath, urlPrefix string) *fileServer {
	return &fileServer{rootPath, urlPrefix}
}

func (f *fileServer) absPath(fp string) string {
	if !strings.HasPrefix(fp, "/") {
		fp = "/" + fp
	}

	return filepath.Join(f.rootPath, path.Clean(fp))
}

func (f *fileServer) serveFile(w http.ResponseWriter, r *http.Request, fp string) {
	fo, err := os.Open(fp)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer fo.Close()

	s, err1 := fo.Stat()
	if err1 != nil {
		http.NotFound(w, r)
		return
	}
	if s.IsDir() {
		http.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, fp, s.ModTime(), fo)
}

func (f *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(f.urlPrefix, "/") {
		f.urlPrefix = "/" + f.urlPrefix
	}

	rel, _ := filepath.Rel(f.urlPrefix, r.URL.Path)

	f.serveFile(w, r, f.absPath(rel))
}
