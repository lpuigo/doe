package ripsites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"strconv"
	"strings"
)

type Site struct {
	Id         int
	Client     string
	Ref        string
	Manager    string
	OrderDate  string
	UpdateDate string
	Status     string
	Comment    string
	Nodes      map[string]*Node
	Troncons   map[string]*Troncon

	Pullings     []*Pulling
	Junctions    []*Junction
	Measurements []*Measurement
}

type IsSiteVisible func(s *Site) bool

func (s *Site) GetInfo() *fm.RipsiteInfo {
	rsi := fm.NewBERipsiteInfo()

	rsi.Id = s.Id
	rsi.Client = s.Client
	rsi.Ref = s.Ref
	rsi.Manager = s.Manager
	rsi.OrderDate = s.OrderDate
	rsi.UpdateDate = s.UpdateDate
	rsi.Status = s.Status
	rsi.Comment = s.Comment

	rsi.NbPulling, rsi.NbPullingBlocked, rsi.NbPullingDone = s.GetPullingNumbers()
	rsi.NbJunction, rsi.NbJunctionBlocked, rsi.NbJunctionDone = s.GetJunctionNumbers()
	rsi.NbMeasurement, rsi.NbMeasurementBlocked, rsi.NbMeasurementDone = s.GetMeasurementNumbers()

	var searchBuilder strings.Builder
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Client", s.Client)
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Ref", s.Ref)
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Manager", s.Manager)
	fmt.Fprintf(&searchBuilder, "%s:%s,", "OrderDate", s.OrderDate)
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Comment", s.Comment)
	for _, node := range s.Nodes {
		fmt.Fprintf(&searchBuilder, "nRef:%s,", node.Ref)
		fmt.Fprintf(&searchBuilder, "nName:%s,", node.Name)
		fmt.Fprintf(&searchBuilder, "nAddr:%s,", node.Address)
	}
	rsi.Search = searchBuilder.String()

	return rsi
}

// GetPullingNumbers returns total, blocked and done number of Pullings
func (s *Site) GetPullingNumbers() (total, blocked, done int) {
	for _, p := range s.Pullings {
		t, b, d := p.State.TodoBlockedDone()
		if t {
			dist := p.GetTotalAggrDist()
			total += dist
			if b {
				blocked += dist
			}
			if d {
				done += dist
			}
		}
	}
	return
}

// GetJunctionNumbers returns total, blocked and done number of Junctions
func (s *Site) GetJunctionNumbers() (total, blocked, done int) {
	for _, j := range s.Junctions {
		t, b, d := j.State.TodoBlockedDone()
		if t {
			nbFiber := j.GetNbFiber()
			total += nbFiber
			if b {
				blocked += nbFiber
			}
			if d {
				done += nbFiber
			}
		}
	}
	return
}

// GetMeasurementNumbers returns total, blocked and done number of Measurements
func (s *Site) GetMeasurementNumbers() (total, blocked, done int) {
	for _, m := range s.Measurements {
		t, b, d := m.State.TodoBlockedDone()
		if t {
			nbMeas := m.GetNbMeas()
			total += nbMeas
			if b {
				blocked += nbMeas
			}
			if d {
				done += nbMeas
			}
		}
	}
	return
}

func (s *Site) Itemize(bpu *bpu.Bpu) ([]*items.Item, error) {
	res := []*items.Item{}
	pullItems, _ := s.itemizePullings(bpu)
	junctItems, _ := s.itemizeJunctions(bpu)
	measItems, _ := s.itemizeMeasurements(bpu)

	res = append(res, pullItems...)
	res = append(res, junctItems...)
	res = append(res, measItems...)

	return res, nil
}

const (
	activityPulling     string = "Tirage"
	activityJunction    string = "Racco"
	activityMeasurement string = "Mesures"

	catPullUnderground string = activityPulling + " Souterain"
	catPullAerial      string = activityPulling + " Aérien"
	catPullBuilding    string = activityPulling + " Façade"
)

func (s *Site) itemizePullings(bpu *bpu.Bpu) ([]*items.Item, error) {
	res := []*items.Item{}

	pullingArticles := bpu.GetCategoryArticles(activityPulling)

	for _, pulling := range s.Pullings {
		cableSize, err := getCableSize(pulling.CableName)
		if err != nil {
			return nil, err
		}
		todo, done := pulling.State.GetTodoDone()
		l, u, a, b := pulling.GetTotalDists()
		// Item for underground cable pulling
		if l+u > 0 {
			article, err := pullingArticles.GetArticleFor(catPullUnderground, cableSize)
			if err != nil {
				return nil, fmt.Errorf("can not define Underground Pulling Item: %s", err.Error())
			}
			res = append(res, items.NewItem(
				activityPulling,
				pulling.Chuncks[0].TronconName,
				fmt.Sprintf("Tirage %s (%dml)", pulling.CableName, l+u),
				pulling.State.DateEnd,
				pulling.State.Team,
				article,
				l+u,
				l+u,
				todo,
				done,
			))
		}

		// Item for aerial cable pulling
		if a+b > 0 {
			article, err := pullingArticles.GetArticleFor(catPullAerial, cableSize)
			if err != nil {
				return nil, fmt.Errorf("can not define Aerial Pulling Item: %s", err.Error())
			}
			res = append(res, items.NewItem(
				activityPulling,
				pulling.Chuncks[0].TronconName,
				fmt.Sprintf("Tirage %s (%dml)", pulling.CableName, a+b),
				pulling.State.DateEnd,
				pulling.State.Team,
				article,
				a+b,
				a+b,
				todo,
				done,
			))
		}

		// Item for building cable pulling
		if b > 0 {
			article, err := pullingArticles.GetArticleFor(catPullBuilding, cableSize)
			if err != nil {
				return nil, fmt.Errorf("can not define Building Pulling Item: %s", err.Error())
			}
			res = append(res, items.NewItem(
				activityPulling,
				pulling.Chuncks[0].TronconName,
				fmt.Sprintf("Tirage %s (%dml)", pulling.CableName, b),
				pulling.State.DateEnd,
				pulling.State.Team,
				article,
				b,
				b,
				todo,
				done,
			))
		}

	}

	return res, nil
}

func (s *Site) itemizeJunctions(bpu *bpu.Bpu) ([]*items.Item, error) {
	res := []*items.Item{}
	return res, nil
}

func (s *Site) itemizeMeasurements(bpu *bpu.Bpu) ([]*items.Item, error) {
	res := []*items.Item{}
	return res, nil
}

func getCableSize(cableName string) (int, error) {
	parts := strings.Split(cableName, "_")
	if len(parts) < 2 {
		return 0, fmt.Errorf("misformatted cable type '%': can not detect _nnFO_ chunk", cableName)
	}
	size, e := strconv.ParseInt(strings.TrimSuffix(parts[1], "FO"), 10, 64)
	if e != nil {
		return 0, fmt.Errorf("misformatted cable type: can not get number of fiber in '%'", parts[1])
	}
	return int(size), nil
}
