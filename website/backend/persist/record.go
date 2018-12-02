package persist

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Record struct {
	id       int
	dirty    bool
	marshall func(w io.Writer) error
}

func (r Record) GetId() int {
	return r.id
}

func (r *Record) SetId(id int) {
	r.id = id
}

func (r *Record) Dirty() {
	r.dirty = true
}

func (r *Record) Persist(path string) error {
	file := filepath.Join(path, fmt.Sprintf("%06d.json", r.id))
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		return err
	}
	err = r.marshall(f)
	if err != nil {
		return fmt.Errorf("error marshalling: %v", err)
	}
	r.dirty = false
	return nil
}

func (r Record) Remove(path string) error {
	file := filepath.Join(path, fmt.Sprintf("%06d.json", r.id))
	return os.Remove(file)
}

func NewRecord(marshall func(w io.Writer) error) *Record {
	return &Record{marshall: marshall}
}

func (r *Record) SetIdFromFile(file string) error {
	_, err := fmt.Sscanf(filepath.Base(file), "%d.json", &r.id)
	return err
}
