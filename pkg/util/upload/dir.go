package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type UploadDir struct {
	root string
	path string
}

func CreateDir(root, baseMime string) (*UploadDir, error) {
	ud := &UploadDir{root, baseMime}
	ud.preparePath(baseMime)

	if err := ud.create(); err != nil {
		return nil, err
	}

	return ud, nil
}

// abs returns absolute path for upload dir.
func (ud *UploadDir) Abs() string {
	return filepath.Join(ud.root, ud.path)
}

// create creates a directory for the upload.
func (ud *UploadDir) create() error {
	return os.MkdirAll(ud.Abs(), 0755)
}

func (ud *UploadDir) preparePath(baseMime string) {
	now := time.Now()
	ud.path = fmt.Sprintf("/%s/%d/%d/%d", baseMime, now.Year(), now.Month(), now.Day())
}
