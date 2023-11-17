package resizeImage

import (
	"fmt"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"sync"

	"github.com/nfnt/resize"
	"github.com/schollz/progressbar/v3"
)

type ImageOptions struct {
	Width int
	Height int
	Quality int
}

type ResizeProgress struct {
	Progress *progressbar.ProgressBar
	Wg       *sync.WaitGroup
}

func (rp *ResizeProgress) Increment() {
	rp.Progress.Add(1)
}


func resizeImage(path , outPath string,  options ImageOptions, resizeInfo *ResizeProgress) error {
	file, err := os.Open(path)

	if err != nil {
  	fmt.Fprintln(os.Stderr, err, "cannot open image")
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err, "cannot read image")
	}
	file.Close()
	m := resize.Resize(uint(options.Width), uint(options.Height), img, resize.Lanczos3)
	out, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err, "cannot resize image")
	}
	defer out.Close()
	err = jpeg.Encode(out, m, &jpeg.Options{Quality: options.Quality})
	if err != nil {
		fmt.Fprintln(os.Stderr, err, "cannot save image")
	}
	resizeInfo.Increment()
	return nil
}

func sortFiles(files []os.DirEntry) {
	reg, _ := regexp.Compile(`\d+`)
	sort.Slice(files, func(i, j int) bool {
		file1 := files[i].Name()
		file2 := files[j].Name()
		numbers1 := reg.FindAllString(file1, -1)
		numbers2 := reg.FindAllString(file2, -1)
		minLength := math.Min(float64(len(numbers1)), float64(len(numbers2)))
		for i := 0; i < int(minLength); i++ {
			number1, _ := strconv.Atoi(numbers1[i])
			number2, _ := strconv.Atoi(numbers2[i])
			if number1 != number2 {
				if number1 < number2 {
					return true
				} else if number1 > number2 {
					return false
				}
			}
		}
		return len(numbers1) > len(numbers2)
	})
}

func ProcessDirs(inputDir, outDir string, imgOptions ImageOptions, resizeInfo *ResizeProgress) error {
	defer resizeInfo.Wg.Done()
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		return err
	}
	sortFiles(files)
	count := 1
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".jpg" {
			continue
		}
		vrName := fmt.Sprintf("VR%02d.jpg", count)
		err = resizeImage(filepath.Join(inputDir, file.Name()), filepath.Join(outDir, vrName), imgOptions, resizeInfo)
		if err != nil {
			return err
		}
		count++
	}
	return nil
}
