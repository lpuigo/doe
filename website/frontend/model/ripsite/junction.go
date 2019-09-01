package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/ripsite/ripconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

type Operation struct {
	*js.Object

	Type        string `js:"Type"`
	TronconName string `js:"TronconName"`
	NbFiber     int    `js:"NbFiber"`
	NbSplice    int    `js:"NbSplice"`
	State       *State `js:"State"`
}

func NewOperation() *Operation {
	o := &Operation{Object: tools.O()}
	o.Type = ""
	o.TronconName = ""
	o.NbFiber = 0
	o.NbSplice = 0
	o.State = NewState()

	return o
}

type Junction struct {
	*js.Object

	NodeName   string       `js:"NodeName"`
	Operations []*Operation `js:"Operations"`
	State      *State       `js:"State"`
}

func NewJunction() *Junction {
	j := &Junction{Object: tools.O()}
	j.NodeName = ""
	j.Operations = nil
	j.State = NewState()

	return j
}

func (j *Junction) Clone() *Junction {
	return &Junction{Object: json.Parse(json.Stringify(j))}
}

func (j *Junction) GetNbFiber() int {
	nb := 0
	for _, op := range j.Operations {
		nb += op.NbFiber
	}
	return nb
}

func (j *Junction) SearchString(filter string) string {
	searchItem := func(prefix, typ, value string) string {
		if value == "" {
			return ""
		}
		if filter != ripconst.FilterValueAll && filter != typ {
			return ""
		}
		return prefix + typ + value
	}

	res := searchItem("", ripconst.FilterValueComment, j.State.Comment)
	res += searchItem(",", ripconst.FilterValuePtRef, j.NodeName)
	for _, ope := range j.Operations {
		res += searchItem(",", ripconst.FilterValueTrRef, ope.TronconName)
		res += searchItem(",", ripconst.FilterValueOpe, ope.Type)
	}
	return res
}
