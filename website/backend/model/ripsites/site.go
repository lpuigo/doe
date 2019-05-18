package ripsites

import (
	"fmt"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
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
			dist := p.GetTotalDist()
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
