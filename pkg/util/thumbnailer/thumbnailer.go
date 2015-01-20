package thumbnailer

import (
	"fmt"
	"os/exec"
)

type Thumbnail struct {
	Filepath  string
	ImageType string
	Width     int
	Height    int
	Size      int64
}

type IdentifiedImage struct {
	Size   int64
	Width  int
	Height int
}

type Thumbnailer interface {
	Create(srcPath, dst string, w, h int) (*Thumbnail, error)
}

func VipsCreate(src, dst string, w, h int) (*Thumbnail, error) {
	vt := newVipsthumbnail()
	return vt.Create(src, dst, w, h)
}

func IdentifyImage(fpath string) (*IdentifiedImage, error) {
	cmd := exec.Command("identify", "-format", `"%w:%h:%b"`, fpath)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var w, h int
	var s int64

	fmt.Sscanf(string(out), `"%d:%d:%dB"`, w, h, s)

	return &IdentifiedImage{s, w, h}, nil
}
