package processor

import (
	"errors"

	"github.com/gedex/simdoc/pkg/util/thumbnailer"
	"github.com/gedex/simdoc/pkg/util/upload"
)

type resizer struct {
	name string
	src  string
	dst  string // filepath destination
	w    int
	h    int
}

func Resizer(name, src, dst string, w, h int) upload.Processor {
	return &resizer{name, src, dst, w, h}
}

func (r *resizer) Process(src *upload.File) (*upload.File, error) {
	t, err := thumbnailer.VipsCreate(src.Filepath, r.dst, r.w, r.h)
	if err != nil {
		return nil, errors.New("thumbnailer.VipsCreate returns error: " + err.Error())
	}

	out := *src
	out.Filepath = t.Filepath
	out.Size = t.Size

	return &out, nil
}

func (r *resizer) GetName() string {
	return r.name
}

func (r *resizer) GetSource() string {
	return r.src
}

func (r *resizer) CanProcess(baseMime string) bool {
	if baseMime == "image" {
		return true
	}
	return false
}
