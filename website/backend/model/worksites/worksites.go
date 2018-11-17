package worksites

import (
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"time"
)

type WorkSitesPersister struct {
	*persist.Persister
}

func NewWorkSitesPersist(dir string) *WorkSitesPersister {
	return &WorkSitesPersister{
		Persister: persist.NewPersister(dir, 2*time.Second),
	}
}

func (wsp WorkSitesPersister) GetAll() {

}

// Add adds the given WorkSiteRecord to the WorkSitesPersister and return its (new) id
func (wsp *WorkSitesPersister) Add(ws *WorkSiteRecord) int {
	// give the record its new ID
	id := wsp.Persister.Add(ws)
	return id
}
