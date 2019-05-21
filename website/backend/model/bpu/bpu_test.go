package bpu

import (
	"encoding/json"
	"os"
	"testing"
)

const (
	testBpuXlsFile = `test/BPU.xlsx`
)

func TestNewBpuFromXLS(t *testing.T) {
	bpu, err := NewBpuFromXLS(testBpuXlsFile)
	if err != nil {
		t.Fatalf("NewBpuFromXLS returns unexpected: %s", err.Error())
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	enc.Encode(bpu)
}
