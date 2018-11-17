package worksites

import (
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/persist"
)

type WorkSiteRecord struct {
	*persist.Record
	model.Worksite
}
