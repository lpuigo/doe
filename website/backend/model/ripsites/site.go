package ripsites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
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

func (s *Site) GetRef() string {
	return s.Ref
}

func (s *Site) GetClient() string {
	return s.Client
}

func (s *Site) GetType() string {
	return "ripsite"
}

func (s *Site) GetUpdateDate() string {
	return s.UpdateDate
}

func (s *Site) GetInfo(clientByName clients.ClientByName) *fm.RipsiteInfo {
	rsi := fm.NewBERipsiteInfo()

	rsi.Id = s.Id
	rsi.Client = s.Client
	rsi.Ref = s.Ref
	rsi.Manager = s.Manager
	rsi.OrderDate = s.OrderDate
	rsi.UpdateDate = s.UpdateDate
	rsi.Status = s.Status
	rsi.Comment = s.Comment

	var searchBuilder strings.Builder
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Client", strings.ToUpper(s.Client))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Ref", strings.ToUpper(s.Ref))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Manager", strings.ToUpper(s.Manager))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "OrderDate", strings.ToUpper(s.OrderDate))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Comment", strings.ToUpper(s.Comment))
	for _, node := range s.Nodes {
		fmt.Fprintf(&searchBuilder, "nRef:%s,", strings.ToUpper(node.Ref))
		fmt.Fprintf(&searchBuilder, "nName:%s,", strings.ToUpper(node.Name))
		fmt.Fprintf(&searchBuilder, "nAddr:%s,", strings.ToUpper(node.Address))
	}
	rsi.Search = searchBuilder.String()

	client := clientByName(s.Client)
	if client == nil {

	}
	itms, err := s.Itemize(client.Bpu, false)
	if err != nil {
		rsi.NbPoints, rsi.NbPointsBlocked, rsi.NbPointsDone = 0, 0, 0
		rsi.NbPulling, rsi.NbPullingBlocked, rsi.NbPullingDone = 0, 0, 0
		rsi.NbJunction, rsi.NbJunctionBlocked, rsi.NbJunctionDone = 0, 0, 0
		return rsi
	}
	points, pulling, junction := s.getPointsNumbers(itms)
	rsi.NbPoints, rsi.NbPointsBlocked, rsi.NbPointsDone = points.total, points.blocked, points.done
	rsi.NbPulling, rsi.NbPullingBlocked, rsi.NbPullingDone = pulling.total, pulling.blocked, pulling.done
	rsi.NbJunction, rsi.NbJunctionBlocked, rsi.NbJunctionDone = junction.total, junction.blocked, junction.done

	//pullingComplete := rsi.NbPullingBlocked + rsi.NbPullingDone == rsi.NbPulling
	//junctionComplete := rsi.NbJunctionBlocked + rsi.NbJunctionDone == rsi.NbJunction
	//
	//if pullingComplete && junctionComplete && rsi.NbPointsBlocked + rsi.NbPointsDone != rsi.NbPoints {
	//	rsi.NbPoints = rsi.NbPointsBlocked + rsi.NbPointsDone
	//}
	return rsi
}

type counts struct {
	total   int
	blocked int
	done    int
}

