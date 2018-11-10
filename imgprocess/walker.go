package imgprocess

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProcessFileFunc func(path string) error

func Process(path string, pfunc ProcessFileFunc) error {
	err := filepath.Walk(path, processFn(pfunc))
	if err != nil {
		return err
	}
	return nil
}

func processFn(pfunc ProcessFileFunc) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// silently skip file with error
			return nil
		}
		if info.IsDir() {
			// skip directory
			return nil
		}
		name := info.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".jpg" && ext != ".jpeg" {
			return nil
		}
		return pfunc(path)
	}
}

type FilterFileFunc func(path string) (bool, ImgInfo)

type ImgInfo struct {
	Info   os.FileInfo
	Width  int
	Height int
}

func (i ImgInfo) FileSize() int {
	return int(i.Info.Size() / 1024)
}

func (i ImgInfo) MaxSize() (max int) {
	max = i.Height
	if i.Width > max {
		max = i.Width
	}
	return max
}

func (i ImgInfo) String() string {
	return fmt.Sprintf("%dKB (%d x %d)", i.FileSize(), i.Width, i.Height)
}

func GetImageInfo(path string) (ImgInfo, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return ImgInfo{}, err
	}
	c, err := Config(path)
	if err != nil {
		return ImgInfo{}, err
	}
	return ImgInfo{Info: fi, Width: c.Width, Height: c.Height}, nil
}

type ImgLog struct {
	Path   string
	Init   ImgInfo
	Result ImgInfo
	Err    error
}

func GetImgList(path string, filter FilterFileFunc) (list []ImgLog, err error) {
	pfunc := func(path string) error {
		if ok, imginfo := filter(path); ok {
			list = append(list, ImgLog{Path: path, Init: imginfo})
		}
		return nil
	}

	err = Process(path, pfunc)
	return
}

type ProcessImgFunc func(il *ImgLog)

func ProcessImgList(list []ImgLog, limit int, pfunc ProcessImgFunc) {
	wip := make(chan struct{}, limit)
	defer close(wip)

	for i, _ := range list {
		wip <- struct{}{}
		go func(n int, imgLog *ImgLog) {
			pfunc(imgLog)
			<-wip
		}(i, &(list[i]))
	}

	// wait for completion
	for n := limit; n > 0; n-- {
		wip <- struct{}{}
	}
}
