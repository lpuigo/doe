package vehicule

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/vehicule/vehiculeconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

// Type Vehicule reflects ewin/doe/website/backend/model/vehicules.Vehicule
type Vehicule struct {
	*js.Object

	Id          int            `js:"Id"`
	Type        string         `js:"Type"`
	Company     string         `js:"Company"`
	Immat       string         `js:"Immat"`
	InCharge    []ActorHistory `js:"InCharge"`
	ServiceDate string         `js:"ServiceDate"`
	Comment     string         `js:"Comment"`

	Inventories []Inventory `js:"Inventories"`
	Events      []Event     `js:"Events"`
}

func NewVehicule() *Vehicule {
	nv := &Vehicule{Object: tools.O()}
	nv.Id = -1
	nv.Type = ""
	nv.Company = ""
	nv.Immat = ""
	nv.InCharge = []ActorHistory{}
	nv.ServiceDate = ""
	nv.Comment = ""

	nv.Inventories = []Inventory{}
	nv.Events = []Event{}

	return nv
}

func VehiculeFromJS(obj *js.Object) *Vehicule {
	return &Vehicule{Object: obj}
}

func (v *Vehicule) Copy() *Vehicule {
	return VehiculeFromJS(json.Parse(json.Stringify(v.Object)))
}

func (v *Vehicule) SearchString(filter string) string {
	searchItem := func(prefix, typ, value string) string {
		if value == "" {
			return ""
		}
		if filter != vehiculeconst.FilterValueAll && filter != typ {
			return ""
		}
		return prefix + typ + value
	}

	res := searchItem("", vehiculeconst.FilterValueCompany, v.Company)
	res += searchItem("", vehiculeconst.FilterValueImmat, v.Immat)
	res += searchItem("", vehiculeconst.FilterValueType, v.Type)
	res += searchItem("", vehiculeconst.FilterValueComment, v.Comment)
	return res
}

func GetFilterTypeValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(vehiculeconst.FilterValueAll, vehiculeconst.FilterLabelAll),
		elements.NewValueLabel(vehiculeconst.FilterValueCompany, vehiculeconst.FilterLabelCompany),
		elements.NewValueLabel(vehiculeconst.FilterValueImmat, vehiculeconst.FilterLabelImmat),
		elements.NewValueLabel(vehiculeconst.FilterValueType, vehiculeconst.FilterLabelType),
		elements.NewValueLabel(vehiculeconst.FilterValueComment, vehiculeconst.FilterLabelComment),
	}
}
