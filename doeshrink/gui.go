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
	sharpen  *walk.CheckBox
	msg      chan string
}

const (
	ParallelRoutine = 12
	MaxDim          = 3200
	MaxSize         = 1024 * 1024
	Quality         = 75
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
			CheckBox{
				Text:           "Améliorer la netteté",
				TextOnLeftSide: true,
				AssignTo:       &gc.sharpen,
				Alignment:      AlignHNearVCenter,
				OnClicked: func() {
					mode := "OFF"
					if gc.sharpen.Checked() {
						mode = "ON"
					}
					gc.AppendText(fmt.Sprintf("Filtre de netteté: %s\r\n", mode))
				},
			},
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

func (gc *GuiContext) AppendText(msg string) {
	gc.textEdit.AppendText(msg)
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
			gc.AppendText(msg)
		}
		wg.Done()
	}()

	gc.ProcessPaths(files)
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

func (gc GuiContext) ProcessPaths(filenames []string) {
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
		gc.Logf("Processing %s :\r\n", filepath.Base(filename))
		gc.Process(filename)
	}
	gc.Logln("Done")
}

func (gc GuiContext) genImgList(path string) ([]imp.ImgLog, error) {
	return imp.GetImgList(path, func(path string) (sharpen, resize, downquality bool, imInf imp.ImgInfo) {
		imInf, err := imp.GetImageInfo(path)
		if err != nil {
			return
		}
		if imInf.MaxSize() > MaxDim {
			resize = true
			return
		}
		if imInf.Info.Size() > MaxSize {
			downquality = true
			return
		}
		sharpen = gc.sharpen.Checked()
		return
	})
}

func subPath(root, path string) string {
	return strings.Replace(path, root, ".", 1)
}

func (gc GuiContext) Process(path string) {
	l, err := gc.genImgList(path)
	if err != nil {
		gc.Logf("Error: %v\r\n", err)
		return
	}
	imp.ProcessImgList(l, ParallelRoutine, func(il *imp.ImgLog) {
		shortPath := subPath(path, il.Path)
		switch {
		case il.DoResize:
			err := imp.ResizeChangeQuality(il, il.Init.Width/2, il.Init.Height/2, Quality)
			if err != nil {
				il.Err = fmt.Errorf("  could not resize %s: %v", shortPath, err)
				gc.Logln(il.Err.Error())
				return
			}
			gc.Logf("%s:\tImg dim reduced %d -> %d (%8d KB)\r\n", shortPath, il.Init.MaxSize(), il.Result.MaxSize(), il.Result.Info.Size()/1024)
		case il.DoDowngradeQual:
			err := imp.ChangeQuality(il, Quality)
			if err != nil {
				il.Err = fmt.Errorf("  could not change quality for %s: %v", shortPath, err)
				gc.Logln(il.Err.Error())
				return
			}
			gc.Logf("%s:\tImg file size reduced %8d -> %8d KB\r\n", shortPath, il.Init.Info.Size()/1024, il.Result.Info.Size()/1024)
		case il.DoSharpen:
			err := imp.Sharpen(il, Quality)
			if err != nil {
				il.Err = fmt.Errorf("  could not change quality for %s: %v", shortPath, err)
				gc.Logln(il.Err.Error())
				return
			}
			gc.Logf("%s:\tImg sharpened (%8d KB)\r\n", shortPath, il.Result.Info.Size()/1024)
		default:
			gc.Logf("%s:\tskipped\r\n", shortPath)
		}
	})
}
