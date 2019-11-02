package timesheets

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"io"
	"os"
)

type TimeSheetRecord struct {
	*persist.Record
	*TimeSheet
}

// NewTimeSheetRecord returns a new TimeSheetRecord
func NewTimeSheetRecord() *TimeSheetRecord {
	tsr := &TimeSheetRecord{}
	tsr.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(tsr.TimeSheet)
	})
	return tsr
}

// NewTimeSheetRecordFrom returns a TimeSheetRecord populated from the given reader
func NewTimeSheetRecordFrom(r io.Reader) (ar *TimeSheetRecord, err error) {
	ar = NewTimeSheetRecord()
	err = json.NewDecoder(r).Decode(ar)
	if err != nil {
		ar = nil
		return
	}
	ar.SetId(ar.Id)
	return
}

// NewTimeSheetRecordFromFile returns a TimeSheetRecord populated from the given file
func NewTimeSheetRecordFromFile(file string) (tsr *TimeSheetRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	tsr, err = NewTimeSheetRecordFrom(f)
	if err != nil {
		tsr = nil
		return
	}
	return
}

// NewTimeSheetRecordFromTimeSheet returns a TimeSheetRecord populated from given TimeSheet
func NewTimeSheetRecordFromTimeSheet(act *TimeSheet) *TimeSheetRecord {
	tsr := NewTimeSheetRecord()
	tsr.TimeSheet = act
	tsr.SetId(tsr.Id)
	return tsr
}
