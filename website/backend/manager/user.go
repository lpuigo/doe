package manager

import (
	"encoding/json"
	"fmt"
	"io"

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
