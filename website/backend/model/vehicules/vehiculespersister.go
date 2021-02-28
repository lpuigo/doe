package vehicules

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/lpuig/ewin/doe/website/backend/model/archives"
	"github.com/lpuig/ewin/doe/website/backend/persist"
)

type VehiculesPersister struct {
	sync.RWMutex
	persister *persist.Persister

	Vehicules []*VehiculeRecord
}

func NewVehiculesPersister(dir string) (*VehiculesPersister, error) {
	gp := &VehiculesPersister{
		persister: persist.NewPersister("Vehicules", dir),
	}
	err := gp.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	gp.persister.SetPersistDelay(1 * time.Second)
	return gp, nil
}

func (vp *VehiculesPersister) NbVehicules() int {
	return len(vp.Vehicules)
}

// LoadDirectory loads all persisted Vehicules Records
func (vp *VehiculesPersister) LoadDirectory() error {
	vp.Lock()
	defer vp.Unlock()

	vp.persister.Reinit()
	vp.Vehicules = []*VehiculeRecord{}

	files, err := vp.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from vehicules persister: %v", err)
	}

	for _, file := range files {
		ar, err := NewVehiculeRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not instantiate vehicule from '%s': %v", filepath.Base(file), err)
		}
		err = vp.persister.Load(ar)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
		vp.Vehicules = append(vp.Vehicules, ar)
	}
	return nil
}

// Add adds the given VehiculeRecord to the VehiculesPersister and return its (updated with new id) VehiculeRecord
func (vp *VehiculesPersister) Add(nv *VehiculeRecord) *VehiculeRecord {
	vp.Lock()
	defer vp.Unlock()

	// give the record its new ID
	vp.persister.Add(nv)
	nv.Id = nv.GetId()
	vp.Vehicules = append(vp.Vehicules, nv)
	return nv
}

// Update updates the given VehiculeRecord
func (vp *VehiculesPersister) Update(uvr *VehiculeRecord) error {
	vp.RLock()
	defer vp.RUnlock()

	ovr := vp.GetById(uvr.Id)
	if ovr == nil {
		return fmt.Errorf("vehicule id not found")
	}
	ovr.Vehicule = uvr.Vehicule
	vp.persister.MarkDirty(ovr)
	return nil
}

func (vp *VehiculesPersister) findIndex(vr *VehiculeRecord) int {
	for i, rec := range vp.Vehicules {
		if rec.GetId() == vr.GetId() {
			return i
		}
	}
	return -1
}

// Remove removes the given VehiculeRecord from the VehiculesPersister (pertaining file is moved to deleted dir)
func (vp *VehiculesPersister) Remove(rvr *VehiculeRecord) error {
	vp.Lock()
	defer vp.Unlock()

	err := vp.persister.Remove(rvr)
	if err != nil {
		return err
	}

	i := vp.findIndex(rvr)
	copy(vp.Vehicules[i:], vp.Vehicules[i+1:])
	vp.Vehicules[len(vp.Vehicules)-1] = nil // or the zero value of T
	vp.Vehicules = vp.Vehicules[:len(vp.Vehicules)-1]
	return nil
}

// GetById returns the VehiculeRecord with given Id (or nil if Id not found)
func (vp *VehiculesPersister) GetById(id int) *VehiculeRecord {
	vp.RLock()
	defer vp.RUnlock()

	for _, gr := range vp.Vehicules {
		if gr.Id == id {
			return gr
		}
	}
	return nil
}

func (vp *VehiculesPersister) GetVehicules() []*Vehicule {
	vp.RLock()
	defer vp.RUnlock()

	res := make([]*Vehicule, len(vp.Vehicules))
	for i, vr := range vp.Vehicules {
		res[i] = vr.Vehicule
	}
	return res
}

func (vp *VehiculesPersister) UpdateVehicules(updatedVehicules []*Vehicule) error {
	for _, updV := range updatedVehicules {
		uvr := NewVehiculeRecordFromVehicule(updV)
		if updV.Id < 0 { // New Vehicule, add it instead of update
			vp.Add(uvr)
			continue
		}
		err := vp.Update(uvr)
		if err != nil {
			fmt.Errorf("could not update Vehicule '%s %s' (id: %d)", uvr.Type, uvr.Immat, uvr.Id)
		}
	}
	return nil
}

func (vp *VehiculesPersister) GetAllSites() []archives.ArchivableRecord {
	vp.RLock()
	defer vp.RUnlock()

	archivableSites := make([]archives.ArchivableRecord, len(vp.Vehicules))
	for i, site := range vp.Vehicules {
		archivableSites[i] = site
	}
	return archivableSites
}

func (vp *VehiculesPersister) GetName() string {
	return "Vehicules"
}
