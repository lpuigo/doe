package teamproductivitychart

import (
	"github.com/gopherjs/gopherjs/js"
)

//type serie struct {
//	*js.Object
//	Name string `js:"name"`
//	Color string `js:"color"`
//	Data []int `js:"data"`
//}
//
func newSerie(style, name, color string, pointpadding float64, data []int) js.M {
	//s := &serie{Object: tools.O()}
	//s.Name = name
	//s.Color = color
	//s.Data = data

	return js.M{
		"type":         style,
		"name":         name,
		"color":        color,
		"data":         data,
		"pointPadding": pointpadding,
	}
}
