package actors

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/archives"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"path/filepath"
	"sync"
	"time"
)

type ActorsPersister struct {
	sync.RWMutex
	persister *persist.Persister

	actors []*ActorRecord
}

func NewActorsPersister(dir string) (*ActorsPersister, error) {
	ap := &ActorsPersister{
		persister: persist.NewPersister("Actors", dir),
	}
	err := ap.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	ap.persister.SetPersistDelay(1 * time.Second)
	return ap, nil
}

func (ap ActorsPersister) NbActors() int {
	return len(ap.actors)
}

// LoadDirectory loads all persisted Actors Records
func (ap *ActorsPersister) LoadDirectory() error {
	ap.Lock()
	defer ap.Unlock()

	ap.persister.Reinit()
	ap.actors = []*ActorRecord{}

	files, err := ap.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from actors persister: %v", err)
	}

	for _, file := range files {
		ar, err := NewActorRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not instantiate actor from '%s': %v", filepath.Base(file), err)
		}
		err = ap.persister.Load(ar)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
		ap.actors = append(ap.actors, ar)
	}
	return nil
}

// Add adds the given ActorRecord to the ActorsPersister and return its (updated with new id) ActorRecord
func (ap *ActorsPersister) Add(na *ActorRecord) *ActorRecord {
	ap.Lock()
	defer ap.Unlock()

	// give the record its new ID
	ap.persister.Add(na)
	na.Id = na.GetId()
	ap.actors = append(ap.actors, na)
	return na
}

// Update updates the given ActorRecord
func (ap *ActorsPersister) Update(ua *ActorRecord) error {
	ap.RLock()
	defer ap.RUnlock()

	oa := ap.GetById(ua.Id)
	if oa == nil {
		return fmt.Errorf("actor id not found")
	}
	oa.Actor = ua.Actor
	ap.persister.MarkDirty(oa)
	return nil
}

func (ap *ActorsPersister) findIndex(ar *ActorRecord) int {
	for i, rec := range ap.actors {
		if rec.GetId() == ar.GetId() {
			return i
		}
	}
	return -1
}

// Remove removes the given ActorRecord from the ActorsPersister (pertaining file is moved to deleted dir)
func (ap *ActorsPersister) Remove(ra *ActorRecord) error {
	ap.Lock()
	defer ap.Unlock()

	err := ap.persister.Remove(ra)
	if err != nil {
		return err
	}

	i := ap.findIndex(ra)
	copy(ap.actors[i:], ap.actors[i+1:])
	ap.actors[len(ap.actors)-1] = nil // or the zero value of T
	ap.actors = ap.actors[:len(ap.actors)-1]
	return nil
}

// GetById returns the ActorRecord with given Id (or nil if Id not found)
func (ap *ActorsPersister) GetById(id int) *ActorRecord {
	ap.RLock()
	defer ap.RUnlock()

	for _, ar := range ap.actors {
		if ar.Id == id {
			return ar
		}
	}
	return nil
}

// GetByRef returns the ActorRecord with given Ref (or nil if Id not found)
func (ap *ActorsPersister) GetByRef(ref string) *ActorRecord {
	ap.RLock()
	defer ap.RUnlock()

	for _, ar := range ap.actors {
		if ar.Ref == ref {
			return ar
		}
	}
	return nil
}

// GetAllActors returns a slice containing all persisted Actors
func (ap *ActorsPersister) GetAllActors() []*Actor {
	ap.RLock()
	defer ap.RUnlock()

	res := make([]*Actor, len(ap.actors))
	for i, ar := range ap.actors {
		res[i] = ar.Actor
	}
	return res
}

// GetActorsByClient returns all Actors (active as today if activeOnly is true), acting for given clients (all actors if clients is empty)
func (ap *ActorsPersister) GetActorsByClient(activeOnly bool, clients ...string) []*Actor {
	ap.RLock()
	defer ap.RUnlock()

	res := []*Actor{}
	today := date.Today().String()
	for _, ar := range ap.actors {
		if ar.Actor.WorksForClient(clients...) && (!activeOnly || ar.IsActiveOn(today)) {
			res = append(res, ar.Actor)
		}
	}
	return res
}

// UpdateActors updates all given updated actors
func (ap *ActorsPersister) UpdateActors(updatedActors []*Actor) error {
	for _, actor := range updatedActors {
		ar := NewActorRecordFromActor(actor)
		if actor.Id < 0 { // New Actor, add it instead of update
			ap.Add(ar)
			continue
		}
		err := ap.Update(ar)
		if err != nil {
			return fmt.Errorf("could not update actor '%s' (id: %d)", actor.Ref, actor.Id)
		}
	}
	return nil
}

func (ap *ActorsPersister) GetAllSites() []archives.ArchivableRecord {
	ap.RLock()
	defer ap.RUnlock()

	archivableSites := make([]archives.ArchivableRecord, len(ap.actors))
	for i, site := range ap.actors {
		archivableSites[i] = site
	}
	return archivableSites
}

func (ap *ActorsPersister) GetName() string {
	return "Actors"
}
