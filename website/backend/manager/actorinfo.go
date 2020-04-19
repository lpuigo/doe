package manager

import (
	"encoding/json"
	"fmt"
	"io"
)

func (m Manager) GetActorInfos(writer io.Writer) error {
	if !m.CurrentUser.HasPermissionHR() {
		return fmt.Errorf("%s ne dispose pas des droits suffisants", m.CurrentUser.Name)
	}
	clientsNames := m.GetCurrentUserClientsName()
	actors := m.Actors.GetActorsByClient(false, clientsNames...)
	actorinfos := m.ActorInfos.GetActorInfosByActors(actors)
	return json.NewEncoder(writer).Encode(actorinfos)
}
