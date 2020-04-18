package imgprocess

import (
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
