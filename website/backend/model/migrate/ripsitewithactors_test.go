package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/model/ripsites"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
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

type ActorsDict map[string][]*actors.Actor

func genActorsDict(t *testing.T, ap *actors.ActorsPersister) ActorsDict {
	res := make(ActorsDict)

	for _, actor := range ap.GetAllActors() {
		luName := strings.Trim(strings.ToUpper(actor.LastName), " ")
		if res[luName] == nil {
			res[luName] = []*actors.Actor{actor}
			continue
		}
		res[luName] = append(res[luName], actor)
	}
	return res
}

// Match returns actors id if name match, empty string otherwise
func (ad ActorsDict) Match(name string, state *ripsites.State) string {
	luName := strings.Trim(strings.ToUpper(name), " ")
	matchActor := ad[luName]
	if matchActor == nil {
		return ""
	}
	if len(matchActor) == 1 {
		return strconv.Itoa(matchActor[0].Id)
	}
	checkDate := state.DateEnd
	if checkDate == "" {
		checkDate = state.DateStart
	}
	for _, actor := range matchActor {
		if actor.Period.Begin > checkDate || (actor.Period.End != "" && actor.Period.End < checkDate) {
			continue
		}
		return strconv.Itoa(actor.Id)
	}
	return ""
}

func migrateState(t *testing.T, ref1, ref2 string, state *ripsites.State, actorsDict ActorsDict) {
	names := strings.Split(state.Team, ",")
	actors := []string{}
	for _, name := range names {
		if id := actorsDict.Match(name, state); id != "" {
			actors = append(actors, id)
		}
	}
	t.Logf("%s.%s: %s => %v", ref1, ref2, state.Team, actors)
	state.Actors = actors
}

func initState(state *ripsites.State) {
	state.Actors = []string{}
}

func migrateRipsite(t *testing.T, ripsite *ripsites.SiteRecord, actorsDict ActorsDict) {
	// migrate measurement
	for _, pulling := range ripsite.Pullings {
		initState(&pulling.State)
		for i, _ := range pulling.Chuncks {
			initState(&pulling.Chuncks[i].State)
		}
		if tools.Empty(pulling.State.Team) {
			continue
		}
		migrateState(t, "pull "+ripsite.Ref, pulling.Chuncks[0].EndingNodeName, &pulling.State, actorsDict)
	}

	// migrate Junctions
	for _, junction := range ripsite.Junctions {
		initState(&junction.State)
		for i, _ := range junction.Operations {
			initState(&junction.Operations[i].State)
		}
		if tools.Empty(junction.State.Team) {
			continue
		}
		migrateState(t, "junc "+ripsite.Ref, junction.NodeName, &junction.State, actorsDict)
	}

	// migrate Junctions
	for _, measurement := range ripsite.Measurements {
		initState(&measurement.State)
		if tools.Empty(measurement.State.Team) {
			continue
		}
		migrateState(t, "meas "+ripsite.Ref, measurement.DestNodeName, &measurement.State, actorsDict)
	}

	fname := filepath.Join(migrated_ripsite_dir, fmt.Sprintf("%06d.json", ripsite.Id))
	f, err := os.Create(fname)
	if err != nil {
		t.Fatalf("can not create file id %d: %s", ripsite.Id, err.Error())
	}
	defer f.Close()
	ripsite.Marshall(f)
}
