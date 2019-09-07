package actorupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

type ActorModalModel struct {
	*js.Object

	Visible bool     `js:"visible"`
	VM      *hvue.VM `js:"VM"`

	User         *fm.User     `js:"user"`
	EditedActor  *actor.Actor `js:"edited_actor"`
	CurrentActor *actor.Actor `js:"current_actor"`

	ShowConfirmDelete bool `js:"showconfirmdelete"`
}

func NewActorModalModel(vm *hvue.VM) *ActorModalModel {
	aumm := &ActorModalModel{Object: tools.O()}
	aumm.Visible = false
	aumm.VM = vm

	aumm.User = fm.NewUser()
	aumm.EditedActor = actor.NewActor()
	aumm.CurrentActor = actor.NewActor()
	aumm.ShowConfirmDelete = false

	return aumm
}

func ActorModalModelFromJS(o *js.Object) *ActorModalModel {
	return &ActorModalModel{Object: o}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (aumm *ActorModalModel) HasChanged() bool {
	if aumm.EditedActor.Object == js.Undefined {
		return true
	}
	return json.Stringify(aumm.CurrentActor) != json.Stringify(aumm.EditedActor)
}

func (aumm *ActorModalModel) Show(act *actor.Actor, user *fm.User) {
	aumm.EditedActor = act
	aumm.CurrentActor = act.Copy()
	aumm.User = user
	aumm.ShowConfirmDelete = false
	aumm.Visible = true
}

func (aumm *ActorModalModel) HideWithControl() {
	if aumm.HasChanged() {
		message.ConfirmWarning(aumm.VM, "OK pour perdre les changements effectués ?", aumm.Hide)
		return
	}
	aumm.Hide()
}

func (aumm *ActorModalModel) Hide() {
	aumm.Visible = false
	aumm.ShowConfirmDelete = false
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Action Button Methods

func (aumm *ActorModalModel) ConfirmChange() {
	aumm.EditedActor.Clone(aumm.CurrentActor)
	aumm.EditedActor.UpdateState()
	aumm.Hide()
}

func (aumm *ActorModalModel) UndoChange() {
	aumm.CurrentActor.Clone(aumm.EditedActor)
}
