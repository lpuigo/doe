package measurementreport

import (
	"fmt"
	"strings"
)

type MeasurementReport struct {
	Troncon       string
	PtName        string
	Date          string
	Time          string
	FiberWarning1 int
	FiberWarning2 int
	FiberKO       int
	FiberOK       int
	ConnectorKO   int
	Results       []string
}

func NewMeasurementReport() *MeasurementReport {
	return &MeasurementReport{Results: []string{}}
}

func (mr *MeasurementReport) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s %s: %s (%s) %2d fibers : ", mr.Date, mr.Time, mr.PtName, mr.Troncon, mr.FiberOK+mr.FiberWarning1+mr.FiberWarning2+mr.FiberKO))
	sb.WriteString(fmt.Sprintf("%2d OK %2d Warn1 %2d Warn2 %2d KO %2d Connect. KO\n", mr.FiberOK, mr.FiberWarning1, mr.FiberWarning2, mr.FiberKO, mr.ConnectorKO))
	if len(mr.Results) > 0 {
		sb.WriteString(fmt.Sprintf("%s\n", mr.Comments()))
	}
	return sb.String()
}

func (mr *MeasurementReport) Comments() string {
	return strings.Join(mr.Results, "\n")
}
