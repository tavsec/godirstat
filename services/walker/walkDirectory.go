package walker

import (
	"io/fs"
	"path/filepath"
)

func Walk(path string) (error, []fs.FileInfo) {
	files := make([]fs.FileInfo, 0)
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		fileInfo, err := d.Info()
		files = append(files, fileInfo)
		return err
	})
	if err != nil {
		return err, nil
	}
	return nil, files
}
