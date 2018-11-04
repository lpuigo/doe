package model

import (
	"fmt"
	"strings"
)

type FormatedInfo struct {
	string
	prefix string
}

// SetPrefix sets the formatedinfo prefix.
func (fi *FormatedInfo) SetPrefix(f string) {
	fi.prefix = f
}

// Set sets the formatedinfo value.
//
// Remove prefix if found (must be called once prefix is set using SetPrefix)
func (fi *FormatedInfo) Set(i string) {
	fi.string = strings.TrimPrefix(i, fi.prefix)
}

func (fi FormatedInfo) String() string {
	return fmt.Sprintf("%s%s", fi.prefix, fi.string)
}

func (fi FormatedInfo) MarshalJSON() ([]byte, error) {
	return []byte("\"" + fi.String() + "\""), nil
}
