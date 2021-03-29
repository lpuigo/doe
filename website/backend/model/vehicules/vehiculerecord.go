package vehicules

import (
	"encoding/json"
	"io"
	"os"

	"github.com/lpuig/ewin/doe/website/backend/persist"
)

type VehiculeRecord struct {
	*persist.Record
	*Vehicule
}

// NewVehiculeRecord returns a new VehiculeRecord
func NewVehiculeRecord() *VehiculeRecord {
	vr := &VehiculeRecord{}
	vr.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(vr.Vehicule)
	})
	return vr
}

// NewVehiculeRecordFrom returns a VehiculeRecord populated from the given reader
func NewVehiculeRecordFrom(r io.Reader) (vr *VehiculeRecord, err error) {
	vr = NewVehiculeRecord()
	err = json.NewDecoder(r).Decode(vr)
	if err != nil {
		vr = nil
		return
	}
	vr.SetId(vr.Id)
	vr.CheckConsistency()
	return
}

// NewVehiculeRecordFromFile returns a VehiculeRecord populated from the given file
func NewVehiculeRecordFromFile(file string) (vr *VehiculeRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	vr, err = NewVehiculeRecordFrom(f)
	if err != nil {
		vr = nil
		return
	}
	return
}

// NewVehiculeRecordFromVehicule returns a VehiculeRecord populated from given Vehicule
func NewVehiculeRecordFromVehicule(grp *Vehicule) *VehiculeRecord {
	gr := NewVehiculeRecord()
	gr.Vehicule = grp
	gr.SetId(gr.Id)
	return gr
}