// GetMeasurementNumbers returns total, blocked and done number of Measurements
func (s *Site) getPointsNumbers(itms []*items.Item) (points, pulling, junction counts) {
	const factor = 1000
	pullingComplete := true
	junctionComplete := true
	for _, itm := range itms {
		if !itm.Todo {
			continue
		}
		pts := int(itm.Work() * factor)
		switch itm.Activity {
		case activityPulling:
			pulling.total += pts
			if itm.Done {
				pulling.done += pts
			}
			if itm.Blocked {
				pulling.blocked += pts
			}
			if !itm.Blocked && !itm.Done {
				pullingComplete = false
			}
		case activityJunction, activityMeasurement:
			junction.total += pts
			if itm.Done {
				junction.done += pts
			}
			if itm.Blocked {
				junction.blocked += pts
			}
			if !itm.Blocked && !itm.Done {
				junctionComplete = false
			}
		}
	}
	pulling.total /= factor
	pulling.blocked /= factor
	pulling.done /= factor
	if pullingComplete && pulling.total != pulling.done+pulling.blocked {
		pulling.total = pulling.done + pulling.blocked
	}

	junction.total /= factor
	junction.blocked /= factor
	junction.done /= factor
	if junctionComplete && junction.total != junction.done+junction.blocked {
		junction.total = junction.done + junction.blocked
	}
	points.total = pulling.total + junction.total
	points.blocked = pulling.blocked + junction.blocked
	points.done = pulling.done + junction.done
	return
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

func (s *Site) Itemize(bpu *bpu.Bpu, doneOnly bool) ([]*items.Item, error) {
	res := []*items.Item{}
	pullItems, err := s.itemizePullings(bpu, doneOnly)
	if err != nil {
		return nil, err
	}
	junctItems, err := s.itemizeJunctions(bpu, doneOnly)
	if err != nil {
		return nil, err
	}
	measItems, err := s.itemizeMeasurements(bpu, doneOnly)
	if err != nil {
		return nil, err
	}

	res = append(res, pullItems...)
	res = append(res, junctItems...)
	res = append(res, measItems...)

	return res, nil
}

const (
	activityPulling     string = "Tirage"
	activityJunction    string = "Racco"
	activityMeasurement string = "Mesures"

	catPullUnderground string = activityPulling + " Souterrain"
	catPullAerial      string = activityPulling + " Aérien"
	catPullBuilding    string = activityPulling + " Façade"

	catJuncPM  string = "PM"
	catJuncBPE string = "BPE"
	catJuncPBO string = "PBO"

	catMeasurement string = "Mesure"
)

func (s *Site) itemizePullings(currentBpu *bpu.Bpu, doneOnly bool) ([]*items.Item, error) {
	res := []*items.Item{}

	pullingArticles := currentBpu.GetCategoryArticles(activityPulling)

	for _, pulling := range s.Pullings {
		cableSize, err := getCableSize(pulling.CableName)
		if err != nil {
			return nil, err
		}
		todo, done, blocked := pulling.State.GetTodoDoneBlocked()
		if doneOnly && !done {
			continue
		}

		l, u, a, b := pulling.GetTotalDists()
		// Item for underground cable pulling
		if l+u > 0 {
			article, err := pullingArticles.GetArticleFor(catPullUnderground, cableSize)
			if err != nil {
				return nil, fmt.Errorf("can not define Underground Pulling Item: %s", err.Error())
			}
			item := items.NewItem(
				s.Client, s.Ref,
				activityPulling,
				pulling.Chuncks[0].TronconName,
				fmt.Sprintf("Tirage %s (%dml)", pulling.CableName, l+u),
				pulling.State.DateEnd,
				"",
				article,
				l+u,
				l+u,
				todo,
				done,
				blocked,
				false,
			)
			item.StartDate = pulling.State.DateStart
			item.Actors = pulling.State.Actors
			res = append(res, item)
		}

		// Item for aerial cable pulling
		if a+b > 0 {
			article, err := pullingArticles.GetArticleFor(catPullAerial, cableSize)
			if err != nil {
				return nil, fmt.Errorf("can not define Aerial Pulling Item: %s", err.Error())
			}
			item := items.NewItem(
				s.Client, s.Ref,
				activityPulling,
				pulling.Chuncks[0].TronconName,
				fmt.Sprintf("Tirage %s (%dml)", pulling.CableName, a+b),
				pulling.State.DateEnd,
				"",
				article,
				a+b,
				a+b,
				todo,
				done,
				blocked,
				false,
			)
			item.StartDate = pulling.State.DateStart
			item.Actors = pulling.State.Actors
			res = append(res, item)
		}

		// Item for building cable pulling
		if b > 0 {
			article, err := pullingArticles.GetArticleFor(catPullBuilding, cableSize)
			if err != nil {
				return nil, fmt.Errorf("can not define Building Pulling Item: %s", err.Error())
			}
			item := items.NewItem(
				s.Client, s.Ref,
				activityPulling,
				pulling.Chuncks[0].TronconName,
				fmt.Sprintf("Tirage %s (%dml)", pulling.CableName, b),
				pulling.State.DateEnd,
				"",
				article,
				b,
				b,
				todo,
				done,
				blocked,
				false,
			)
			item.StartDate = pulling.State.DateStart
			item.Actors = pulling.State.Actors
			res = append(res, item)
		}

	}

	return res, nil
}

func (s *Site) itemizeJunctions(currentBpu *bpu.Bpu, doneOnly bool) ([]*items.Item, error) {
	res := []*items.Item{}
	junctionArticles := currentBpu.GetCategoryArticles(activityJunction)

	var e error

	for _, junction := range s.Junctions {
		todo, done, blocked := junction.State.GetTodoDoneBlocked()
		if doneOnly && !done {
			continue
		}

		node, nodeFound := s.Nodes[junction.NodeName]
		if !nodeFound {
			return nil, fmt.Errorf("unknow node '%s'", junction.NodeName)
		}
		if !currentBpu.IsBoxDefined(node.Type, node.BoxType) {
			return nil, fmt.Errorf("unknow Box Type '%s' for '%s'", node.BoxType, node.Type)
		}
		_, nbSplice := junction.GetNbFiberSplice()
		tronconIn, tronconInFound := s.Troncons[node.TronconInName]
		var boxSize int
		if tronconInFound {
			boxSize = tronconIn.Size
		}

		var mainArticle, optArticle *bpu.Article
		var qty1, qty2 int

		switch strings.ToUpper(node.Type) {
		case catJuncPM:
			pmArticles := junctionArticles.GetArticles(catJuncPM)
			mainArticle, optArticle = pmArticles[0], pmArticles[1]

			qty1 = nbSplice / mainArticle.Unit
			// check for missing modules
			qty2 = 0
			if qty1*mainArticle.Unit < nbSplice {
				qty1++
				nbMissingSplice := qty1*mainArticle.Unit - nbSplice
				qty2 = nbMissingSplice / optArticle.Unit
			}
		case catJuncBPE, catJuncPBO:
			mainArticle, optArticle, e = getJunctionBoxArticles(currentBpu, activityJunction, node.Type, node.BoxType)
			qty1, qty2 = 1, nbSplice
			if e != nil {
				return nil, fmt.Errorf("shit hit the fence: %s", e.Error())
			}
		default:
			return nil, fmt.Errorf("unexpected box category '%s'", node.Type)
		}

		info := fmt.Sprintf("Install. %s: %s", node.Type, node.BoxType)
		if boxSize > 0 {
			info += fmt.Sprintf(" (%dFO)", boxSize)
		}

		item := items.NewItem(s.Client, s.Ref, activityJunction, junction.NodeName, info, junction.State.DateEnd, "", mainArticle, qty1, qty1, todo, done, blocked, false)
		item.StartDate = junction.State.DateStart
		item.Actors = junction.State.Actors
		res = append(res, item)

		if optArticle != nil {
			item2 := items.NewItem(s.Client, s.Ref, activityJunction, junction.NodeName, info, junction.State.DateEnd, "", optArticle, qty2, qty2, todo, done, blocked, false)
			item2.StartDate = junction.State.DateStart
			item2.Actors = junction.State.Actors
			res = append(res, item2)
		}
	}

	return res, nil
}

func (s *Site) itemizeMeasurements(currentBpu *bpu.Bpu, doneOnly bool) ([]*items.Item, error) {
	res := []*items.Item{}
	measurementArticles := currentBpu.GetCategoryArticles(activityMeasurement)

	mainArticle, err := measurementArticles.GetArticleFor(catMeasurement, 1)
	if err != nil {
		return nil, err
	}
	qty1 := 1

	for _, measurement := range s.Measurements {
		todo, done, blocked := measurement.State.GetTodoDoneBlocked()
		if doneOnly && !done {
			continue
		}
		//actors := []string{}
		//for _, actorId := range measurement.State.Actors {
		//	actors = append(actors, actorById(actorId))
		//}
		//actorsString := strings.Join(actors, ", ")

		qty2 := measurement.NbFiber
		info := fmt.Sprintf("Mesure %d fibres - %d epissures", qty2, measurement.NbSplice())
		item := items.NewItem(s.Client, s.Ref, activityMeasurement, measurement.DestNodeName, info, measurement.State.DateEnd, "", mainArticle, qty1, qty2, todo, done, blocked, false)
		item.StartDate = measurement.State.DateStart
		item.Actors = measurement.State.Actors
		res = append(res, item)
	}
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

// getJunctionBoxArticles returns Article applicable for given Bpe or Pbo type
func getJunctionBoxArticles(currentBpu *bpu.Bpu, activity, category, boxType string) (boxArticle, spliceArticle *bpu.Article, err error) {
	// box lookup
	box := currentBpu.GetBox(category, boxType)
	if box == nil {
		err = fmt.Errorf("unknow box type '%s' for category '%s'", boxType, category)
		return
	}

	catArticles := currentBpu.GetCategoryArticles(activity)
	if catArticles == nil {
		err = fmt.Errorf("unknow activity '%s'", activity)
		return
	}

	switch strings.ToUpper(category) {
	case catJuncBPE:
		boxArticle, err = catArticles.GetArticleFor(category, box.Size)
		if err != nil {
			return
		}
		spliceArticle, err = catArticles.GetArticleFor(category+" Splice", box.Size)
		if err != nil {
			return
		}
	case catJuncPBO:
		boxArticle, err = catArticles.GetArticleFor(category+" "+box.Usage, box.Size)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("category '%s' is not handled", category)
		return
	}
	return
}

type IsSiteVisible func(s *Site) bool

const (
	RipStatSerieWork  string = "Work"
	RipStatSeriePrice string = "Price"
)
