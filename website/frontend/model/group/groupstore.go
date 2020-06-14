package group

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"honnef.co/go/js/xhr"
)

type GroupStore struct {
	*js.Object

	Groups    []*Group `js:"Groups"`
	Reference string   `js:"Reference"`
	Dirty     bool     `js:"Dirty"`
}

func NewGroupStore() *GroupStore {
	gs := &GroupStore{Object: tools.O()}
	gs.Groups = []*Group{}
	gs.Reference = ""
	gs.Dirty = false
	return gs
}

func (gs *GroupStore) GetReference() string {
	return json.Stringify(gs.Groups)
}

func (gs *GroupStore) SetReference() {
	gs.Reference = gs.GetReference()
	gs.Dirty = false
}

func (gs *GroupStore) GetReferenceGroups() []*Group {
	refGroups := []*Group{}
	json.Parse(gs.Reference).Call("forEach", func(item *js.Object) {
		grp := GroupFromJS(item)
		refGroups = append(refGroups, grp)
	})
	return refGroups
}

func (gs *GroupStore) IsDirty() bool {
	gs.Dirty = gs.Reference == gs.GetReference()
	return gs.Dirty
}

func (gs *GroupStore) CallGetGroups(vm *hvue.VM, onSuccess func()) {
	go gs.callGetGroups(vm, onSuccess)
}

func (gs *GroupStore) callGetGroups(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("GET", "/api/groups")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(vm, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(vm, req)
		return
	}
	loadedGroups := []*Group{}
	req.Response.Call("forEach", func(item *js.Object) {
		grp := GroupFromJS(item)
		loadedGroups = append(loadedGroups, grp)
	})
	gs.Groups = loadedGroups
	gs.SetReference()
	onSuccess()
}

func (gs *GroupStore) CallUpdateGroups(vm *hvue.VM, onSuccess func()) {
	updatedGroups := gs.getUpdatedGroups()
	if len(updatedGroups) == 0 {
		message.ErrorStr(vm, "Could not find any updated actors", false)
		return
	}

	req := xhr.NewRequest("PUT", "/api/groups")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(updatedGroups))
	if err != nil {
		message.ErrorStr(vm, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(vm, req)
		return
	}
	message.NotifySuccess(vm, "Groupes", "Modifications sauvegardées")
	gs.SetReference()
	onSuccess()
}

func (gs *GroupStore) getUpdatedGroups() []*Group {
	updGrps := []*Group{}
	refGroups := gs.GetReferenceGroups()
	groupById := map[int]*Group{}
	for _, grp := range refGroups {
		groupById[grp.Id] = grp
	}
	for _, grp := range gs.Groups {
		refgrp := groupById[grp.Id]
		if !(refgrp != nil && json.Stringify(grp) != json.Stringify(refgrp)) {
			updGrps = append(updGrps, grp)
		}
	}
	return updGrps
}

// GetGroupById returns Group with given Id (nil if Id does not exist)
func (gs *GroupStore) GetGroupById(id int) *Group {
	for _, group := range gs.Groups {
		if group.Id == id {
			return group
		}
	}
	return nil
}

// GetGroupNameById returns Group"Name with given Id ("" if Id does not exist)
func (gs *GroupStore) GetGroupNameById(id int) string {
	gr := gs.GetGroupById(id)
	if gr != nil {
		return gr.Name
	}
	return ""
}
