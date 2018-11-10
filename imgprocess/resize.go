package imgprocess

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"os"
)

func Config(file string) (config image.Config, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	config, _, err = image.DecodeConfig(f) // Image Struct
	return
}

type FileInfo struct {
}

func readImg(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f) // Image Struct
	if err != nil {
		return nil, fmt.Errorf("error decoding file:%v\n", err)
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

func ChangeQuality(file string, quality int) error {
	img, err := readImg(file)
	if err != nil {
		return err
	}
	tmpFile := file + "_tmp"
	err = saveImg(tmpFile, img, quality)
	if err != nil {
		return err
	}
	if err := os.Remove(file); err != nil {
		return fmt.Errorf("error removing initial file: %v\n", err)
	}
	if err := os.Rename(tmpFile, file); err != nil {
		return fmt.Errorf("error renaming temp file: %v\n", err)
	}
	return nil
}

func ReduceChangeQuality(il *ImgLog, w, h, quality int) error {
	img, err := readImg(il.Path)
	if err != nil {
		return err
	}

	resImg := imaging.Sharpen(imaging.Resize(img, w, h, imaging.Box), 1) // imaging.NearestNeighbor / imaging.Linear / imaging.CatmullRom

	tmpFile := il.Path + "_tmp"
	err = saveImg(tmpFile, resImg, quality)
	if err != nil {
		return err
	}
	//if err := os.Remove(file); err != nil {
	//	return fmt.Errorf("error removing initial file: %v\n", err)
	//}
	if err := os.Rename(tmpFile, il.Path); err != nil {
		return fmt.Errorf("error renaming temp file: %v\n", err)
	}
	t := il.Init.Info.ModTime()
	return os.Chtimes(il.Path, t, t)
}
