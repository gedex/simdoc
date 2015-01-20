package processor

import (
	"os"

	"github.com/gedex/simdoc/pkg/util/upload"
)

type mover struct {
	name string
	src  string
	dst  string // filepath destination
}

func Mover(name, src, dst string) upload.Processor {
	return &mover{name, src, dst}
}

func (r *mover) Process(src *upload.File) (*upload.File, error) {
	err := os.Rename(src.Filepath, r.dst)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(r.dst)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	dst := *src
	dst.Filepath = r.dst
	dst.Size = fi.Size()

	return &dst, nil
}

func (r *mover) GetName() string {
	return r.name
}

func (r *mover) GetSource() string {
	return r.src
}

func (r *mover) CanProcess(baseMime string) bool {
	return true
}
