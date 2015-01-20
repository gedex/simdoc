package thumbnailer

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/gedex/simdoc/pkg/util/mimetype"
)

type vipsCmd struct {
	cmd string
}

func newVipsthumbnail() *vipsCmd {
	return &vipsCmd{"vipthumbnail"}
}

func (v vipsCmd) Create(src, dst string, w, h int) (*Thumbnail, error) {
	if src == dst {
		return v.ret(src)
	}

	ddir, _ := filepath.Split(dst)

	// If only filename is given in destination dst, use base dir from source src.
	if ddir == "" {
		sdir, _ := filepath.Split(src)

		dst = filepath.Join(sdir, dst)
	}

	// Creates the thumbnail.
	args := []string{
		"-s",
		fmt.Sprintf("%dx%d", w, h),
		"-o",
		dst,
		src,
	}
	fmt.Println(args)
	cmd := exec.Command("vipsthumbnail", args...)

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return v.ret(dst)
}

func (v vipsCmd) ret(fpath string) (*Thumbnail, error) {
	mt, err := mimetype.FromFilepath(fpath)
	if err != nil {
		return nil, err
	}

	id, err := IdentifyImage(fpath)
	if err != nil {
		return nil, err
	}

	thumb := &Thumbnail{
		Filepath:  fpath,
		ImageType: mimetype.Base(mt),
		Width:     id.Width,
		Height:    id.Height,
		Size:      id.Size,
	}

	return thumb, nil
}
