package actorinfos

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"io"
	"os"
)

type ActorInfoRecord struct {
	*persist.Record
	*ActorInfo
}

// NewActorInfoRecord returns a new ActorInfoRecord
func NewActorInfoRecord() *ActorInfoRecord {
	air := &ActorInfoRecord{}
	air.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(air.ActorInfo)
	})
	air.ActorInfo = NewActorInfo()
	return air
}

// NewActorInfoRecordFrom returns a ActorInfoRecord populated from the given reader
func NewActorInfoRecordFrom(r io.Reader) (air *ActorInfoRecord, err error) {
	air = NewActorInfoRecord()
	err = json.NewDecoder(r).Decode(air)
	if err != nil {
		air = nil
		return
	}
	air.SetId(air.Id)
	return
}

// NewActorInfoRecordFromFile returns a ActorInfoRecord populated from the given file
func NewActorInfoRecordFromFile(file string) (air *ActorInfoRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	air, err = NewActorInfoRecordFrom(f)
	if err != nil {
		air = nil
		return
	}
	return
}

// NewActorInfoRecordFromActorInfo returns a ActorInfoRecord populated from given actor
func NewActorInfoRecordFromActorInfo(actinf *ActorInfo) *ActorInfoRecord {
	air := NewActorInfoRecord()
	air.ActorInfo = actinf
	air.SetId(air.Id)
	return air
}

// NewActorInfoRecordForActor returns a New ActorInfoRecord set from given actor
func NewActorInfoRecordForActor(actor *actors.Actor) *ActorInfoRecord {
	air := NewActorInfoRecord()
	air.ActorId = actor.Id
	return air

}
