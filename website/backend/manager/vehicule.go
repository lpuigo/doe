package manager

import (
	"encoding/json"
	"io"

	"github.com/lpuig/ewin/doe/website/backend/model/vehicules"
)

func (m Manager) GetVehicules(writer io.Writer) error {
	vehicules := m.Vehicules.GetVehicules()
	return json.NewEncoder(writer).Encode(vehicules)
}

func (m Manager) UpdateVehicules(updatedVehicules []*vehicules.Vehicule) error {
	return m.Vehicules.UpdateVehicules(updatedVehicules)
}
