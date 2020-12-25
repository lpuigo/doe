package users

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"io"
	"os"
)

type UserRecord struct {
	*persist.Record
	*User
}

// NewUserRecord returns a new UserRecord
func NewUserRecord() *UserRecord {
	ur := &UserRecord{}
	ur.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(ur.User)
	})
	return ur
}

// NewUserRecordFrom returns a UserRecord populated from the given reader
func NewUserRecordFrom(r io.Reader) (ur *UserRecord, err error) {
	ur = NewUserRecord()
	err = json.NewDecoder(r).Decode(ur)
	if err != nil {
		ur = nil
		return
	}
	ur.SetId(ur.Id)
	return
}

// NewUserRecordFromFile returns a UserRecord populated from the given file
func NewUserRecordFromFile(file string) (ur *UserRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	ur, err = NewUserRecordFrom(f)
	if err != nil {
		ur = nil
		return
	}
	return
}

// NewUserRecordFromUser returns a UserRecord populated from given User
func NewUserRecordFromUser(usr *User) *UserRecord {
	ur := NewUserRecord()
	ur.User = usr
	ur.SetId(ur.Id)
	return ur
}
