package main

import (
	"fmt"
	imp "github.com/lpuig/ewin/doe/imgprocess"
	"github.com/lxn/walk"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	. "github.com/lxn/walk/declarative"
)

type GuiContext struct {
	textEdit *walk.TextEdit
	msg      chan string
}

const (
	ParallelRoutine = 12
)

func main() {
	gc := GuiContext{}
	_, err := MainWindow{
		Title:   "EWIN Services DOE shrink",
		MinSize: Size{640, 480},
		Layout:  VBox{},
		OnDropFiles: func(files []string) {
			go gc.GoShrink(files)
		},
		Children: []Widget{
			Label{Text: "Glisser un répertoire DOE ici ..."},
			TextEdit{
				AssignTo:  &gc.textEdit,
				ReadOnly:  true,
				VScroll:   true,
				MaxLength: 100 * 1024,
			},
		},
	}.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (gc *GuiContext) GoShrink(files []string) {
	if gc.msg != nil {
		return
	}
	wg := sync.WaitGroup{}
	gc.msg = make(chan string)
	wg.Add(1)
	go func() {
		for msg := range gc.msg {
			gc.textEdit.AppendText(msg)
		}
		wg.Done()
	}()

	gc.Shrink(files)
	close(gc.msg)
	wg.Wait()
	gc.msg = nil
}

func (gc GuiContext) Logln(text string) {
	gc.msg <- text + "\r\n"
}

func (gc GuiContext) Logf(format string, arg ...interface{}) {
	gc.msg <- fmt.Sprintf(format, arg...)
}

func (gc GuiContext) Shrink(filenames []string) {
	for _, filename := range filenames {
		fi, err := os.Stat(filename)
		if err != nil {
			gc.Logf("Error : %v\r\n", err)
			continue
		}
		if !fi.IsDir() {
			gc.Logf("Skip file : %s\r\n", filename)
			continue
		}
		gc.Logf("Shrinking %s :\r\n", filepath.Base(filename))
		gc.Process(filename)
	}
	gc.Logln("Done")
}

func (gc GuiContext) genImgList(path string) ([]imp.ImgLog, error) {
	return imp.GetImgList(path, func(path string) (ok bool, imInf imp.ImgInfo) {
		imInf, err := imp.GetImageInfo(path)
		if err != nil {
			ok = false
			return
		}
		if imInf.MaxSize() > 3200 {
			ok = true
		}
		return
	})
}

func truncpath(root, path string) string {
	return strings.Replace(path, root, ".", 1)
}

func (gc GuiContext) Process(path string) {
	l, err := gc.genImgList(path)
	if err != nil {
		gc.Logf("Error: %v\r\n", err)
		return
	}
	imp.ProcessImgList(l, ParallelRoutine, func(il *imp.ImgLog) {
		//err := ChangeQuality(il.Path, 40)
		err := imp.ReduceChangeQuality(il, il.Init.Width/2, il.Init.Height/2, 70)
		shortPath := truncpath(path, il.Path)
		if err != nil {
			il.Err = fmt.Errorf("  could not change quality for %s: %v", shortPath, err)
			gc.Logln(il.Err.Error())
			return
		}
		iires, err := imp.GetImageInfo(il.Path)
		if err != nil {
			il.Err = fmt.Errorf("  could not get image info for %s: %v", shortPath, err)
			gc.Logln(il.Err.Error())
			return
		}
		il.Result = iires
		gc.Logf("  process %s\r\n", shortPath)
	})
}
