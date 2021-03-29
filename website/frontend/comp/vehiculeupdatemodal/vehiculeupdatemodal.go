package vehiculeupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/vehicule"
	"github.com/lpuig/ewin/doe/website/frontend/model/vehicule/vehiculeconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"strings"
)

type VehiculeUpdateModalModel struct {
	*VehiculeModalModel

	ActorStore *actor.ActorStore `js:"ActorStore"`

	VehiculeList              []*vehicule.Vehicule `js:"VehiculeList"`
	InventoryNum              int                  `js:"InventoryNum"`
	Control                   bool                 `js:"Control"`
	AddInventoryModelVisible  bool                 `js:"AddInventoryModelVisible"`
	AddInventoryModelSameType bool                 `js:"AddInventoryModelSameType"`
	AddInventoryModelVehicId  int                  `js:"AddInventoryModelVehicId"`
}

func NewVehiculeUpdateModalModel(vm *hvue.VM) *VehiculeUpdateModalModel {
	vumm := &VehiculeUpdateModalModel{VehiculeModalModel: NewVehiculeModalModel(vm)}
	vumm.ActorStore = actor.NewActorStore()
	vumm.VehiculeList = []*vehicule.Vehicule{}
	vumm.InventoryNum = -1
	vumm.Control = true
	vumm.AddInventoryModelVisible = false
	vumm.AddInventoryModelSameType = true
	vumm.AddInventoryModelVehicId = -1
	return vumm
}

