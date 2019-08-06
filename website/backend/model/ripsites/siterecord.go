package ripsites

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"io"
	"os"
)

type SiteRecord struct {
	*persist.Record
	*Site
}

// NewSiteRecord returns a new SiteRecord
func NewSiteRecord() *SiteRecord {
	sr := &SiteRecord{}
	sr.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(sr.Site)
	})
	return sr
}

// NewSiteRecordFrom returns a SiteRecord populated from the given reader
func NewSiteRecordFrom(r io.Reader) (sr *SiteRecord, err error) {
	sr = NewSiteRecord()
	err = json.NewDecoder(r).Decode(sr)
	if err != nil {
		sr = nil
		return
	}
	sr.SetId(sr.Id)
	return
}

// NewSiteRecordFromFile returns a SiteRecord populated from the given file
func NewSiteRecordFromFile(file string) (sr *SiteRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	sr, err = NewSiteRecordFrom(f)
	if err != nil {
		sr = nil
		return
	}
	if sr.Site.UpdateDate != "" {
		return
	}
	fs, err := f.Stat()
	if err != nil {
		return nil, err
	}
	sr.Site.UpdateDate = date.Date(fs.ModTime()).String()
	return
}
