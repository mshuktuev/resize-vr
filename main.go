package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"sync"

	resizeImage "github.com/mshuktuev/resize-vr/resize"
	"github.com/nfnt/resize"
	"github.com/schollz/progressbar/v3"
)


type PathInfo struct {
	Path string
	OutDir string
	SplitPath string
}

func main()  {
	images, dirs, err := getDirsInfo("./test")
	if err != nil {
		fmt.Println(err)
		return
	}
	_ =images
	_ =dirs

	bar := progressbar.Default(int64(images))

	resizeInfo := resizeImage.ResizeProgress{
		Progress: bar,
		Mutex:    &sync.Mutex{},
		Wg:       &sync.WaitGroup{},
	}


	// fmt.Println("Total images: ", images)
	// fmt.Println(dirs)

	file, err := os.Open("./test/VR_TH_Type_A_0001.jpg")

	resizeImage.ResizeImage("./test/TEst_2/Type_A/VR_TH_Type_A_0001.jpg", "output", "test", resizeImage.ImageOptions{
		Width: 4000,
		Height: 2000,
		Quality: 85,
	}, &resizeInfo)

	if err != nil {
  	fmt.Fprintln(os.Stderr, err, "cannot read image")
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err, "cannot read image")
	}
	file.Close()
	m := resize.Resize(4000, 2000, img, resize.Lanczos3)
	out, err := os.Create("test_resized.jpg")
	if err != nil {
		fmt.Fprintln(os.Stderr, err, "cannot read image")
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
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

func processDirs (path string, resizeInfo *resizeImage.ResizeProgress) error {
	resizeInfo.Wg.Add(1)
	dirs, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, dir := range dirs {

	}
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
