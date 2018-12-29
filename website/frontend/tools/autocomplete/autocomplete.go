package autocomplete

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"strconv"
)

type Result struct {
	*js.Object
	Value string `js:"value"`
}

func NewResult(val string) *Result {
	r := &Result{Object: tools.O()}
	r.Value = val
	return r
}

func GenResults(prefix, val string, number int) []*Result {
	res := []*Result{}
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		v = 0
	}
	valLength := len(val)
	zeroPrefix := "000000000000000000000000"
	for i := 1; i <= number; i++ {
		newval := strconv.Itoa(int(v) + i)
		newval = zeroPrefix[:valLength-len(newval)] + newval
		res = append(res, NewResult(prefix+newval))
	}
	return res
}
