package manager

import "github.com/lpuig/ewin/doe/website/backend/model/items"

func (m Manager) NewStatContext(freq string) (*items.StatContext, error) {
	statContext, err := items.NewStatContext(freq)
	if err != nil {
		return nil, err
	}

	isActorVisible, err := m.genIsActorVisible()
	if err != nil {
		return nil, err
	}

	statContext.IsTeamVisible = isActorVisible
	statContext.ClientByName = m.genGetClient()
	statContext.ActorById = m.genActorById()
	statContext.ShowTeam = !m.CurrentUser.Permissions["Review"]

	statContext.SetSerieTeamSiteConf()

	return statContext, nil
}
