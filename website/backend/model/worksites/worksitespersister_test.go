package worksites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/model"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	testDir = "test"
	numWS   = 10
)

func genWorksite(ref string) *model.Worksite {
	ws := model.MakeWorksite(
		ref,
		"2018-11-30",
		model.MakePT("PMZ-"+ref, "PT-007605", "02, Rue Kléber, CROIX"),
		model.MakePT("PA-"+ref, "PT-008020", "02, Rue Jean Jaurès, CROIX"),
	)
	return &ws
}

func cleanTest(t *testing.T) {
	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		os.Remove(path)
		return nil
	})
	if err != nil {
		t.Errorf("cleanTest return: %v", err)
	}
}

func getId(t *testing.T, wsr *WorkSiteRecord) int {
	var id int
	_, err := fmt.Sscanf(wsr.Ref, "WS-%d", &id)
	if err != nil {
		t.Fatalf("could not retrieved worksite id from ref '%s': %v", wsr.Ref, err)
	}
	return id
}

func TestNewWorkSitesPersist_Empty(t *testing.T) {
	cleanTest(t)
	wsp, err := NewWorkSitesPersist(testDir)
	if err != nil {
		t.Fatalf("could not create NewWorkSitesPersist: %v", err)
	}
	wsp.persister.SetPersistDelay(10 * time.Millisecond)

	var wsr *WorkSiteRecord
	for i := 1; i <= numWS; i++ {
		wsr = wsp.Add(genWorksite(fmt.Sprintf("WS-%04d", i)))
	}

	if len(wsp.workSites) != numWS {
		t.Fatalf("WorkSites array has unexpected length %d (expected %d)", len(wsp.workSites), numWS)
	}

	if wsr.GetId() != numWS-1 {
		t.Fatalf("last added WorkSitesRecord has unexpected id %d (expected %d)", wsr.GetId(), numWS-1)
	}

	var id int
	wslist := wsp.GetAll(func(ws *model.Worksite) bool {
		_, err := fmt.Sscanf(ws.Ref, "WS-%d", &id)
		if err != nil {
			t.Fatalf("could not retrieved worksite id from ref '%s': %v", ws.Ref, err)
		}
		return id%2 == 0
	})
	if len(wslist) != numWS/2 {
		t.Fatalf("WorkSitespersister.GetAll returns unexpected length %d (expected %d)", len(wslist), numWS/2)
	}
	time.Sleep(100 * time.Millisecond)
}

func initPopulatedWorkSitesPersister(t *testing.T) *WorkSitesPersister {
	TestNewWorkSitesPersist_Empty(t)

	wsp, err := NewWorkSitesPersist(testDir)
	if err != nil {
		t.Fatalf("could not create NewWorkSitesPersist: %v", err)
	}
	wsp.persister.SetPersistDelay(10 * time.Millisecond)
	err = wsp.LoadDirectory()
	if err != nil {
		t.Fatalf("WorkSitesPersist.LoadDirectory returns: %v", err)
	}
	return wsp
}

func TestWorkSitesPersister_LoadDirectory(t *testing.T) {
	wsp := initPopulatedWorkSitesPersister(t)
	wslist := wsp.GetAll(func(ws *model.Worksite) bool { return true })
	if len(wslist) != numWS {
		t.Fatalf("WorkSitesPersist.GetAll returns unexpected length %d (expected %d)", len(wslist), numWS)
	}
	for i, wsr := range wslist {
		id := getId(t, wsr)
		if id != wsr.GetId()+1 {
			t.Fatalf("WorkSiteRecord (Pos %d) Ref '%s' has unexpected GetId %d (expected %d)", i, wsr.Ref, id, wsr.GetId()+1)
		}
	}
}

func initHalfPopulatedWorkSitesPersister(t *testing.T) *WorkSitesPersister {
	wsp := initPopulatedWorkSitesPersister(t)
	wslist := wsp.GetAll(func(ws *model.Worksite) bool { return true })
	for _, wsr := range wslist[0 : numWS/2] {
		wsp.Remove(wsr)
	}
	return wsp
}

func TestWorkSitesPersister_Remove(t *testing.T) {
	wsp := initHalfPopulatedWorkSitesPersister(t)

	wslist := wsp.GetAll(func(ws *model.Worksite) bool { return true })
	if len(wslist) != numWS/2 {
		t.Fatalf("WorkSitesPersist.GetAll returns unexpected length %d (expected %d)", len(wslist), numWS/2)
	}
	time.Sleep(50 * time.Millisecond)
	files, err := wsp.persister.GetFilesList()
	if err != nil {
		t.Fatalf("GetFilesList returns : %v", err)
	}
	if len(files) != numWS/2 {
		t.Fatalf("GetFilesList returns unexpected number of files %d (expected %d)", len(files), numWS/2)
	}

	for i, wsr := range wslist {
		id := getId(t, wsr)
		if id != wsr.GetId()+1 {
			t.Fatalf("WorkSiteRecord (Pos %d) Ref '%s' has unexpected GetId %d (expected %d)", i, wsr.Ref, id, wsr.GetId()+1)
		}
	}
}

func TestWorkSitesPersister_Add(t *testing.T) {
	wsp := initHalfPopulatedWorkSitesPersister(t)
	wsp.Add(genWorksite(fmt.Sprintf("WS-%04d", numWS+1)))

	wslist := wsp.GetAll(func(ws *model.Worksite) bool { return true })
	if len(wslist) != numWS/2+1 {
		t.Fatalf("WorkSitesPersist.GetAll returns unexpected length %d (expected %d)", len(wslist), numWS/2+1)
	}

	time.Sleep(50 * time.Millisecond)
	files, err := wsp.persister.GetFilesList()
	if err != nil {
		t.Fatalf("GetFilesList returns : %v", err)
	}
	if len(files) != numWS/2+1 {
		t.Fatalf("GetFilesList returns unexpected number of files %d (expected %d)", len(files), numWS/2+1)
	}

	for i, wsr := range wslist {
		id := getId(t, wsr)
		if id != wsr.GetId()+1 {
			t.Fatalf("WorkSiteRecord (Pos %d) Ref '%s' has unexpected GetId %d (expected %d)", i, wsr.Ref, id, wsr.GetId()+1)
		}
	}
}
