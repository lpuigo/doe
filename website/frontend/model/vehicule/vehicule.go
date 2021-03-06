package vehicule

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/vehicule/vehiculeconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

// Type Vehicule reflects ewin/doe/website/backend/model/vehicules.Vehicule
type Vehicule struct {
	*js.Object

	Id             int             `js:"Id"`
	Type           string          `js:"Type"`
	Model          string          `js:"Model"`
	Company        string          `js:"Company"`
	Immat          string          `js:"Immat"`
	InCharge       []*ActorHistory `js:"InCharge"`
	ServiceDate    string          `js:"ServiceDate"`
	EndServiceDate string          `js:"EndServiceDate"`
	FuelCard       string          `js:"FuelCard"`
	Comment        string          `js:"Comment"`

	TravelledKms []*TravelledHistory `js:"TravelledKms"`
	Inventories  []*Inventory        `js:"Inventories"`
	Events       []*Event            `js:"Events"`
}

func NewVehicule() *Vehicule {
	nv := &Vehicule{Object: tools.O()}
	nv.Id = -1
	nv.Type = vehiculeconst.TypeCar
	nv.Model = ""
	nv.Company = ""
	nv.Immat = ""
	nv.InCharge = []*ActorHistory{NewActorHistory()}
	nv.ServiceDate = date.TodayAfter(0)
	nv.EndServiceDate = ""
	nv.FuelCard = ""
	nv.Comment = ""

	nv.TravelledKms = []*TravelledHistory{}
	nv.Inventories = []*Inventory{}
	nv.Events = []*Event{}

	return nv
}

func VehiculeFromJS(obj *js.Object) *Vehicule {
	return &Vehicule{Object: obj}
}

func (v *Vehicule) Copy() *Vehicule {
	return VehiculeFromJS(json.Parse(json.Stringify(v.Object)))
}

func (v *Vehicule) Clone(ov *Vehicule) {
	v.Id = ov.Id
	v.Type = ov.Type
	v.Model = ov.Model
	v.Company = ov.Company
	v.Immat = ov.Immat

	inCharge := make([]*ActorHistory, len(ov.InCharge))
	for i, ah := range ov.InCharge {
		inCharge[i] = ah.Copy()
	}
	v.InCharge = inCharge

	v.ServiceDate = ov.ServiceDate
	v.EndServiceDate = ov.EndServiceDate
	v.FuelCard = ov.FuelCard
	v.Comment = ov.Comment

	travelledKms := make([]*TravelledHistory, len(ov.TravelledKms))
	for i, th := range ov.TravelledKms {
		travelledKms[i] = th.Copy()
	}
	v.TravelledKms = travelledKms

	inventories := make([]*Inventory, len(ov.Inventories))
	for i, inv := range ov.Inventories {
		inventories[i] = inv.Copy()
	}
	v.Inventories = inventories

	v.Events = []*Event{}
	events := make([]*Event, len(ov.Events))
	for i, ev := range ov.Events {
		events[i] = ev.Copy()
	}
	v.Events = events
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

func (v *Vehicule) Status() string {
	if !tools.Empty(v.EndServiceDate) && v.EndServiceDate < date.TodayAfter(0) {
		return vehiculeconst.StatusReturned
	}
	//TODO detect vehiculeconst.StatusInRepair
	return vehiculeconst.StatusInUse
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

func (v *Vehicule) SortInCharge() {
	v.Get("InCharge").Call("sort", CompareActorHistory)
}

func (v *Vehicule) GetInChargeActorId(day string) int {
	for _, ah := range v.InCharge {
		if day >= ah.Date {
			return ah.ActorId
		}
	}
	return -1
}

// GetCurrentTravelledKms returns current (latest) TravelledHistory if exists, nil otherwise
func (v *Vehicule) GetCurrentTravelledKms() *TravelledHistory {
	if len(v.TravelledKms) > 0 {
		return v.TravelledKms[0]
	}
	return nil
}

func (v *Vehicule) SortTravelledKmsByDate() {
	v.Get("TravelledKms").Call("sort", CompareTravelledHistory)
}

func (v *Vehicule) AddInventory(ni *Inventory) {
	v.Inventories = append(v.Inventories, ni)
	v.SortInventoriesByDate()
}

func (v *Vehicule) SortInventoriesByDate() {
	v.Get("Inventories").Call("sort", CompareInventoryDate)
}

func (v *Vehicule) InventoryIndexByDate(day string) int {
	for iNum, inventory := range v.Inventories {
		if day >= inventory.ReferenceDate {
			return iNum
		}
	}
	return -1
}

func (v *Vehicule) SortEventsByDate() {
	v.Get("Events").Call("sort", CompareEventDate)
}

// GetInterestEvent returns the next future event (that is startdate >= today), or the lastest one (that is lastest event with startdate < today). Returns nil if no event available
func (v *Vehicule) GetInterestEvent() (*Event, bool) {
	today := date.TodayAfter(0)
	var resEvent *Event = nil
	future := false
	for _, event := range v.Events {
		if event.StartDate >= today {
			future = true
			resEvent = event
			continue
		}
		// now assert event.StartDate < today
		if resEvent == nil {
			//future = false
			resEvent = event
		}
		break
	}
	return resEvent, future
}
