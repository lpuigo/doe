package actor

import "github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"

func GetDefaultInfoTraining() []string {
	return []string{
		actorconst.InfoTrainingBirth,
		actorconst.InfoTrainingMedic,
		actorconst.InfoTrainingProBTP,
		actorconst.InfoTrainingAIPR,
		actorconst.InfoTrainingElec,
		actorconst.InfoTrainingCacesNacelle,
		actorconst.InfoTrainingCacesHauteur,
		actorconst.InfoTrainingCacesChantier,
		actorconst.InfoTrainingCacesGrue,
	}
}
