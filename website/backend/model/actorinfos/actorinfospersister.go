package actorinfos

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/model/archives"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"path/filepath"
	"sync"
	"time"
)

type ActorInfosPersister struct {
	sync.RWMutex
	persister *persist.Persister

	actorInfosById     map[int]*ActorInfoRecord
	actorInfoByActorId map[int]*ActorInfoRecord
}

func NewActorInfosPersister(dir string) (*ActorInfosPersister, error) {
	aip := &ActorInfosPersister{
		persister:          persist.NewPersister("ActorInfos", dir),
		actorInfosById:     make(map[int]*ActorInfoRecord),
		actorInfoByActorId: make(map[int]*ActorInfoRecord),
	}
	err := aip.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	aip.persister.SetPersistDelay(1 * time.Second)
	return aip, nil
}

func (aip ActorInfosPersister) NbActorInfos() int {
	return len(aip.actorInfosById)
}

// LoadDirectory loads all persisted Actors Records
func (aip *ActorInfosPersister) LoadDirectory() error {
	aip.Lock()
	defer aip.Unlock()

	aip.persister.Reinit()

	aip.actorInfosById = make(map[int]*ActorInfoRecord)
	aip.actorInfoByActorId = make(map[int]*ActorInfoRecord)

	files, err := aip.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from actorinfos persister: %v", err)
	}

	for _, file := range files {
		air, err := NewActorInfoRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not instantiate actorinfo from '%s': %v", filepath.Base(file), err)
		}
		err = aip.persister.Load(air)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
		aip.addActorInfo(air)
	}
	return nil
}

// addActorInfo adds given ActorInfoRecord to ActorInfosPersister
func (aip *ActorInfosPersister) addActorInfo(nair *ActorInfoRecord) {
	aip.actorInfosById[nair.Id] = nair
	aip.actorInfoByActorId[nair.ActorId] = nair
}

// deleteActorInfo deletes given ActorInfoRecord from ActorInfosPersister
func (aip *ActorInfosPersister) deleteActorInfo(rai *ActorInfoRecord) {
	delete(aip.actorInfosById, rai.Id)
	delete(aip.actorInfoByActorId, rai.ActorId)
}

// addNewActorInfo adds the given ActorInfoRecord to the ActorInfosPersister and return its (updated with new id) ActorInfoRecord
func (aip *ActorInfosPersister) addNewActorInfo(nair *ActorInfoRecord) *ActorInfoRecord {
	// give the record its new ID
	aip.persister.Add(nair)
	nair.Id = nair.GetId()
	aip.addActorInfo(nair)
	return nair
}

// GetById returns the ActorInfoRecord with given Id (or nil if Id not found)
func (aip *ActorInfosPersister) GetById(id int) *ActorInfoRecord {
	aip.RLock()
	defer aip.RUnlock()

	air, found := aip.actorInfosById[id]
	if !found {
		return nil
	}
	return air
}

// Update updates the given ActorInfoRecord
func (aip *ActorInfosPersister) Update(uair *ActorInfoRecord) error {
	aip.RLock()
	defer aip.RUnlock()

	oai := aip.GetById(uair.Id)
	if oai == nil {
		return fmt.Errorf("actorinfo id not found")
	}
	oai.ActorInfo = uair.ActorInfo
	aip.persister.MarkDirty(oai)
	return nil
}

// Remove removes the given ActorInfoRecord from the ActorInfosPersister (pertaining file is moved to deleted dir)
func (aip *ActorInfosPersister) Remove(rai *ActorInfoRecord) error {
	aip.Lock()
	defer aip.Unlock()

	err := aip.persister.Remove(rai)
	if err != nil {
		return err
	}

	aip.deleteActorInfo(rai)
	return nil
}

// GetAllActorInfos returns a slice containing all persisted ActorInfos
func (aip *ActorInfosPersister) GetAllActorInfos() []*ActorInfo {
	aip.RLock()
	defer aip.RUnlock()

	res := make([]*ActorInfo, len(aip.actorInfosById))
	for i, ar := range aip.actorInfosById {
		res[i] = ar.ActorInfo
	}
	return res
}

// GetActorInfosByActors returns a slice containing all persisted ActorInfos pertinaining to given actors
func (aip *ActorInfosPersister) GetActorInfosByActors(actors []*actors.Actor) []*ActorInfo {
	aip.RLock()
	defer aip.RUnlock()

	actorInfos := make([]*ActorInfo, len(actors))
	for i, actor := range actors {
		air, found := aip.actorInfoByActorId[actor.Id]
		if !found {
			air = NewActorInfoRecordForActor(actor)
		}
		actorInfos[i] = air.ActorInfo
	}
	return actorInfos
}

// GetActorInfosByActors returns a slice containing all persisted ActorInfos pertinaining to given actors
func (aip *ActorInfosPersister) GetActorHRsByActors(actors []*actors.Actor, addHRInfo bool) []*ActorHr {
	aip.RLock()
	defer aip.RUnlock()

	fakeActorInfo := NewActorInfo()

	actorHrs := make([]*ActorHr, len(actors))
	for i, actor := range actors {
		ahr := &ActorHr{Actor: actor}
		if addHRInfo {
			air, found := aip.actorInfoByActorId[actor.Id]
			if !found {
				air = NewActorInfoRecordForActor(actor)
			}
			ahr.Info = air.ActorInfo
		} else {
			ahr.Info = fakeActorInfo
		}
		actorHrs[i] = ahr
	}
	return actorHrs
}

// UpdateActorInfos adds (if not already known) or updates all given actorinfos
func (aip *ActorInfosPersister) UpdateActorInfos(updatedActorInfos []*ActorInfo) error {
	aip.Lock()
	defer aip.Unlock()

	for _, uai := range updatedActorInfos {
		oai, found := aip.actorInfoByActorId[uai.ActorId]
		if !found { // add actorinfo for new actorid
			oai = aip.addNewActorInfo(NewActorInfoRecordFromActorInfo(uai))
		} else {
			oai.ActorInfo = uai
			aip.persister.MarkDirty(oai)
		}
	}
	return nil
}

// Archive methods =====================================================================================================

func (aip *ActorInfosPersister) GetAllSites() []archives.ArchivableRecord {
	aip.RLock()
	defer aip.RUnlock()

	archivableSites := make([]archives.ArchivableRecord, len(aip.actorInfosById))
	i := 0
	for _, site := range aip.actorInfosById {
		archivableSites[i] = site
		i++
	}
	return archivableSites
}

func (aip *ActorInfosPersister) GetName() string {
	return "ActorInfos"
}
