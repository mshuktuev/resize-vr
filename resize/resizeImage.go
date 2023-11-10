package resizeImage

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type ImageOptions struct {
	Width int
	Height int
	Quality int
}

type ResizeProgress struct {
	Progress *progressbar.ProgressBar
	Mutex    *sync.Mutex
	Wg       *sync.WaitGroup
}

func (rp *ResizeProgress) Increment() {
	rp.Mutex.Lock()
	defer rp.Mutex.Unlock()
	rp.Progress.Add(1)
}

func ResizeImage(path , outDir, splitPath string,  options ImageOptions, resizeInfo *ResizeProgress) error {
	dir := subpath(splitPath, path)
	outPath := outDir + string(filepath.Separator) + dir
	fmt.Println(outPath)

	return nil
}

func subpath(homeDir, prevDir string) string {
	subFiles := ""
	for {
			dir, file := filepath.Split(prevDir)
			if file == homeDir {
					break
			}
			if len(subFiles) > 0 {
				subFiles = file + string(filepath.Separator) + subFiles
			} else {
					subFiles = file
			}
			if len(dir) == 0 || dir == prevDir {
					break
			}
			prevDir = dir[:len(dir) - 1]
	}
	return subFiles
}
