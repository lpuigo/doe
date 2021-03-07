package vehiculeupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/vehicule"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

type VehiculeModalModel struct {
	*js.Object

	Visible bool     `js:"visible"`
	VM      *hvue.VM `js:"VM"`

	User            *fm.User           `js:"user"`
	EditedVehicule  *vehicule.Vehicule `js:"edited_vehicule"`
	CurrentVehicule *vehicule.Vehicule `js:"current_vehicule"`

	ShowConfirmDelete bool `js:"showconfirmdelete"`
}

func NewVehiculeModalModel(vm *hvue.VM) *VehiculeModalModel {
	vumm := &VehiculeModalModel{Object: tools.O()}
	vumm.Visible = false
	vumm.VM = vm

	vumm.User = fm.NewUser()
	vumm.EditedVehicule = vehicule.NewVehicule()
	vumm.CurrentVehicule = vehicule.NewVehicule()
	vumm.ShowConfirmDelete = false

	return vumm
}

func VehiculeModalModelFromJS(o *js.Object) *VehiculeModalModel {
	return &VehiculeModalModel{Object: o}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (vumm *VehiculeModalModel) HasChanged() bool {
	if vumm.EditedVehicule.Object == js.Undefined {
		return true
	}
	return json.Stringify(vumm.CurrentVehicule) != json.Stringify(vumm.EditedVehicule)
}

func (vumm *VehiculeModalModel) Show(vehic *vehicule.Vehicule, user *fm.User) {
	vumm.EditedVehicule = vehic
	vumm.CurrentVehicule = vehic.Copy()
	vumm.User = user
	vumm.ShowConfirmDelete = false
	vumm.Visible = true
}

func (vumm *VehiculeModalModel) HideWithControl() {
	if vumm.HasChanged() {
		message.ConfirmWarning(vumm.VM, "OK pour perdre les changements effectués ?", vumm.Hide)
		return
	}
	vumm.Hide()
}

func (vumm *VehiculeModalModel) Hide() {
	vumm.Visible = false
	vumm.ShowConfirmDelete = false
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Action Button Methods

func (vumm *VehiculeModalModel) ConfirmChange() {
	vumm.EditedVehicule.Clone(vumm.CurrentVehicule)
	//vumm.EditedVehicule.UpdateState()
	vumm.Hide()
}

func (vumm *VehiculeModalModel) UndoChange() {
	vumm.CurrentVehicule.Clone(vumm.EditedVehicule)
}
