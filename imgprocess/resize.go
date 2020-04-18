package imgprocess

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"os"
)

func ImageConfig(file string) (config image.Config, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	config, _, err = image.DecodeConfig(f) // Image Struct
	return
}

//type FileInfo struct {
//}

func readImg(il *ImgLog) (image.Image, error) {
	img, err := imaging.Open(il.Path, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("error opening file:%v\n", err)
	}
	if img.Bounds().Size().X != il.Init.Width {
		il.Init.Width, il.Init.Height = il.Init.Height, il.Init.Width
	}
	return img, nil
}

func saveImg(file string, img image.Image, quality int) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return fmt.Errorf("error encoding image: %v\n", err)
	}
	return nil
}

func swapSaveImg(il *ImgLog, img image.Image, quality int) error {
	tmpFile := il.Path + "_tmp"
	err := saveImg(tmpFile, img, quality)
	if err != nil {
		return err
	}
	if err := os.Remove(il.Path); err != nil {
		return fmt.Errorf("error removing initial file: %v\n", err)
	}
	if err := os.Rename(tmpFile, il.Path); err != nil {
		return fmt.Errorf("error renaming temp file: %v\n", err)
	}
	il.Result, err = GetImageInfo(il.Path)
	if err != nil {
		return err
	}
	t := il.Init.Info.ModTime()
	return os.Chtimes(il.Path, t, t)
}

func ChangeQuality(il *ImgLog, quality int) error {
	img, err := readImg(il)

	if err != nil {
		return err
	}
	return swapSaveImg(il, img, quality)
}

func ResizeChangeQuality(il *ImgLog, ratio, quality int) error {
	img, err := readImg(il)
	if err != nil {
		return err
	}
	resImg := imaging.Resize(img, il.Init.Width/ratio, il.Init.Height/ratio, imaging.Lanczos) // imaging.NearestNeighbor / imaging.Linear / imaging.CatmullRom
	return swapSaveImg(il, resImg, quality)
}

func Sharpen(il *ImgLog, quality int) error {
	img, err := readImg(il)
	if err != nil {
		return err
	}
	resImg := imaging.Sharpen(img, 1) // imaging.NearestNeighbor / imaging.Linear / imaging.CatmullRom
	return swapSaveImg(il, resImg, quality)
}
