package migrate

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/model/ripsites"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const (
	migrated_ripsite_dir string = `migrated_ripsites`
	ripsite_dir          string = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Ripsites`
	actors_dir           string = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Actors`
)

func Test_MigrateRipSiteWithActors(t *testing.T) {
	rp, err := ripsites.NewSitesPersit(ripsite_dir)
	if err != nil {
		t.Fatalf("NewSitesPersit returns unexpected :%s", err.Error())
	}

	err = rp.LoadDirectory()
	if err != nil {
		t.Fatalf("NewSitesPersit.LoadDirectory returns unexpected :%s", err.Error())
	}

	ap, err := actors.NewActorsPersister(actors_dir)
	if err != nil {
		t.Fatalf("NewActorsPersister returns unexpected :%s", err.Error())
	}

	err = ap.LoadDirectory()
	if err != nil {
		t.Fatalf("NewActorsPersister.LoadDirectory returns unexpected :%s", err.Error())
	}

	actorsDict := genActorsDict(t, ap)

	for _, ripsite := range rp.GetAll(func(s *ripsites.Site) bool {
		return true
	}) {
		migrateRipsite(t, ripsite, actorsDict)
	}
}

type ActorsDict map[string]string

func genActorsDict(t *testing.T, ap *actors.ActorsPersister) ActorsDict {
	res := make(ActorsDict)

	for _, actor := range ap.GetAllActors() {
		res[strings.Trim(strings.ToUpper(actor.LastName), " ")] = strconv.Itoa(actor.Id)
	}
	return res
}

// Match returns actors id if name match, empty string otherwise
func (ad ActorsDict) Match(name string) string {
	return ad[strings.Trim(strings.ToUpper(name), " ")]
}

func migrateRipsite(t *testing.T, ripsite *ripsites.SiteRecord, actorsDict ActorsDict) {
	// migrate pulling
	for _, pulling := range ripsite.Pullings {
		if tools.Empty(pulling.State.Team) {
			continue
		}
		names := strings.Split(pulling.State.Team, ",")
		actors := []string{}
		for _, name := range names {
			if id := actorsDict.Match(name); id != "" {
				actors = append(actors, id)
			}
		}
		t.Logf("%s: pulling %s => %s => %v", ripsite.Ref, pulling.CableName, pulling.State.Team, actors)
		pulling.State.Actors = actors
	}

	fname := filepath.Join(migrated_ripsite_dir, fmt.Sprintf("%06d.json", ripsite.Id))
	f, err := os.Create(fname)
	if err != nil {
		t.Fatalf("can not create file id %d: %s", ripsite.Id, err.Error())
	}
	defer f.Close()
	ripsite.Marshall(f)
}
