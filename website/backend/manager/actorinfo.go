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
	actors := m.GetCurrentUserActors()
	actorinfos := m.ActorInfos.GetActorInfosByActors(actors)
	return json.NewEncoder(writer).Encode(actorinfos)
}
