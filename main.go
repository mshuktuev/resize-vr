package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	resizeImage "github.com/mshuktuev/resize-vr/resize"
	"github.com/schollz/progressbar/v3"
)


type PathInfo struct {
	Path string
	OutDir string
	SplitPath string
}

func main()  {
	runtime.GOMAXPROCS(6)
	start := time.Now()
	images, dirs, err := getDirsInfo("./test")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.RemoveAll("output")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Images: %d, dirs: %d\n", images, len(dirs))

	bar := progressbar.Default(int64(images) * 2)

	resizeInfo := resizeImage.ResizeProgress{
		Progress: bar,
		Wg:       &sync.WaitGroup{},
	}

	for dir := range dirs {
		outDir := subpath("test", dir)
		previewDir := filepath.Join("output", "preview", outDir)
		texturesDir := filepath.Join("output", "textures", outDir)
		resizeInfo.Wg.Add(2)
		go resizeDir(dir, previewDir, resizeImage.ImageOptions{
			Width: 2048,
			Height: 1024,
			Quality: 85,
		}, &resizeInfo)

		go resizeDir(dir, texturesDir, resizeImage.ImageOptions{
			Width: 4096,
			Height: 2048,
			Quality: 85,
		}, &resizeInfo)
	}
	resizeInfo.Wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Done in %s", elapsed)
}

func resizeDir(dir, outDir string, options resizeImage.ImageOptions, resizeInfo *resizeImage.ResizeProgress) {
	err := resizeImage.ProcessDirs(dir, outDir, options, resizeInfo)
	if err != nil {
		log.Fatal(err)
	}
}



func getDirsInfo(path string) (int, map[string]bool, error) {
	images := 0
	dirs := make(map[string]bool)
	err := filepath.Walk(path, func(wPath string, info os.FileInfo, err error) error {
		if wPath != path {
			if (filepath.Ext(wPath) == ".jpg") {
				images++
				rootDir := filepath.Dir(wPath)
				if _, ok := dirs[rootDir]; !ok {
					dirs[rootDir] = true
				}
			}
		}
		return nil
	})
	if err != nil {
		return 0, nil, err
	}

	return images, dirs, nil
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
