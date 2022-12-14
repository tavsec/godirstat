package walker

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

var WG sync.WaitGroup
var Files = make([]fs.FileInfo, 0)

func WalkDir(dir string) []fs.FileInfo {
	defer WG.Done()

	visit := func(path string, f os.FileInfo, err error) error {
		Files = append(Files, f)
		if f.IsDir() && path != dir {
			WG.Add(1)
			go WalkDir(path)
			return filepath.SkipDir
		}
		return nil
	}

	err := filepath.Walk(dir, visit)
	if err != nil {
		return nil
	}
	return nil
}
