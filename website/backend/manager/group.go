package manager

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/model/groups"
	"io"
)

func (m Manager) GetGroups(writer io.Writer) error {
	grps := m.Groups.GetGroups()
	return json.NewEncoder(writer).Encode(grps)
}

func (m Manager) UpdateGroups(updatedGroups []*groups.Group) error {
	return m.Groups.UpdateGroups(updatedGroups)
}
