package ripteamproductivitychart

import (
	"github.com/gopherjs/gopherjs/js"
)

func newSerie(style, name, color, stack string, pointpadding float64, data map[string][]float64) []js.M {
	//s := &serie{Object: tools.O()}
	//s.Name = name
	//s.Color = color
	//s.Data = data
	res := []js.M{}
	for key, value := range data {
		res = append(res, js.M{
			"type":         style,
			"name":         name + " " + key,
			"color":        color,
			"data":         value,
			"pointPadding": pointpadding,
			"stack":        stack,
		})
	}
	return res
}
