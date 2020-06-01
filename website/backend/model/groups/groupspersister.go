package groups

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/archives"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"path/filepath"
	"sync"
	"time"
)

type GroupsPersister struct {
	sync.RWMutex
	persister *persist.Persister

	groups []*GroupRecord
}

func NewGroupsPersister(dir string) (*GroupsPersister, error) {
	gp := &GroupsPersister{
		persister: persist.NewPersister("Groups", dir),
	}
	err := gp.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	gp.persister.SetPersistDelay(1 * time.Second)
	return gp, nil
}

func (gp *GroupsPersister) NbGroups() int {
	return len(gp.groups)
}

// LoadDirectory loads all persisted Actors Records
func (gp *GroupsPersister) LoadDirectory() error {
	gp.Lock()
	defer gp.Unlock()

	gp.persister.Reinit()
	gp.groups = []*GroupRecord{}

	files, err := gp.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from groups persister: %v", err)
	}

	for _, file := range files {
		ar, err := NewGroupRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not instantiate group from '%s': %v", filepath.Base(file), err)
		}
		err = gp.persister.Load(ar)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
		gp.groups = append(gp.groups, ar)
	}
	return nil
}

// Add adds the given GroupRecord to the GroupsPersister and return its (updated with new id) GroupRecord
func (gp *GroupsPersister) Add(ng *GroupRecord) *GroupRecord {
	gp.Lock()
	defer gp.Unlock()

	// give the record its new ID
	gp.persister.Add(ng)
	ng.Id = ng.GetId()
	gp.groups = append(gp.groups, ng)
	return ng
}

// Update updates the given GroupRecord
func (gp *GroupsPersister) Update(ugr *GroupRecord) error {
	gp.RLock()
	defer gp.RUnlock()

	ogr := gp.GetById(ugr.Id)
	if ogr == nil {
		return fmt.Errorf("actor id not found")
	}
	ogr.Group = ugr.Group
	gp.persister.MarkDirty(ogr)
	return nil
}

func (gp *GroupsPersister) findIndex(gr *GroupRecord) int {
	for i, rec := range gp.groups {
		if rec.GetId() == gr.GetId() {
			return i
		}
	}
	return -1
}

// Remove removes the given ActorRecord from the GroupsPersister (pertaining file is moved to deleted dir)
func (gp *GroupsPersister) Remove(ra *GroupRecord) error {
	gp.Lock()
	defer gp.Unlock()

	err := gp.persister.Remove(ra)
	if err != nil {
		return err
	}

	i := gp.findIndex(ra)
	copy(gp.groups[i:], gp.groups[i+1:])
	gp.groups[len(gp.groups)-1] = nil // or the zero value of T
	gp.groups = gp.groups[:len(gp.groups)-1]
	return nil
}

// GetById returns the GroupRecord with given Id (or nil if Id not found)
func (gp *GroupsPersister) GetById(id int) *GroupRecord {
	gp.RLock()
	defer gp.RUnlock()

	for _, gr := range gp.groups {
		if gr.Id == id {
			return gr
		}
	}
	return nil
}

func (gp *GroupsPersister) GetAllSites() []archives.ArchivableRecord {
	gp.RLock()
	defer gp.RUnlock()

	archivableSites := make([]archives.ArchivableRecord, len(gp.groups))
	for i, site := range gp.groups {
		archivableSites[i] = site
	}
	return archivableSites
}

func (gp *GroupsPersister) GetName() string {
	return "Groups"
}
