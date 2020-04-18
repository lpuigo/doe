package imgprocess

import (
	"fmt"
	"os"
)

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

func (i *ImgInfo) CalcImgSize() {

}

func GetImageInfo(path string) (ImgInfo, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return ImgInfo{}, err
	}
	ic, err := ImageConfig(path)
	if err != nil {
		return ImgInfo{}, err
	}
	return ImgInfo{Info: fi, Width: ic.Width, Height: ic.Height}, nil
}

type ImgLog struct {
	Path string

	DoSharpen       bool
	DoResize        bool
	DoDowngradeQual bool

	Init   ImgInfo
	Result ImgInfo
	Err    error
}

type FilterFileFunc func(path string) (bool, bool, bool, ImgInfo)

func GetImgList(path string, filter FilterFileFunc) (list []ImgLog, err error) {
	pfunc := func(path string) error {
		sharpen, resize, downquality, imginfo := filter(path)
		if sharpen || resize || downquality {
			list = append(list, ImgLog{
				Path:            path,
				DoSharpen:       sharpen,
				DoResize:        resize,
				DoDowngradeQual: downquality,
				Init:            imginfo,
			})
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
		go func(imgLog *ImgLog) {
			pfunc(imgLog)
			<-wip
		}(&(list[i]))
	}

	// wait for goroutines completion
	for n := limit; n > 0; n-- {
		wip <- struct{}{}
	}
}
