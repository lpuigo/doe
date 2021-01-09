package clients

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/archives"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"path/filepath"
	"sync"
	"time"
)

type ClientsPersister struct {
	sync.RWMutex
	persister *persist.Persister

	clients []*ClientRecord
}

func NewClientsPersister(dir string) (*ClientsPersister, error) {
	wsp := &ClientsPersister{
		persister: persist.NewPersister("Clients", dir),
	}
	err := wsp.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	wsp.persister.SetPersistDelay(1 * time.Second)
	return wsp, nil
}

func (cp ClientsPersister) NbClients() int {
	return len(cp.clients)
}

// LoadDirectory loads all persisted Clients Records
func (cp *ClientsPersister) LoadDirectory() error {
	cp.Lock()
	defer cp.Unlock()

	cp.persister.Reinit()
	cp.clients = []*ClientRecord{}

	files, err := cp.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from ClientsPersister: %v", err)
	}

	for _, file := range files {
		ur, err := NewClientRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not create client from '%s': %v", filepath.Base(file), err)
		}
		err = cp.persister.Load(ur)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
		cp.clients = append(cp.clients, ur)
	}
	return nil
}

// Add adds the given ClientRecord to the ClientsPersister and return its (updated with new id) ClientRecord
func (cp *ClientsPersister) Add(nc *ClientRecord) *ClientRecord {
	cp.Lock()
	defer cp.Unlock()

	// give the record its new ID
	cp.persister.Add(nc)
	nc.Id = nc.GetId()
	cp.clients = append(cp.clients, nc)
	return nc
}

// Update updates the given ClientRecord
func (cp *ClientsPersister) Update(uc *ClientRecord) error {
	cp.RLock()
	defer cp.RUnlock()

	oc := cp.GetById(uc.Id)
	if oc == nil {
		return fmt.Errorf("client id not found")
	}
	oc.Client = uc.Client
	cp.persister.MarkDirty(oc)
	return nil
}

func (cp *ClientsPersister) findIndex(ur *ClientRecord) int {
	for i, rec := range cp.clients {
		if rec.GetId() == ur.GetId() {
			return i
		}
	}
	return -1
}

// Remove removes the given ClientRecord from the ClientsPersister (pertaining file is moved to deleted dir)
func (cp *ClientsPersister) Remove(ru *ClientRecord) error {
	cp.Lock()
	defer cp.Unlock()

	err := cp.persister.Remove(ru)
	if err != nil {
		return err
	}

	i := cp.findIndex(ru)
	copy(cp.clients[i:], cp.clients[i+1:])
	cp.clients[len(cp.clients)-1] = nil // or the zero value of T
	cp.clients = cp.clients[:len(cp.clients)-1]
	return nil
}

// GetById returns the ClientRecord with given Id (or nil if Id not found)
func (cp *ClientsPersister) GetById(id int) *ClientRecord {
	cp.RLock()
	defer cp.RUnlock()

	for _, cr := range cp.clients {
		if cr.Id == id {
			return cr
		}
	}
	return nil
}

// GetByRef returns the ClientRecord with given Name (or nil if Id not found)
func (cp *ClientsPersister) GetByName(name string) *ClientRecord {
	cp.RLock()
	defer cp.RUnlock()

	for _, cr := range cp.clients {
		if cr.Name == name {
			return cr
		}
	}
	return nil
}

func (cp *ClientsPersister) GetAllClients() []*Client {
	cp.RLock()
	defer cp.RUnlock()

	res := []*Client{}
	for _, cr := range cp.clients {
		res = append(res, cr.Client)
	}
	return res
}

func (cp *ClientsPersister) UpdateClients(updatedClients []*Client) error {
	for _, updClt := range updatedClients {
		ucl := NewClientRecordFromClient(updClt)
		if updClt.Id < 0 { // New Group, add it instead of update
			cp.Add(ucl)
			continue
		}
		err := cp.Update(ucl)
		if err != nil {
			fmt.Errorf("could not update group '%s' (id: %d)", ucl.Name, ucl.Id)
		}
	}
	return nil
}

func (cp *ClientsPersister) CalcPriceByClientArticleGetter() func(clientName, articleName string, qty int) (float64, error) {
	cp.RLock()
	defer cp.RUnlock()

	clts := make(map[string]map[string]bpu.Article)
	for _, cr := range cp.clients {
		articles := make(map[string]bpu.Article)
		for _, article := range cr.Client.GetOrangeArticles() {
			articles[article.Name] = *article
		}
		clts[cr.Client.Name] = articles
	}

	return func(clientName, articleName string, qty int) (float64, error) {
		if articleName == "" {
			articleName = "CEM42"
		}
		articles := clts[clientName]
		if articles == nil {
			return 0, fmt.Errorf("unknown client name: %s", clientName)
		}
		article := articles[articleName]
		if article.Name == "" {
			return 0, fmt.Errorf("unknown article name: %s", articleName)
		}
		return article.CalcPrice(qty), nil
	}
}

func (cp *ClientsPersister) GetAllSites() []archives.ArchivableRecord {
	cp.RLock()
	defer cp.RUnlock()

	archivableSites := make([]archives.ArchivableRecord, len(cp.clients))
	for i, site := range cp.clients {
		archivableSites[i] = site
	}
	return archivableSites
}

func (cp *ClientsPersister) GetName() string {
	return "Clients"
}
