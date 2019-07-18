package polesites

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"io"
	"os"
)

type PoleSiteRecord struct {
	*persist.Record
	*PoleSite
}

// NewSiteRecord returns a new PoleSiteRecord
func NewPoleSiteRecord() *PoleSiteRecord {
	psr := &PoleSiteRecord{}
	psr.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(psr.PoleSite)
	})
	return psr
}

// NewPoleSiteRecordFrom returns a PoleSiteRecord populated from the given reader
func NewPoleSiteRecordFrom(r io.Reader) (psr *PoleSiteRecord, err error) {
	psr = NewPoleSiteRecord()
	err = json.NewDecoder(r).Decode(psr)
	if err != nil {
		psr = nil
		return
	}
	psr.SetId(psr.Id)
	return
}

// NewPoleSiteRecordFromFile returns a PoleSiteRecord populated from the given file
func NewPoleSiteRecordFromFile(file string) (psr *PoleSiteRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	psr, err = NewPoleSiteRecordFrom(f)
	if err != nil {
		psr = nil
	}
	if psr.PoleSite.UpdateDate != "" {
		return
	}
	fs, err := f.Stat()
	if err != nil {
		return nil, err
	}
	psr.PoleSite.UpdateDate = date.Date(fs.ModTime()).String()
	return
}
