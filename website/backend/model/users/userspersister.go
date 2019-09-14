package users

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"path/filepath"
	"sync"
	"time"
)

type UsersPersister struct {
	sync.RWMutex
	persister *persist.Persister

	users []*UserRecord
}

func NewUsersPersister(dir string) (*UsersPersister, error) {
	wsp := &UsersPersister{
		persister: persist.NewPersister("Users", dir),
	}
	err := wsp.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	wsp.persister.SetPersistDelay(1 * time.Second)
	return wsp, nil
}

func (up UsersPersister) NbUsers() int {
	return len(up.users)
}

// LoadDirectory loads all persisted Users Records
func (up *UsersPersister) LoadDirectory() error {
	up.Lock()
	defer up.Unlock()

	up.persister.Reinit()
	up.users = []*UserRecord{}

	files, err := up.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from UsersPersister: %v", err)
	}

	for _, file := range files {
		ur, err := NewUserRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not create user from '%s': %v", filepath.Base(file), err)
		}
		err = up.persister.Load(ur)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
		up.users = append(up.users, ur)
	}
	return nil
}

// Add adds the given UserRecord to the USersPersister and return its (updated with new id) UserRecord
func (up *UsersPersister) Add(nu *UserRecord) *UserRecord {
	up.Lock()
	defer up.Unlock()

	// give the record its new ID
	up.persister.Add(nu)
	nu.Id = nu.GetId()
	up.users = append(up.users, nu)
	return nu
}

// Update updates the given UserRecord
func (up *UsersPersister) Update(uu *UserRecord) error {
	up.RLock()
	defer up.RUnlock()

	ou := up.GetById(uu.Id)
	if ou == nil {
		return fmt.Errorf("user id not found")
	}
	ou.User = uu.User
	up.persister.MarkDirty(ou)
	return nil
}

func (up *UsersPersister) findIndex(ur *UserRecord) int {
	for i, rec := range up.users {
		if rec.GetId() == ur.GetId() {
			return i
		}
	}
	return -1
}

// Remove removes the given UserRecord from the UsersPersister (pertaining file is moved to deleted dir)
func (up *UsersPersister) Remove(ru *UserRecord) error {
	up.Lock()
	defer up.Unlock()

	err := up.persister.Remove(ru)
	if err != nil {
		return err
	}

	i := up.findIndex(ru)
	copy(up.users[i:], up.users[i+1:])
	up.users[len(up.users)-1] = nil // or the zero value of T
	up.users = up.users[:len(up.users)-1]
	return nil
}

// GetById returns the UserRecord with given Id (or nil if Id not found)
func (up *UsersPersister) GetById(id int) *UserRecord {
	for _, ur := range up.users {
		if ur.Id == id {
			return ur
		}
	}
	return nil
}

// GetByRef returns the UserRecord with given Name (or nil if Id not found)
func (up *UsersPersister) GetByName(name string) *UserRecord {
	for _, ur := range up.users {
		if ur.Name == name {
			return ur
		}
	}
	return nil
}
