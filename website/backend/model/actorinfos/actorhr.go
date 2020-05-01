package actorinfos

import (
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
)

type ActorHr struct {
	*actors.Actor
	Info *ActorInfo
}
