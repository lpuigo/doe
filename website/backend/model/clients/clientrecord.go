package clients

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"io"
	"os"
)

type ClientRecord struct {
	*persist.Record
	*Client
}

// NewClientRecord returns a new ClientRecord
func NewClientRecord() *ClientRecord {
	ur := &ClientRecord{}
	ur.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(ur.Client)
	})
	return ur
}

// NewClientRecordFrom returns a ClientRecord populated from the given reader
func NewClientRecordFrom(r io.Reader) (cr *ClientRecord, err error) {
	cr = NewClientRecord()
	err = json.NewDecoder(r).Decode(cr)
	if err != nil {
		cr = nil
		return
	}
	cr.SetId(cr.Id)
	return
}

// NewClientRecordFromFile returns a ClientRecord populated from the given file
func NewClientRecordFromFile(file string) (ur *ClientRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	ur, err = NewClientRecordFrom(f)
	if err != nil {
		ur = nil
		return
	}
	return
}

// NewClientRecordFromClient returns a ClientRecord populated from the given client
func NewClientRecordFromClient(clt *Client) *ClientRecord {
	cr := NewClientRecord()
	cr.Client = clt
	cr.SetId(cr.Id)
	return cr
}
