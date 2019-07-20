package polesites

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"os"
	"testing"
)

func TestPoleSiteEncode(t *testing.T) {
	ps := &PoleSite{
		Id:         0,
		Client:     "Axians Moselle",
		Ref:        "Rechicourt",
		Manager:    "Goeffrey Wecker",
		OrderDate:  "2019-06-15",
		UpdateDate: "2019-07-11",
		Status:     "20 InProgress",
		Comment:    "test",
		Poles:      nil,
	}

	for i, bp := range polesite.Poles {
		p := &Pole{
			Ref:      bp.Ref,
			City:     bp.City,
			Address:  "",
			Lat:      bp.Lat,
			Long:     bp.Long,
			State:    bp.State,
			DtRef:    "",
			DictRef:  "",
			Height:   8,
			Product:  map[string]int{},
			DictInfo: map[string]string{},
		}
		if i < 10 {
			p.Product["Enrobé"] = 1
		}
		ps.Poles = append(ps.Poles, p)
	}

	json.NewEncoder(os.Stdout).Encode(ps)
}
