package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/users"
)

func (m Manager) GetUsers(writer io.Writer) error {
	usrs := m.Users.GetUsers()
	return json.NewEncoder(writer).Encode(usrs)
}

func (m Manager) UpdateUsers(updatedUsers []*users.User) error {
	return m.Users.UpdateUsers(updatedUsers)
}

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

// genActorById returns a ActorById function: func(actorId string) *actors.Actor. Returned *actors.Actor is nil if actorId is not found
func (m *Manager) genActorById() func(actorId string) *actors.Actor {
	return func(actorId string) *actors.Actor {
		var ar *actors.ActorRecord
		if actId, err := strconv.Atoi(actorId); err == nil {
			ar = m.Actors.GetById(actId)
		} else {
			ar = m.Actors.GetByRef(actorId)
		}
		if ar == nil {
			return nil
		}
		return ar.Actor
	}
}

// genActorNameById returns a ActorNameById function: func(actorId string) string. Returned string (actor ref) is "" if actorId is not found
func (m *Manager) genActorNameById() clients.ActorNameById {
	getActorById := m.genActorById()
	return func(actorId string) string {
		act := getActorById(actorId)
		if act == nil {
			return ""
		}
		return act.GetLabel()
	}
}

// genActorInfoById returns a ActorInfoById function: func(actorId string) []string which returns nil if actorId is not known, or [0] Actor Role [1] Actor Ref
func (m *Manager) genActorInfoById() clients.ActorInfoById {
	getActorById := m.genActorById()
	return func(actorId string) []string {
		act := getActorById(actorId)
		if act == nil {
			return nil
		}
		return []string{act.Role, act.Ref}
	}
}

// GetCurrentUserClientsName returns Clients' names visible by current user (if user has no client, returns empty list)
func (m Manager) GetCurrentUserClientsName() []string {
	if m.CurrentUser == nil {
		return nil
	}
	if len(m.CurrentUser.Clients) > 0 {
		return m.CurrentUser.Clients
	}
	return []string{}
}

// GetCurrentUserClients returns slice of Clients visible by current user
//
// Rules:
//
// - if current user has attached groups, visible Clients are identified by User's groups Clients list
//
// - otherwise, returns User's Client list
func (m Manager) GetCurrentUserClients() ([]*clients.Client, error) {
	clientDict := make(map[string]*clients.Client)
	if m.CurrentUser == nil {
		return nil, nil
	}
	getClientByName := m.genGetClient()
	// Set Clients via User's Groups
	if len(m.CurrentUser.Groups) != 0 {
		getGroupById := m.GenGroupById()
		for _, groupId := range m.CurrentUser.Groups {
			group := getGroupById(groupId)
			if group == nil {
				continue
			}
			for _, clientName := range group.Clients {
				client := getClientByName(clientName)
				if client == nil {
					continue
				}
				clientDict[clientName] = client
			}
		}
		res := make([]*clients.Client, len(clientDict))
		i := 0
		for _, client := range clientDict {
			res[i] = client
			i++
		}
		return res, nil
	}
	// Set Clients via Users's Clients
	if len(m.CurrentUser.Clients) == 0 {
		return m.Clients.GetAllClients(), nil
	}
	res := []*clients.Client{}
	for _, clientName := range m.CurrentUser.Clients {
		client := getClientByName(clientName)
		if client == nil {
			return nil, fmt.Errorf("could not retrieve client '%s' info", clientName)
		}
		res = append(res, client)
	}
	return res, nil
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

// genIsActorVisible returns a IsTeamVisible function: func(ClientTeam) bool, which is true when current user is allowed to see clientteam (by actorId) related activity
func (m Manager) genIsActorVisible() (clients.IsTeamVisible, error) {
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
