package converter

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ConvertEtx converts the image files in the specified directories to specified extension.
func ConvertEtx(src, from, to string) (int, error) {
	from = strings.ToLower(from)
	to = strings.ToLower(to)

	if err := validateArgs(from, to); err != nil {
		return 0, err
	}

	fileNames := make(chan string)
	go func() {
		walkDir(src, from, fileNames)
		close(fileNames)
	}()

	fileCnt := 0
	uniqCheck := make(map[string]int)
	for fn := range fileNames {
		file, err := os.Open(fn)
		if err != nil {
			return fileCnt, err
		}

		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			return fileCnt, err
		}

		fileName := filename(fn)
		if _, ok := uniqCheck[fileName]; !ok {
			uniqCheck[fileName] = 0
		} else {
			uniqCheck[fileName]++
			fileName = fileName + "(" + strconv.Itoa(uniqCheck[fileName]) + ")"
		}

		dstFile, err := os.Create(fmt.Sprintf("output/%s.%s", fileName, to))
		if err != nil {
			return fileCnt, err
		}
		defer dstFile.Close()

		switch to {
		case "jpg", "jpeg":
			err = jpeg.Encode(dstFile, img, nil)
		case "png":
			err = png.Encode(dstFile, img)
		}
		if err != nil {
			return fileCnt, err
		}
		fileCnt++
	}

	return fileCnt, nil
}

func filename(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func walkDir(dir, ext string, fileNames chan<- string) {
	ue := strings.ToUpper(ext)
	for _, ent := range dirents(dir) {
		if !strings.HasSuffix(ent.Name(), ext) && !strings.HasSuffix(ent.Name(), ue) && !ent.IsDir() {
			continue
		}
		if ent.IsDir() {
			subdir := filepath.Join(dir, ent.Name())
			walkDir(subdir, ext, fileNames)
		} else {
			fileNames <- filepath.Join(dir, ent.Name())
		}
	}
}

func dirents(dir string) []os.FileInfo {
	ents, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}

	return ents
}

func validateArgs(from, to string) error {
	ae := allowedExt{"jpeg", "png", "jpg"}

	if from == to {
		return errors.New("from and to are same")
	}
	if !ae.contains(from) {
		return errors.New("from is not supported")
	}
	if !ae.contains(to) {
		return errors.New("to is not supported")
	}

	return nil
}

type allowedExt []string

func (e allowedExt) contains(name string) bool {
	set := make(map[string]struct{}, len(e))
	for _, v := range e {
		set[v] = struct{}{}
	}

	_, ok := set[name]
	return ok
}
