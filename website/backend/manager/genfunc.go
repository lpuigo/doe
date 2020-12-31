package manager

import (
	"strconv"

	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/groups"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
)

// ====================================================================================================
// Client Related Function ============================================================================

// genGetClient returns a GetClientByName function: func(clientName string) *clients.Client. Returned client is nil if clientName is not found
func (m *Manager) genGetClient() clients.ClientByName {
	clientByName := make(map[string]*clients.Client)
	for _, client := range m.Clients.GetAllClients() {
		clientByName[client.Name] = client
	}
	return func(clientName string) *clients.Client {
		return clientByName[clientName]
	}
}

// ====================================================================================================
// Group Related Function =============================================================================

// GenGroupById returns a groups.GroupById function
func (m Manager) GenGroupById() groups.GroupById {
	groupById := make(map[int]*groups.Group)
	for _, group := range m.Groups.GetGroups() {
		groupById[group.Id] = group
	}
	return func(groupId int) *groups.Group {
		return groupById[groupId]
	}
}

// GenGroupByName returns a groups.GroupByName function
func (m Manager) GenGroupByName() groups.GroupByName {
	groupByName := make(map[string]*groups.Group)
	for _, group := range m.Groups.GetGroups() {
		groupByName[group.Name] = group
	}
	return func(groupName string) *groups.Group {
		return groupByName[groupName]
	}
}

func (m Manager) GenActorsByGroupId() actors.ActorsByGroupId {
	today := date.Today().String()
	actorsByGroupId := make(map[int][]*actors.Actor)
	for _, actor := range m.Actors.GetAllActors() {
		groupId := actor.Groups.ActiveGroupOnDate(today)
		if groupId == -1 {
			continue
		}
		actorsList := actorsByGroupId[groupId]
		actorsByGroupId[groupId] = append(actorsList, actor)
	}
	return func(groupId int) []*actors.Actor {
		return actorsByGroupId[groupId]
	}
}

// ====================================================================================================
// Actor Related Function =============================================================================

// genActorById returns a ActorById function: func(actorId string, day string) *actors.Actor.
//
// Returned *actors.Actor is nil if actorId is not found or not Visible by manager.CurrentUser
func (m *Manager) genActorById(visibleActorsOnly bool) func(actorId string, day string) *actors.Actor {
	actDict := m.Actors.GetActorsDict()
	if !(visibleActorsOnly && !m.CurrentUser.IsSeeingAllGroups()) {
		// all existing actors
		return func(actorId string, day string) *actors.Actor {
			return actDict[actorId]
		}
	}
	// only visible actors (via current user groups)
	groupDict := make(map[int]bool)
	for _, grId := range m.CurrentUser.Groups {
		groupDict[grId] = true
	}
	return func(actorId string, day string) *actors.Actor {
		actor := actDict[actorId]
		if actor == nil {
			return nil
		}
		if groupDict[actor.Groups.ActiveGroupOnDate(day)] {
			return actor
		}
		return nil
	}
}

// genActorNameById returns a ActorNameById function: func(actorId string) string. Returned string (actor ref) is "" if actorId is not found
func (m *Manager) genActorNameById(visibleActorsOnly bool) clients.ActorNameById {
	getActorById := m.genActorById(visibleActorsOnly)
	return func(actorId string, day string) string {
		act := getActorById(actorId, day)
		if act == nil {
			return ""
		}
		return act.GetLabel()
	}
}

// genActorInfoById returns a ActorInfoById function: func(actorId string) []string which returns nil if actorId is not known, or [0] Actor Role [1] Actor Ref
func (m *Manager) genActorInfoById(visibleActorsOnly bool) clients.ActorInfoById {
	getActorById := m.genActorById(visibleActorsOnly)
	return func(actorId string, day string) []string {
		act := getActorById(actorId, day)
		if act == nil {
			return nil
		}
		return []string{act.Role, act.Ref}
	}
}

// SetGraphNameByClient sets given StatContext.GraphName to group by Group / Actors
func (m *Manager) SetGraphNameByGroup(sc *items.StatContext) {
	getActById := m.genActorById(true)
	getGroupById := m.GenGroupById()
	sc.GraphName = func(item *items.Item) []items.NamePct {
		defaultGroupName := getGroupById(0).Name
		gName := defaultGroupName
		res := []items.NamePct{}
		globPct := 1.0
		if sc.ShowTeam && len(item.Actors) > 0 {
			pct := 1.0 / float64(len(item.Actors))
			globPct = 0.0
			for _, actId := range item.Actors {
				actor := getActById(actId, item.Date)
				if actor == nil { // Skip unknown or not visible Actors
					continue
				}
				gName = defaultGroupName
				grp := getGroupById(actor.Groups.ActiveGroupOnDate(item.Date))
				if grp != nil {
					gName = grp.Name
				}
				res = append(res, items.NamePct{
					Name: gName + " : " + actor.GetLabel(),
					Pct:  pct,
				})
				globPct += pct
			}
		}
		res = append(res, items.NamePct{
			Name: gName,
			Pct:  globPct,
		})
		return res
	}
}

// ====================================================================================================
// Team Related Function ==============================================================================

// genIsTeamVisibleViaActors returns a IsTeamVisible function: func(ClientTeam) bool, which is true when current user is allowed to see clientteam (by actorId) related activity
func (m Manager) genIsTeamVisibleViaActors() (clients.IsTeamVisible, error) {
	if len(m.CurrentUser.Clients) == 0 {
		return func(clients.ClientTeam) bool { return true }, nil
	}

	actorVisible := make(map[clients.ClientTeam]bool)
	clts, err := m.GetCurrentUserClients()
	if err != nil {
		return nil, err
	}
	for _, client := range clts {
		allowedActors := m.Actors.GetActorsByClient(false, client.Name)
		for _, actor := range allowedActors {
			actorVisible[clients.ClientTeam{Client: client.Name, Team: strconv.Itoa(actor.Id)}] = true
			actorVisible[clients.ClientTeam{Client: client.Name, Team: actor.LastName}] = true
		}
	}
	return func(ct clients.ClientTeam) bool {
		return actorVisible[ct]
	}, nil
}

// genIsTeamVisible returns a IsTeamVisible function: func(ClientTeam) bool, which is true when current user is allowed to see clientteam related activity
func (m Manager) genIsTeamVisible() (clients.IsTeamVisible, error) {
	if len(m.CurrentUser.Clients) == 0 {
		return func(clients.ClientTeam) bool { return true }, nil
	}

	teamVisible := make(map[clients.ClientTeam]bool)
	clts, err := m.GetCurrentUserClients()
	if err != nil {
		return nil, err
	}
	for _, client := range clts {
		for _, team := range client.Teams {
			teamVisible[clients.ClientTeam{Client: client.Name, Team: team.Members}] = true
		}
	}
	return func(ct clients.ClientTeam) bool {
		return teamVisible[ct]
	}, nil
}
