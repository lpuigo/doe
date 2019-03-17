package persist

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

//TODO To be tested

type Container struct {
	sync.RWMutex
	persister *Persister

	name                 string
	newContaineeFromFile func(file string) (Recorder, error)
	records              []Recorder
}

func NewContainer(name, dir string) (*Container, error) {
	c := &Container{
		name:      name,
		persister: NewPersister(dir),
		newContaineeFromFile: func(file string) (recorder Recorder, e error) {
			return nil, fmt.Errorf("%s.newContaineeFromFile is not defined", name)
		},
	}
	err := c.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	c.persister.SetPersistDelay(1 * time.Second)
	return c, nil
}

func (c Container) Len() int {
	return len(c.records)
}

func (c *Container) LoadDirectory() error {
	c.Lock()
	defer c.Unlock()

	files, err := c.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from %s container: %v", c.name, err)
	}

	for _, file := range files {
		nr, err := c.newContaineeFromFile(file)
		if err != nil {
			return fmt.Errorf("could not create record from '%s': %v", filepath.Base(file), err)
		}
		c.persister.Load(nr)
		c.records = append(c.records, nr)
	}
	return nil
}

// GetById returns the Record with given Id (or nil if Id not found)
func (c *Container) GetById(id int) Recorder {
	for _, cr := range c.records {
		if cr.GetId() == id {
			return cr
		}
	}
	return nil
}

// Add adds the given ClientRecord to the ClientsPersister and return its (updated with new id) ClientRecord
func (c *Container) Add(nr Recorder) Recorder {
	c.Lock()
	defer c.Unlock()

	// give the record its new ID
	c.persister.Add(nr)
	c.records = append(c.records, nr)
	return nr
}

// Update updates the given ClientRecord
func (c *Container) Update(ur Recorder) error {
	c.RLock()
	defer c.RUnlock()

	or := c.GetById(ur.GetId())
	if or == nil {
		return fmt.Errorf("record id not found")
	}
	or = ur
	c.persister.MarkDirty(or)
	return nil
}
