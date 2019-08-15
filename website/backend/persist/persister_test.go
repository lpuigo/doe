package persist

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

const (
	persistDir = "test"
)

type Payload struct {
	Text  string
	Value int
}

type TestRecord struct {
	*Record
	Payload
}

func NewTestRecord(text string, val int) *TestRecord {
	tr := &TestRecord{
		Payload: Payload{
			Text:  text,
			Value: val,
		},
	}
	tr.Record = NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(&tr.Payload)
	})

	return tr
}

func NewTestRecordFromFile(file string) (tr *TestRecord, err error) {
	tr = NewTestRecord("empty", -1)
	err = tr.SetIdFromFile(file)
	if err != nil {
		tr = nil
		return
	}
	f, err := os.Open(file)
	if err != nil {
		tr = nil
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(tr)
	return
}

func genTestPersister(t *testing.T, numrec int, delay time.Duration) (*Persister, map[int]int) {
	trp := NewPersister(persistDir)
	trp.SetPersistDelay(delay)

	if err := trp.CheckDirectory(); err != nil {
		t.Fatal("checkDirectory return unexpected:", err)
	}

	res := make(chan struct{ id, val int }, numrec)
	for i := 1; i <= numrec; i++ {
		go func(n int) {
			id := trp.Add(NewTestRecord(fmt.Sprintf("record number %d", n), n))
			res <- struct{ id, val int }{id: id, val: n}
		}(i)
	}

	index := make(map[int]int)
	for i := 1; i <= numrec; i++ {
		s := <-res
		index[s.id] = s.val
	}
	return trp, index
}

func cleanTest(t *testing.T) {
	err := filepath.Walk(persistDir, func(path string, info os.FileInfo, err error) error {
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

func TestNewPersister(t *testing.T) {
	cleanTest(t)
	numrec := 100
	trp, index := genTestPersister(t, numrec, 10*time.Millisecond)

	if len(trp.records) != numrec {
		t.Errorf("Persister records has unexpected length (expected %d): %d", numrec, len(trp.records))
	}
	if len(trp.dirtyIds) != numrec {
		t.Errorf("Persister dirtyIds has unexpected length (expected %d): %d", numrec, len(trp.dirtyIds))
	}

	for i, r := range trp.records {
		tr, ok := r.(*TestRecord)
		if !ok {
			t.Errorf("record %d can not be casted back to TestRecord (type %v)", i, reflect.TypeOf(r))
			continue
		}
		if index[tr.id] != tr.Value {
			t.Errorf("record %d has unexepected value %d (expected %d)", tr.id, tr.Value, index[tr.id])
		}
	}

	time.Sleep(100 * time.Millisecond)
	if len(trp.dirtyIds) != 0 {
		t.Errorf("Persister dirtyIds has unexpected length (expected 0): %d", len(trp.records))
	}
}

func TestPersister_WaitPersistDone(t *testing.T) {
	cleanTest(t)
	numrec := 30
	trp, _ := genTestPersister(t, numrec, 10*time.Millisecond)

	if len(trp.dirtyIds) != numrec {
		t.Errorf("Persister dirtyIds has unexpected length before persisting (expected %d): %d", numrec, len(trp.records))
	}
	// test WaitPersistDone with active Persister
	trp.WaitPersistDone()
	if len(trp.dirtyIds) != 0 {
		t.Errorf("Persister dirtyIds has unexpected length after persisting (expected 0): %d", len(trp.records))
	}

	// test WaitPersistDone with inactive Persister
	trp.WaitPersistDone()
}

func TestPersister_GetFilesList(t *testing.T) {
	cleanTest(t)
	numrec := 30
	trp, index := genTestPersister(t, numrec, 10*time.Millisecond)
	time.Sleep(100 * time.Millisecond)
	files, err := trp.GetFilesList("deleted")
	if err != nil {
		t.Fatal("GetFilesList returns unexpected error:", err)
	}
	if len(files) != numrec {
		t.Fatalf("GetFilesList returns unexpected number of file %d (expected %d)", len(files), numrec)
	}
	var id int
	format := filepath.Join(persistDir, "%d.json")
	for _, file := range files {
		_, err := fmt.Sscanf(file, format, &id)
		if err != nil {
			t.Fatal("sscanf returns", err)
		}
		tr := getRecordFromFile(t, file)
		if tr.Value != index[id] {
			t.Errorf("file '%s' has unexpected value %d (expected %d)\n", filepath.Base(file), tr.Value, index[id])
		}
		delete(index, id)
	}
	if len(index) != 0 {
		t.Errorf("some id were not found: %v", index)
	}
}

func getRecordFromFile(t *testing.T, file string) TestRecord {
	var tr TestRecord
	f, err := os.Open(file)
	if err != nil {
		t.Fatal("could not open file:", err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&tr)
	if err != nil {
		t.Fatalf("could not decode '%s': %v", filepath.Base(file), err)
	}
	return tr
}

func TestPersister_LoadRecordFromFileList(t *testing.T) {
	cleanTest(t)
	numrec := 10
	trp, index := genTestPersister(t, numrec, 10*time.Millisecond)
	trp.WaitPersistDone()
	files, err := trp.GetFilesList("")
	if err != nil {
		t.Fatal("GetFilesList returns unexpected error:", err)
	}
	for _, file := range files {
		rt, err := NewTestRecordFromFile(file)
		if err != nil {
			t.Fatalf("could not load TestRecord from file '%s': %v\n", filepath.Base(file), err)
		}
		if rt.Value != index[rt.id] {
			t.Errorf("record has unexpected value %d (%d expected)", rt.Value, index[rt.id])
		}
	}
}
