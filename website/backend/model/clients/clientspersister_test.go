package clients

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	testDir   = "test"
	numClient = 10
)

func genClient(name string) *ClientRecord {
	cr := NewClientRecord()
	c := NewClient(name)
	c.Teams = append(c.Teams,
		Team{
			Name:     "1",
			Members:  "a, b",
			IsActive: true,
		},
		Team{
			Name:     "1",
			Members:  "a, c",
			IsActive: false,
		},
	)
	cr.Client = c
	return cr
}

func cleanTest(t *testing.T) {
	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if info.Name() == "deleted" {
				return filepath.SkipDir
			}
			return nil
		}
		return os.Remove(path)
	})
	if err != nil {
		t.Fatalf("cleanTest return: %v", err)
	}
}

func TestNewClientsPersister(t *testing.T) {
	cleanTest(t)
	cp, err := NewClientsPersister(testDir)
	if err != nil {
		t.Fatalf("NewClientsPersister returned unexpected: %v", err)
	}
	cp.persister.SetPersistDelay(10 * time.Millisecond)

	var cr *ClientRecord
	for i := 1; i <= numClient; i++ {
		cr = cp.Add(genClient(fmt.Sprintf("client-%04d", i)))
	}

	if len(cp.clients) != numClient {
		t.Fatalf("Clients array has unexpected length %d (expected %d)", len(cp.clients), numClient)
	}

	if cr.GetId() != numClient-1 {
		t.Fatalf("last added ClientRecord has unexpected id %d (expected %d)", cr.GetId(), numClient-1)
	}

	time.Sleep(100 * time.Millisecond)

}
