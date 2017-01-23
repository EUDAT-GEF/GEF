package def

import (
	"os"
	"path/filepath"

	"github.com/pborman/uuid"
)

const (
	tmpDirDefault = "gefdocker"
	tmpDirPerm    = 0700
)

func NewRandomTmpDir(parentComponents ...string) (path string, id string, err error) {
	id = uuid.New()
	path = filepath.Join(append(parentComponents, id)...)
	err = os.MkdirAll(path, os.FileMode(tmpDirPerm))
	if err != nil {
		return "", "", Err(err, "Cannot create temporary dir %s", path)
	}
	return path, id, nil
}

func MakeTmpDir(tmpDir string) (string, error) {
	if tmpDir == "" {
		tmpDir = tmpDirDefault
	}
	if !filepath.IsAbs(tmpDir) {
		tmpDir = filepath.Join(os.TempDir(), tmpDir)
	}
	if err := os.MkdirAll(tmpDir, os.FileMode(tmpDirPerm)); err != nil {
		return tmpDir, Err(err, "cannot create temporary directory %s", tmpDir)
	}
	return tmpDir, nil
}
