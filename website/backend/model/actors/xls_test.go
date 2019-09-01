package actors

import (
	"os"
	"sort"
	"testing"
)

func LoadActorsFromXLS(t *testing.T, file string) []*Actor {
	f, err := os.Open(file)
	if err != nil {
		t.Fatalf("could not open file '%s': %s", file, err.Error())
	}
	defer f.Close()

	actors, err := FromXLS(f)
	if err != nil {
		t.Fatalf("FromXLS returns unexpected: %s", err.Error())
	}
	return actors
}

func TestFromXLS(t *testing.T) {
	file := `test/employees.xlsx`

	actors := LoadActorsFromXLS(t, file)
	for _, actor := range actors {
		t.Logf("%#v", actor)
	}
}

func TestXLStoJSON(t *testing.T) {
	file := `test/employees.xlsx`
	dir := `test`

	actors := LoadActorsFromXLS(t, file)

	ap, err := NewActorsPersister(dir)
	ap.persister.SetPersistDelay(0)
	if err != nil {
		t.Fatalf("NewActorsPersister returns unexpected: %s", err.Error())
	}

	// sort actors by Id
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].Id < actors[j].Id
	})

	for _, actor := range actors {
		ar := NewActorRecord()
		ar.Actor = actor

		ap.Add(ar)
	}

	ap.persister.WaitPersistDone()
}
