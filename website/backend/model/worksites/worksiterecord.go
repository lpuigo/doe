package worksites

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"io"
	"os"
)

type WorkSiteRecord struct {
	*persist.Record
	*model.Worksite
}

// NewWorkSiteRecord returns a new WorkSiteRecord
func NewWorkSiteRecord() *WorkSiteRecord {
	wsr := &WorkSiteRecord{}
	wsr.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(wsr.Worksite)
	})
	return wsr
}

// NewWorkSiteRecordFromFile returns a WorkSiteRecord populated from the given reader
func NewWorkSiteRecordFrom(r io.Reader) (wsr *WorkSiteRecord, err error) {
	wsr = NewWorkSiteRecord()
	err = json.NewDecoder(r).Decode(wsr)
	if err != nil {
		wsr = nil
		return
	}
	wsr.SetId(wsr.Id)
	return
}

// NewWorkSiteRecordFromFile returns a WorkSiteRecord populated from the given file
func NewWorkSiteRecordFromFile(file string) (wsr *WorkSiteRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	wsr, err = NewWorkSiteRecordFrom(f)
	if err != nil {
		wsr = nil
		return
	}
	wsr.SetId(wsr.Id)
	if wsr.Worksite.UpdateDate != "" {
		return
	}
	fs, err := f.Stat()
	if err != nil {
		return nil, err
	}
	wsr.Worksite.UpdateDate = date.Date(fs.ModTime()).String()
	return
}
