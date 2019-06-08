package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"strings"
)

type MeasurementReport struct {
	*js.Object

	Troncon       string   `js:"Troncon"`
	PtName        string   `js:"PtName"`
	Date          string   `js:"Date"`
	Time          string   `js:"Time"`
	FiberWarning1 int      `js:"FiberWarning1"`
	FiberWarning2 int      `js:"FiberWarning2"`
	FiberKO       int      `js:"FiberKO"`
	FiberOK       int      `js:"FiberOK"`
	ConnectorKO   int      `js:"ConnectorKO"`
	Results       []string `js:"Results"`
}

func (mr *MeasurementReport) Comments() string {
	return strings.Join(mr.Results, "\n")
}
