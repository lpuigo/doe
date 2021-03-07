package vehiculeupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/vehicule"
	"github.com/lpuig/ewin/doe/website/frontend/model/vehicule/vehiculeconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"strings"
)

type VehiculeUpdateModalModel struct {
	*VehiculeModalModel

	ActorStore *actor.ActorStore `js:"ActorStore"`
}

func NewVehiculeUpdateModalModel(vm *hvue.VM) *VehiculeUpdateModalModel {
	vumm := &VehiculeUpdateModalModel{VehiculeModalModel: NewVehiculeModalModel(vm)}
	return vumm
}

func VehiculeUpdateModalModelFromJS(o *js.Object) *VehiculeUpdateModalModel {
	return &VehiculeUpdateModalModel{VehiculeModalModel: VehiculeModalModelFromJS(o)}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("vehicule-update-modal", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewVehiculeUpdateModalModel(vm)
		}),
		hvue.MethodsOf(&VehiculeUpdateModalModel{}),
		hvue.Computed("isNewVehicule", func(vm *hvue.VM) interface{} {
			vumm := VehiculeUpdateModalModelFromJS(vm.Object)
			return vumm.CurrentVehicule.Id < 0
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			vumm := VehiculeUpdateModalModelFromJS(vm.Object)
			return vumm.HasChanged()
		}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (vumm *VehiculeUpdateModalModel) Show(vehic *vehicule.Vehicule, user *fm.User, actorStore *actor.ActorStore) {
	vumm.ActorStore = actorStore
	vumm.VehiculeModalModel.Show(vehic, user)
}

func (vumm *VehiculeUpdateModalModel) ConfirmChange(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	vumm.VehiculeModalModel.ConfirmChange()
	vm.Emit("edited-vehicule", vumm.EditedVehicule)
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Tools Button Methods
func (vumm *VehiculeUpdateModalModel) CheckCompany(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)

	vumm.CurrentVehicule.Company = strings.Trim(strings.ToUpper(vumm.CurrentVehicule.Company), " \t")
}

func (vumm *VehiculeUpdateModalModel) GetVehiculeType() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(vehiculeconst.TypeCar, vehiculeconst.TypeCar),
		elements.NewValueLabel(vehiculeconst.TypeFourgon, vehiculeconst.TypeFourgon),
		elements.NewValueLabel(vehiculeconst.TypeNacelle, vehiculeconst.TypeNacelle),
		elements.NewValueLabel(vehiculeconst.TypeTariere, vehiculeconst.TypeTariere),
		elements.NewValueLabel(vehiculeconst.TypePorteTouret, vehiculeconst.TypePorteTouret),
	}
}

// InCharge Methods =============================================================================
func (vumm *VehiculeUpdateModalModel) AddInCharge(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	nah := vehicule.NewActorHistory()
	nah.Date = date.TodayAfter(0)
	vumm.CurrentVehicule.InCharge = append([]*vehicule.ActorHistory{nah}, vumm.CurrentVehicule.InCharge...)
}

func (vumm *VehiculeUpdateModalModel) RemoveInCharge(vm *hvue.VM, pos int) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	vumm.CurrentVehicule.Get("InCharge").Call("splice", pos, 1)
}

func (vumm *VehiculeUpdateModalModel) UpdateInCharge(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	vumm.CurrentVehicule.SortInCharge()
}

func (vumm *VehiculeUpdateModalModel) GetActors(vm *hvue.VM) []*elements.IntValueLabel {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	na := elements.NewIntValueLabel(-1, vehiculeconst.InChargeNotAffected)
	res := []*elements.IntValueLabel{na}
	for _, act := range vumm.ActorStore.GetActorsSortedByName() {
		res = append(res, elements.NewIntValueLabel(act.Id, act.GetRefStatus()))
	}
	return res
}
