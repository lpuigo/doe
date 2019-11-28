package measurementreport

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"io"
	"regexp"
	"strconv"
	"strings"
)

const (
	dbThresholdWarn1       float64 = 0.200001
	dbThresholdWarn2       float64 = 0.3
	dbThresholdKo          float64 = 0.4
	dbConnectorThresholdKo float64 = 0.8
)

func ParseZipMeasurementFiles(r io.ReaderAt, size int64) (map[string]*MeasurementReport, error) {
	zreader, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}

	res := make(map[string]*MeasurementReport)
	for _, file := range zreader.File {
		if strings.HasSuffix(file.Name, ".txt") {
			mr, err := parseZipTxtFile(file)
			if err != nil {
				return nil, fmt.Errorf("Failed to read '%s' from zip: %s", file.Name, err)
			}
			res[mr.PtName] = mr
		}
	}
	return res, nil
}

func parseZipTxtFile(file *zip.File) (*MeasurementReport, error) {
	fileread, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("Failed to open zip %s for reading: %s", file.Name, err)
	}
	defer fileread.Close()
	return parserTxtFile(fileread)
}

func parserTxtFile(r io.Reader) (*MeasurementReport, error) {
	scan := bufio.NewScanner(r)
	mr := NewMeasurementReport()
	var state txtParserState
	state = getDate
	for scan.Scan() {
		state = state(mr, scan.Text())
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return mr, nil
}

func splitLine(line string) []string {
	cols := strings.Split(line, "\t")
	for i, col := range cols {
		cols[i] = strings.Trim(col, " ")
	}
	return cols
}

type txtParserState func(*MeasurementReport, string) txtParserState

func getDate(mr *MeasurementReport, line string) txtParserState {
	if !strings.HasPrefix(line, "Date") {
		return getDate
	}
	col := splitLine(line)
	mr.Date = date.ChangeDDMMYYYYtoYYYYMMDD(col[1])
	return getTime
}

func getTime(mr *MeasurementReport, line string) txtParserState {
	if !strings.HasPrefix(line, "Time") {
		return getTime
	}
	col := splitLine(line)
	mr.Time = col[1]
	return getPTName

}

func getPTName(mr *MeasurementReport, line string) txtParserState {
	if !strings.HasPrefix(line, "File Name") {
		return getPTName
	}
	col := splitLine(line)
	chunks := strings.Split(col[1], "/")
	reTR := regexp.MustCompile(`TR *[0-9]+`)
	rePT := regexp.MustCompile(`PT *[0-9]+`)
	trFound := false
	for _, dir := range chunks {
		if !trFound {
			if trIndex := reTR.FindStringIndex(dir); trIndex != nil {
				mr.Troncon = "TR " + strings.Replace(dir[trIndex[0]+2:trIndex[1]], " ", "", -1)
				trFound = true
				continue
			}
		}
		if trFound {
			if ptIndex := rePT.FindStringIndex(dir); ptIndex != nil {
				mr.PtName = "PT " + strings.Replace(dir[ptIndex[0]+2:ptIndex[1]], " ", "", -1)
				break
			}
		}
	}
	return getAlarms
}

func getAlarms(mr *MeasurementReport, line string) txtParserState {
	if !strings.HasPrefix(line, "Alarms") {
		return getAlarms
	}
	return getMeasurement
}

func getMeasurement(mr *MeasurementReport, line string) txtParserState {
	col := splitLine(line)
	if !(len(col) > 10) {
		return getMeasurement // TODO How to manage End of parsing ???
	}
	msg := []string{}
	maxSplice, err := strconv.ParseFloat(col[7], 64)
	if err != nil {
		if col[7] == "-" {
			maxSplice = 0
		}
	}
	switch {
	case maxSplice >= dbThresholdKo:
		mr.FiberKO++
		msg = append(msg, fmt.Sprintf("KO Max Splice %sdb à %sm", col[7], col[10]))
	case maxSplice >= dbThresholdWarn2:
		mr.FiberWarning2++
		msg = append(msg, fmt.Sprintf("Warn2 Max Splice %sdb à %sm", col[7], col[10]))
	case maxSplice >= dbThresholdWarn1:
		mr.FiberWarning1++
		msg = append(msg, fmt.Sprintf("Warn1 Max Splice %sdb à %sm", col[7], col[10]))
	default:
		mr.FiberOK++
	}

	connector, err := strconv.ParseFloat(col[9], 64)
	if err == nil {
		if connector > dbConnectorThresholdKo {
			mr.ConnectorKO++
			msg = append(msg, fmt.Sprintf("Max Connector %sdb", col[9]))
		}
	} else {
		fmt.Printf("Err Connector: %s", col[9])
	}
	if len(msg) > 0 {
		mr.Results = append(mr.Results, fmt.Sprintf("Fib. #%s: %s", col[1], strings.Join(msg, ", ")))
	}
	return getMeasurement
}
