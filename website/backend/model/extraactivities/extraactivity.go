package extraactivities

import (
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
)

type ExtraActivity struct {
	Name           string
	State          string
	NbPoints       float64
	Income         float64
	Date           string
	AttachmentDate string
	Actors         []string
	Comment        string
}

func (ea *ExtraActivity) MakeArticle() *bpu.Article {
	article := bpu.NewArticle()
	article.Name = "Activité Supplémentaire"
	article.Price = ea.Income
	article.Work = ea.NbPoints
	return article
}

func Itemize(eas []*ExtraActivity, site items.ItemizableSite, doneOnly bool) []*items.Item {
	res := []*items.Item{}
	for _, ea := range eas {
		done := len(ea.Actors) > 0 && ea.Date != ""
		billed := done && ea.AttachmentDate != ""
		if doneOnly && !done {
			continue
		}
		item := items.NewItem(
			site.GetClient(),
			site.GetRef(),
			"Extra",
			ea.Name,
			"",
			ea.Date,
			"",
			ea.MakeArticle(),
			1,
			1,
			true,
			done,
			false,
			billed,
		)
		item.StartDate = ea.Date
		item.Actors = ea.Actors[:]
		item.Work()

		res = append(res, item)
	}
	return res
}

func GetUpdateDate(eas []*ExtraActivity) string {
	updatedate := date.TimeJSMinDate
	for _, ea := range eas {
		if ea.Date != "" && ea.Date > updatedate {
			updatedate = ea.Date
		}
	}
	return updatedate
}
