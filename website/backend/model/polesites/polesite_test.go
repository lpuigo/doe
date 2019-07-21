package polesites

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/model/polesites/test"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
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

	pole_Enrobe := map[string]bool{
		"MF22005": true,
		"MF22009": true,
		"MF22053": true,
		"MF22054": true,
		"MF22072": true,
		"MF22106": true,
		"MF22143": true,
		"MF22146": true,
		"MF22154": true,
		"MF22158": true,
		"MF22174": true,
		"MF22176": true,
		"MF22182": true,
		"MF22183": true,
		"MF22185": true,
		"MF22187": true,
		"MF22245": true,
		"MF22246": true,
		"MF22256": true,
		"MF22260": true,
		"MF22798": true,
		"MF22807": true,
		"MF22820": true,
		"MF22821": true,
		"MF22829": true,
		"MF22830": true,
		"MF22834": true,
		"MF22836": true,
		"MF22837": true,
		"MF22839": true,
		"MF22840": true,
		"MF22841": true,
	}
	pole_9m := map[string]bool{
		"MF22219": true,
		"MF22268": true,
		"MF22272": true,
		"MF22104": true,
		"MF22017": true,
		"MF22003": true,
		"MF22659": true,
		"MF22300": true,
		"MF22011": true,
		"MF22663": true,
		"MF22136": true,
		"MF22186": true,
		"MF22140": true,
		"MF22851": true,
		"MF22833": true,
		"MF22991": true,
		"MF22804": true,
		"MF22773": true,
	}

	for _, bp := range test.Poles {
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
			Material: poleconst.MaterialWood,
			Product:  map[string]int{},
			DictInfo: map[string]string{},
		}
		if pole_Enrobe[p.Ref] {
			p.Product[poleconst.ProductCoated] = 1
		}
		if pole_9m[p.Ref] {
			p.Height = 9
		}
		ps.Poles = append(ps.Poles, p)
	}

	json.NewEncoder(os.Stdout).Encode(ps)
}
