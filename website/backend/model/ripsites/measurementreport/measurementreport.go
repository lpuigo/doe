package measurementreport

import (
	"fmt"
	"strings"
)

type MeasurementReport struct {
	Troncon      string
	PtName       string
	Date         string
	Time         string
	FiberWarning int
	FiberKO      int
	FiberOK      int
	ConnectorKO  int
	Results      []string
}

func (mr *MeasurementReport) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s %s: %s (%s) %d fibers\n", mr.Date, mr.Time, mr.Troncon, mr.PtName, mr.FiberOK+mr.FiberWarning+mr.FiberKO))
	sb.WriteString(fmt.Sprintf("\t%2d OK %2d Warn %2d KO / %2d Connect. KO\n", mr.FiberOK, mr.FiberWarning, mr.FiberKO, mr.ConnectorKO))
	if len(mr.Results) > 0 {
		sb.WriteString(fmt.Sprintf("%s\n", mr.Comments()))
	}
	return sb.String()
}

func (mr *MeasurementReport) Comments() string {
	return strings.Join(mr.Results, "\n")
}
