package nominatim

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"honnef.co/go/js/xhr"
	"strconv"
	"strings"
)

type Nominatim struct {
	*js.Object

	Lat      float64    `js:"Lat"`
	Long     float64    `js:"Long"`
	Found    bool       `js:"Found"`
	Err      string     `js:"Err"`
	VM       *hvue.VM   `js:"VM"`
	Response *js.Object `js:"Response"`
}

func NewNominatim(vm *hvue.VM) *Nominatim {
	n := &Nominatim{Object: tools.O()}
	n.Lat = 0
	n.Long = 0
	n.Found = false
	n.Err = ""
	n.VM = vm
	n.Response = nil

	return n
}

const (
	nominatimUrl        string = "https://nominatim.openstreetmap.org/search?q="
	nominatimFormatJson string = "&format=json"
)

func (n *Nominatim) SearchAdress(addr string, callback func()) {
	adr := strings.Replace(strings.Trim(addr, " "), " ", "+", -1)
	uri := nominatimUrl + adr + nominatimFormatJson

	go n.callSearchAdress(uri, callback)
}

func (n *Nominatim) callSearchAdress(uri string, callback func()) {
	req := xhr.NewRequest("GET", uri)
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(n.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(n.VM, req)
		return
	}
	n.Response = req.Response
	defer callback()
	if n.Response.Length() == 0 {
		return
	}

	n.Found = true
	faddr := n.Response.Index(0)

	n.Lat = n.parseFloat(faddr.Get("lat").String(), "lat")
	if n.Err != "" {
		return
	}
	n.Long = n.parseFloat(faddr.Get("lon").String(), "lon")
	if n.Err != "" {
		return
	}
}

func (n *Nominatim) parseFloat(s, attr string) float64 {
	res, err := strconv.ParseFloat(s, 64)
	if err != nil {
		n.Err = "could not read " + attr + " from '" + s + "'"
		return 0
	}
	return res
}
