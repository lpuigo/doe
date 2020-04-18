package imgprocess

import (
	"fmt"
	"path/filepath"
	"testing"
)

const TestDir string = `C:\Users\Laurent\Desktop\test\DOEs Termines - Copie`

func TestProcess(t *testing.T) {
	err := Process(TestDir, func(path string) error {
		fmt.Println(path)
		return nil
	})
	if err != nil {
		t.Fatal("Process returns", err)
	}
}

func genImgList(t *testing.T) []ImgLog {
	l, err := GetImgList(TestDir, func(path string) (ok bool, imInf ImgInfo) {
		imInf, err := GetImageInfo(path)
		if err != nil {
			ok = false
			return
		}
		if imInf.MaxSize() > 3200 {
			ok = true
		}
		return
	})
	if err != nil {
		t.Fatal("GetImgList returns", err)
	}
	return l
}

func TestGetImgList(t *testing.T) {
	l := genImgList(t)

	for _, il := range l {
		fmt.Printf("%s : %s\n", il.Path, il.Init.String())
	}
}

func TestProcessImgList(t *testing.T) {
	l := genImgList(t)

	ProcessImgList(l, 12, func(il *ImgLog) {
		//err := ChangeQuality(il.Path, 40)
		err := ResizeChangeQuality(il, il.Init.Width/2, il.Init.Height/2, 70)
		if err != nil {
			il.Err = fmt.Errorf("could not change quality for %s: %v\n", il.Path, err)
			return
		}
		iires, err := GetImageInfo(il.Path)
		if err != nil {
			il.Err = fmt.Errorf("could not get image info for %s: %v\n", il.Path, err)
			return
		}
		il.Result = iires
		fmt.Printf("process %s\n", filepath.Base(il.Path))
	})

	for _, il := range l {
		if il.Err != nil {
			t.Errorf("error: %s issued: %v\n", filepath.Base(il.Path), il.Err)
		}
		fmt.Printf("%s : %s -> %s\n", filepath.Base(il.Path), il.Init.String(), il.Result.String())
	}
}
