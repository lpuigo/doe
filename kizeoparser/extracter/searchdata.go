package extracter

import (
	"fmt"
	"github.com/lpuig/ewin/doe/kizeoparser/api"
	"sort"
)

func SortSearchDatasBySroRef(datas []*api.SearchData) {
	sort.Slice(datas, func(i, j int) bool {
		return datas[i].SummarySubtitle < datas[j].SummarySubtitle
	})
}

func getGMAPUrl(rec *api.SearchData) string {
	return fmt.Sprintf("https://maps.google.com/maps?q=%s,%s", rec.Geoloc.Lat, rec.Geoloc.Long)
	// https://maps.google.com/maps?q=43.63292314,%205.58770894
}
