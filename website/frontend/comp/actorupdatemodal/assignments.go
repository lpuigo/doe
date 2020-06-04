package actorupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/group"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"strconv"
)

type Assignment struct {
	*js.Object
	Date    string `js:"Date"`
	GroupId string `js:"GroupId"`
}

func NewAssignment(date string, groupId int) *Assignment {
	a := &Assignment{Object: tools.O()}
	a.Date = date
	a.GroupId = strconv.Itoa(groupId)
	return a
}

func NewAssignmentFromGroup(date string, group *group.Group) *Assignment {
	a := &Assignment{Object: tools.O()}
	a.Date = date
	a.GroupId = strconv.Itoa(group.Id)
	return a
}

type GroupsControl struct {
	*js.Object
	GroupStore  *group.GroupStore `js:"GroupStore"`
	Actor       *actor.Actor      `js:"Actor"`
	Assignments []*Assignment     `js:"Assignments"`
}

func NewGroupsControl(gs *group.GroupStore) *GroupsControl {
	gm := &GroupsControl{Object: tools.O()}
	gm.GroupStore = gs
	gm.Actor = nil
	gm.Assignments = []*Assignment{}
	return gm
}

func (gc *GroupsControl) SetAssignments(actor *actor.Actor) []*Assignment {
	gc.Actor = actor
	gc.Assignments = make([]*Assignment, len(actor.Groups))
	i := 0
	for assignDate, groupId := range actor.Groups {
		grp := gc.GroupStore.GetGroupById(groupId)
		if grp == nil {
			continue
		}
		gc.Get("Assignments").SetIndex(i, NewAssignment(assignDate, groupId))
		i++
	}
	gc.SortAssignments()
	return gc.Assignments
}

func (gc *GroupsControl) SortAssignments() {
	gc.Get("Assignments").Call("sort", func(a, b *Assignment) int {
		if a.Date < b.Date {
			return 1
		}
		if a.Date > b.Date {
			return -1
		}
		return 0
	})
}

func (gc *GroupsControl) GetCurrentAssignment() *Assignment {
	if len(gc.Assignments) == 0 {
		defAssignGroup := gc.GroupStore.GetGroupById(0)
		return NewAssignmentFromGroup(date.TodayAfter(0), defAssignGroup)
	}
	return gc.Assignments[0]
}

func (gc *GroupsControl) Add() {
	defAssignGroup := gc.GroupStore.GetGroupById(0)
	gc.Get("Assignments").Call("unshift", NewAssignmentFromGroup(date.TodayAfter(0), defAssignGroup))
	gc.SortAssignments()
}

func (gc *GroupsControl) Remove(index int) {
	if index >= len(gc.Assignments) {
		return
	}
	gc.Get("Assignments").Call("splice", index, 1)
}

func (gc *GroupsControl) UpdateActor() {
	ngh := actor.NewGroupHistory()
	for _, assignment := range gc.Assignments {
		grpId, err := strconv.Atoi(assignment.GroupId)
		if err != nil {
			continue
		}
		ngh[assignment.Date] = grpId
	}
	gc.Actor.Groups = ngh
}
