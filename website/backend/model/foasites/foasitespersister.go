package foasites

import (
	"archive/zip"
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"io"
	"path/filepath"
	"sync"
	"time"
)

type FoaSitesPersister struct {
	sync.RWMutex
	persister *persist.Persister

	FoaSites []*FoaSiteRecord
}

func NewFoaSitesPersist(dir string) (*FoaSitesPersister, error) {
	fsp := &FoaSitesPersister{
		persister: persist.NewPersister("Foasites", dir),
	}
	err := fsp.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	fsp.persister.SetPersistDelay(1 * time.Second)
	return fsp, nil
}

func (fsp FoaSitesPersister) NbSites() int {
	return len(fsp.FoaSites)
}

// LoadDirectory loads all persisted FoaSiteRecords
func (fsp *FoaSitesPersister) LoadDirectory() error {
	fsp.Lock()
	defer fsp.Unlock()

	fsp.persister.Reinit()
	fsp.FoaSites = []*FoaSiteRecord{}

	files, err := fsp.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from FoaSites persister: %v", err)
	}

	for _, file := range files {
		psr, err := NewFoaSiteRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not instantiate FoaSite from '%s': %v", filepath.Base(file), err)
		}
		err = fsp.persister.Load(psr)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
		fsp.FoaSites = append(fsp.FoaSites, psr)
	}
	return nil
}

// GetAll returns all contained FoaSiteRecord for which keep(*FoaSite) == true
func (fsp FoaSitesPersister) GetAll(keep func(*FoaSite) bool) []*FoaSiteRecord {
	fsp.RLock()
	defer fsp.RUnlock()

	foaSiteRecords := []*FoaSiteRecord{}
	for _, psr := range fsp.FoaSites {
		if keep(psr.FoaSite) {
			foaSiteRecords = append(foaSiteRecords, psr)
		}
	}
	return foaSiteRecords
}

// GetById returns the FoaSiteRecord with given Id (or nil if Id not found)
func (fsp *FoaSitesPersister) GetById(id int) *FoaSiteRecord {
	fsp.RLock()
	defer fsp.RUnlock()

	for _, wsr := range fsp.FoaSites {
		if wsr.Id == id {
			return wsr
		}
	}
	return nil
}

// Add adds the given FoaSiteRecord to the FoaSitesPersister and return its (new) FoaSiteRecord
func (fsp *FoaSitesPersister) Add(psr *FoaSiteRecord) *FoaSiteRecord {
	fsp.Lock()
	defer fsp.Unlock()

	// Set the Update Date
	psr.FoaSite.UpdateDate = date.Today().String()

	// give the record its new ID
	fsp.persister.Add(psr)
	psr.Id = psr.GetId()
	fsp.FoaSites = append(fsp.FoaSites, psr)
	return psr
}

// Update updates the given FoaSiteRecord
func (fsp *FoaSitesPersister) Update(usr *FoaSiteRecord) error {
	fsp.RLock()
	defer fsp.RUnlock()

	osr := fsp.GetById(usr.Id)
	if osr == nil {
		return fmt.Errorf("id not found")
	}
	osr.FoaSite = usr.FoaSite
	osr.FoaSite.UpdateDate = date.Today().String()
	fsp.persister.MarkDirty(osr)
	return nil
}

// Remove removes the given FoaSiteRecord from the FoaSitesPersister (pertaining file is moved to deleted dir)
func (fsp *FoaSitesPersister) Remove(psr *FoaSiteRecord) error {
	fsp.Lock()
	defer fsp.Unlock()

	err := fsp.persister.Remove(psr)
	if err != nil {
		return err
	}

	i := fsp.findIndex(psr)
	copy(fsp.FoaSites[i:], fsp.FoaSites[i+1:])
	fsp.FoaSites[len(fsp.FoaSites)-1] = nil // or the zero value of T
	fsp.FoaSites = fsp.FoaSites[:len(fsp.FoaSites)-1]
	return nil
}

func (fsp FoaSitesPersister) findIndex(sr *FoaSiteRecord) int {
	for i, rec := range fsp.FoaSites {
		if rec.GetId() == sr.GetId() {
			return i
		}
	}
	return -1
}

// ArchiveName returns the FoaSiteArchive file name with today's date
func (fsp FoaSitesPersister) ArchiveName() string {
	return fmt.Sprintf("FoaSites %s.zip", date.Today().String())
}

// CreateArchive writes a zipped archive of all contained FoaSites files to the given writer
func (fsp *FoaSitesPersister) CreateArchive(writer io.Writer) error {
	fsp.RLock()
	defer fsp.RUnlock()

	zw := zip.NewWriter(writer)

	for _, sr := range fsp.FoaSites {
		wfw, err := zw.Create(sr.GetFileName())
		if err != nil {
			return fmt.Errorf("could not create zip entry for FoaSite %d", sr.Id)
		}
		err = sr.Marshall(wfw)
		if err != nil {
			return fmt.Errorf("could not write zip entry for FoaSite %d", sr.Id)
		}
	}

	return zw.Close()
}

// GetStats returns all Stats about all contained RipsiteRecords visible with isWSVisible = true and IsTeamVisible = true
//func (fsp *FoaSitesPersister) GetStats(sc items.StatContext, isPSVisible IsFoaSiteVisible, clientByName clients.ClientByName, actorById clients.ActorById, showprice bool) (*rs.RipsiteStats, error) {
//	fsp.RLock()
//	defer fsp.RUnlock()
//
//	// calc per Team/date/indicator values
//	calcValues := make(items.Stats)
//	for _, pr := range fsp.FoaSites {
//		if isPSVisible(pr.FoaSite) {
//			client := clientByName(pr.FoaSite.Client)
//			if client == nil {
//				continue
//			}
//			err := pr.AddStat(calcValues, sc, actorById, client.Bpu, showprice)
//			if err != nil {
//				return nil, err
//			}
//		}
//	}
//
//	d1 := func(s items.StatKey) string { return s.Serie } // Bars Family
//	d2 := func(s items.StatKey) string { return s.Team }  // Graphs
//	d3 := func(s items.StatKey) string { return s.Site }  // side block
//	f1 := items.KeepAll
//	//f2 := func(e string) bool { return !(!sc.ShowTeam && strings.Contains(e, " : ")) }
//	f2 := items.KeepAll
//	f3 := items.KeepAll
//	return calcValues.Aggregate(sc, d1, d2, d3, f1, f2, f3), nil
//}
