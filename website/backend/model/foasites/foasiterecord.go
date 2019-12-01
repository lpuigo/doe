package foasites

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"io"
	"os"
)

type FoaSiteRecord struct {
	*persist.Record
	*FoaSite
}

// NewFoaSiteRecord returns a new FoaSiteRecord
func NewFoaSiteRecord() *FoaSiteRecord {
	psr := &FoaSiteRecord{}
	psr.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(psr.FoaSite)
	})
	return psr
}

// NewFoaSiteRecordFrom returns a FoaSiteRecord populated from the given reader
func NewFoaSiteRecordFrom(r io.Reader) (psr *FoaSiteRecord, err error) {
	psr = NewFoaSiteRecord()
	err = json.NewDecoder(r).Decode(psr)
	if err != nil {
		psr = nil
		return
	}
	psr.SetId(psr.Id)
	return
}

// NewFoaSiteRecordFromFile returns a FoaSiteRecord populated from the given file
func NewFoaSiteRecordFromFile(file string) (fsr *FoaSiteRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	fsr, err = NewFoaSiteRecordFrom(f)
	if err != nil {
		fsr = nil
		return
	}
	if fsr.FoaSite.UpdateDate != "" {
		return
	}
	fs, err := f.Stat()
	if err != nil {
		return nil, err
	}
	fsr.FoaSite.UpdateDate = date.Date(fs.ModTime()).String()
	return
}
