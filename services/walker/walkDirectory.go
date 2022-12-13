package walker

import (
	"io/fs"
	"path/filepath"
)

func Walk(path string) (error, []string) {
	files := make([]string, 0)
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		files = append(files, path+d.Name())
		return nil
	})
	if err != nil {
		return err, nil
	}
	return nil, files
}
