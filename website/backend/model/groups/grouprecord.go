package groups

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"io"
	"os"
)

type GroupRecord struct {
	*persist.Record
	*Group
}

// NewGroupRecord returns a new GroupRecord
func NewGroupRecord() *GroupRecord {
	gr := &GroupRecord{}
	gr.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(gr.Group)
	})
	return gr
}

// NewGroupRecordFrom returns a GroupRecord populated from the given reader
func NewGroupRecordFrom(r io.Reader) (gr *GroupRecord, err error) {
	gr = NewGroupRecord()
	err = json.NewDecoder(r).Decode(gr)
	if err != nil {
		gr = nil
		return
	}
	gr.SetId(gr.Id)
	return
}

// NewGroupRecordFromFile returns a GroupRecord populated from the given file
func NewGroupRecordFromFile(file string) (gr *GroupRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	gr, err = NewGroupRecordFrom(f)
	if err != nil {
		gr = nil
		return
	}
	return
}

// NewGroupRecordFromGroup returns a GroupRecord populated from given group
func NewGroupRecordFromGroup(grp *Group) *GroupRecord {
	gr := NewGroupRecord()
	gr.Group = grp
	gr.SetId(gr.Id)
	return gr
}
