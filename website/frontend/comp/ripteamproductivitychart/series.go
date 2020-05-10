package ripteamproductivitychart

import (
	"github.com/gopherjs/gopherjs/js"
	"strconv"
)

type ColorMap struct {
	HueStart   int
	HueEnd     int
	Light      int
	Saturation int
}

type siteColors map[string]string

type SiteColorMap map[string]siteColors

func (s SiteColorMap) GetWorkColor(site string) string {
	return s["Price"][site]
}

func SetColor(sites map[string]bool, workCM, priceCM ColorMap) SiteColorMap {
	res := SiteColorMap{}
	wdHue := (workCM.HueEnd - workCM.HueStart) / len(sites)
	wHue := workCM.HueStart
	pdHue := (priceCM.HueEnd - priceCM.HueStart) / len(sites)
	pHue := priceCM.HueStart
	wSC := siteColors{}
	pSC := siteColors{}
	wSC["Total"] = hsl(wHue-30.0, workCM.Light, workCM.Saturation)
	pSC["Total"] = hsl(pHue-30.0, priceCM.Light, priceCM.Saturation)
	for site, _ := range sites {
		wSC[site] = hsl(wHue, workCM.Light, workCM.Saturation)
		pSC[site] = hsl(pHue, priceCM.Light, priceCM.Saturation)
		wHue += wdHue
		pHue += pdHue
	}
	res["Work"] = wSC
	res["Price"] = pSC
	return res
}

func newSerie(style, dash, name, stack, tooltipPrefix, tooltipSuffix string, colormap siteColors, yAxis int, pointpadding float64, data map[string][]float64) []js.M {
	//s := &serie{Object: tools.O()}
	//s.Name = name
	//s.Color = color
	//s.Data = data
	res := []js.M{}
	for key, value := range data {
		res = append(res, js.M{
			//"name":         name + " " + key,
			"name":         key,
			"type":         style,
			"dashStyle":    dash,
			"color":        colormap[key],
			"yAxis":        yAxis,
			"data":         value,
			"pointPadding": pointpadding,
			"stack":        stack,
			"tooltip": js.M{
				"valuePrefix": tooltipPrefix,
				"valueSuffix": tooltipSuffix,
			},
		})
	}
	return res
}

func hsl(hue, light, sat int) string {
	return "hsl(" + strconv.Itoa(hue) + "," + strconv.Itoa(sat) + "%," + strconv.Itoa(light) + "%)"
}
