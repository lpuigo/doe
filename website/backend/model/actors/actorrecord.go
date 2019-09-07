package actors

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"io"
	"os"
)

type ActorRecord struct {
	*persist.Record
	*Actor
}

// NewActorRecord returns a new ActorRecord
func NewActorRecord() *ActorRecord {
	ar := &ActorRecord{}
	ar.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(ar.Actor)
	})
	return ar
}

// NewActorRecordFrom returns a ActorRecord populated from the given reader
func NewActorRecordFrom(r io.Reader) (ar *ActorRecord, err error) {
	ar = NewActorRecord()
	err = json.NewDecoder(r).Decode(ar)
	if err != nil {
		ar = nil
		return
	}
	ar.SetId(ar.Id)
	return
}

// NewActorRecordFromFile returns a ActorRecord populated from the given file
func NewActorRecordFromFile(file string) (ar *ActorRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	ar, err = NewActorRecordFrom(f)
	if err != nil {
		ar = nil
		return
	}
	return
}

// NewActorRecordFromActor returns a ActorRecord populated from given actor
func NewActorRecordFromActor(act *Actor) *ActorRecord {
	ar := NewActorRecord()
	ar.Actor = act
	ar.SetId(ar.Id)
	return ar
}
