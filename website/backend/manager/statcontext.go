package manager

import "github.com/lpuig/ewin/doe/website/backend/model/items"

func (m Manager) NewStatContext(freq string) (*items.StatContext, error) {
	statContext, err := items.NewStatContext(freq)
	if err != nil {
		return nil, err
	}

	isTeamVisible, err := m.genIsTeamVisibleViaActors()
	if err != nil {
		return nil, err
	}

	statContext.IsTeamVisible = isTeamVisible
	statContext.ClientByName = m.genGetClient()
	statContext.ActorNameById = m.genActorNameById(true)
	statContext.ShowTeam = !m.CurrentUser.Permissions["Review"]

	statContext.SetSerieTeamSiteConf()

	return statContext, nil
}