func VehiculeUpdateModalModelFromJS(o *js.Object) *VehiculeUpdateModalModel {
	return &VehiculeUpdateModalModel{VehiculeModalModel: VehiculeModalModelFromJS(o)}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
		hvue.Computed("currentInventory", func(vm *hvue.VM) interface{} {
			vumm := VehiculeUpdateModalModelFromJS(vm.Object)
			if len(vumm.CurrentVehicule.Inventories) == 0 {
				vumm.InventoryNum = -1
			}
			if len(vumm.CurrentVehicule.Inventories) > 0 && vumm.InventoryNum == -1 {
				vumm.InventoryNum = 0
			}
			if vumm.InventoryNum < 0 {
				return nil
			}
			return vumm.CurrentVehicule.Inventories[vumm.InventoryNum]
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (vumm *VehiculeUpdateModalModel) Show(vehic *vehicule.Vehicule, user *fm.User, actorStore *actor.ActorStore, vehicList []*vehicule.Vehicule) {
	vumm.ActorStore = actorStore
	vumm.VehiculeList = vehicList
	vumm.VehiculeModalModel.Show(vehic, user)
	vumm.SetDefaultInventory()
}

func (vumm *VehiculeUpdateModalModel) ConfirmChange(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	vumm.VehiculeModalModel.ConfirmChange()
	vm.Emit("edited-vehicule", vumm.EditedVehicule)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

// InCharge Methods ====================================================================================================

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

// TravelledKms Methods ====================================================================================================

func (vumm *VehiculeUpdateModalModel) AddTravelledKms(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	nth := vehicule.NewTravelledHistory()
	nth.Date = date.TodayAfter(0)
	if len(vumm.CurrentVehicule.TravelledKms) > 0 {
		nth.Kms = vumm.CurrentVehicule.TravelledKms[0].Kms + 1000
	} else {
		nth.Kms = 1000
	}
	vumm.CurrentVehicule.TravelledKms = append([]*vehicule.TravelledHistory{nth}, vumm.CurrentVehicule.TravelledKms...)
}

func (vumm *VehiculeUpdateModalModel) RemoveTravelledKms(vm *hvue.VM, pos int) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	vumm.CurrentVehicule.Get("TravelledKms").Call("splice", pos, 1)
}

func (vumm *VehiculeUpdateModalModel) UpdateTravelledKms(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	vumm.CurrentVehicule.SortTravelledKmsByDate()
}

// Inventory Methods ===================================================================================================

func (vumm *VehiculeUpdateModalModel) SetDefaultInventory() {
	if len(vumm.CurrentVehicule.Inventories) == 0 {
		vumm.InventoryNum = -1
		return
	}
	vumm.InventoryNum = 0
	vumm.Control = true
}

func (vumm *VehiculeUpdateModalModel) UpdateInventoryNum(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	vumm.Control = true
}

func (vumm *VehiculeUpdateModalModel) GetInventoryDates(vm *hvue.VM) []*elements.IntValueLabel {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	if len(vumm.CurrentVehicule.Inventories) == 0 {
		return []*elements.IntValueLabel{elements.NewIntValueLabel(-1, vehiculeconst.InventoryNotFound)}
	}
	res := []*elements.IntValueLabel{}
	for iNum, inventory := range vumm.CurrentVehicule.Inventories {
		res = append(res, elements.NewIntValueLabel(iNum, date.DateString(inventory.ReferenceDate)))
	}
	return res
}

func (vumm *VehiculeUpdateModalModel) addInventory(ni *vehicule.Inventory) {
	vumm.CurrentVehicule.AddInventory(ni)
	vumm.InventoryNum = vumm.CurrentVehicule.InventoryIndexByDate(ni.ReferenceDate)
	vumm.Control = false
}

//func (vumm *VehiculeUpdateModalModel) AddInventory(vm *hvue.VM) {
//	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
//	vumm.addInventory(vehicule.NewInventory())
//}

func (vumm *VehiculeUpdateModalModel) AddInventoryModel(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	defer func() { vumm.AddInventoryModelVisible = false }()
	invModel := vumm.getInventoryModel()
	if invModel == nil {
		vumm.addInventory(vehicule.NewInventory())
		return
	}
	vumm.addInventory(vehicule.NewInventoryFromModel(invModel))
}

func (vumm *VehiculeUpdateModalModel) DeleteInventory(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	if vumm.InventoryNum < 0 {
		return
	}
	vumm.CurrentVehicule.Get("Inventories").Call("splice", vumm.InventoryNum, 1)
}

func (vumm *VehiculeUpdateModalModel) AddItem(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	if vumm.InventoryNum < 0 {
		return
	}
	currentInventory := vumm.CurrentVehicule.Inventories[vumm.InventoryNum]
	currentInventory.Items = append(currentInventory.Items, vehicule.NewInventoryItem())
}

func (vumm *VehiculeUpdateModalModel) RemoveItem(vm *hvue.VM, pos int) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	if vumm.InventoryNum < 0 {
		return
	}
	currentInventory := vumm.CurrentVehicule.Inventories[vumm.InventoryNum]
	currentInventory.Get("Items").Call("splice", pos, 1)
}

func (vumm *VehiculeUpdateModalModel) UpdateControlQuantity(inventItem *vehicule.InventoryItem) {
	inventItem.ControledQuantity = inventItem.ReferenceQuantity
}

func (vumm *VehiculeUpdateModalModel) GetInventoryModelVehiculeId(vm *hvue.VM) []*elements.IntValueLabel {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	res := []*elements.IntValueLabel{elements.NewIntValueLabel(-1, "Nouvel Inventaire")}
	today := date.TodayAfter(0)
	for _, v := range vumm.VehiculeList {
		if len(v.Inventories) == 0 {
			continue
		}
		if vumm.AddInventoryModelSameType && vumm.CurrentVehicule.Type != v.Type {
			continue
		}
		actId := v.GetInChargeActorId(today)
		actorRef := vehiculeconst.InChargeNotAffected
		if actId != -1 {
			act := vumm.ActorStore.GetActorById(actId)
			actorRef = act.Ref
		}

		label := v.Type + " " + v.Immat + " " + actorRef
		res = append(res, elements.NewIntValueLabel(v.Id, label))
	}
	return res
}

func (vumm *VehiculeUpdateModalModel) getInventoryModel() *vehicule.Inventory {
	if vumm.AddInventoryModelVehicId == -1 {
		return nil
	}
	var modelVehic *vehicule.Vehicule = nil
	for _, v := range vumm.VehiculeList {
		if v.Id == vumm.AddInventoryModelVehicId {
			modelVehic = v
			break
		}
	}
	if modelVehic == nil {
		return nil
	}
	if len(modelVehic.Inventories) == 0 {
		return nil
	}
	return modelVehic.Inventories[0]
}

func (vumm *VehiculeUpdateModalModel) GetInventoryModelItems(vm *hvue.VM) []*vehicule.InventoryItem {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	invModel := vumm.getInventoryModel()
	if invModel == nil {
		return []*vehicule.InventoryItem{}
	}
	return invModel.Items
}

func (vumm *VehiculeUpdateModalModel) ValidateInventoryControl(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	selectedInventory := vumm.CurrentVehicule.Inventories[vumm.InventoryNum]
	if tools.Empty(selectedInventory.ControledDate) {
		selectedInventory.ControledDate = date.TodayAfter(0)
	}
	vumm.addInventory(vehicule.NewInventoryFromControledModel(selectedInventory))
	vumm.InventoryNum = 0
	vumm.Control = false
}

func (vumm *VehiculeUpdateModalModel) AddEvent(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	vumm.CurrentVehicule.Events = append(vumm.CurrentVehicule.Events, vehicule.NewEvent())
	vumm.CurrentVehicule.SortEventsByDate()
}

// Event Methods ===================================================================================================

func (vumm *VehiculeUpdateModalModel) RemoveEvent(vm *hvue.VM, pos int) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	vumm.CurrentVehicule.Get("Events").Call("splice", pos, 1)
}

func (vumm *VehiculeUpdateModalModel) UpdateEvent(vm *hvue.VM) {
	vumm = VehiculeUpdateModalModelFromJS(vm.Object)
	vumm.CurrentVehicule.SortEventsByDate()
}

func (vumm *VehiculeUpdateModalModel) GetEventTypes(vm *hvue.VM) []*elements.ValueLabel {
	res := []*elements.ValueLabel{
		elements.NewValueLabel(vehiculeconst.EventTypeMisc, vehiculeconst.EventTypeMisc),
		elements.NewValueLabel(vehiculeconst.EventTypeCheck, vehiculeconst.EventTypeCheck),
		elements.NewValueLabel(vehiculeconst.EventTypeRepair, vehiculeconst.EventTypeRepair),
		elements.NewValueLabel(vehiculeconst.EventTypeIncident, vehiculeconst.EventTypeIncident),
	}
	return res
}
