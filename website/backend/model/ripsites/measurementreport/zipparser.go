package measurementreport

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"io"
	"strconv"
	"strings"
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
	mr := &MeasurementReport{}
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
	for i, dir := range chunks {
		if i > 3 && strings.HasPrefix(dir, "TR") {
			mr.Troncon = dir
		}
		if i > 4 && strings.HasPrefix(dir, "PT") {
			mr.PtName = dir
			break
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
	if maxSplice > 0.3 {
		mr.FiberKO++
		msg = append(msg, fmt.Sprintf("KO Max Splice %sdb à %sm", col[7], col[10]))
	} else if maxSplice > 0.2 {
		mr.FiberWarning++
		msg = append(msg, fmt.Sprintf("Warn Max Splice %sdb à %sm", col[7], col[10]))
	} else {
		mr.FiberOK++
	}

	connector, err := strconv.ParseFloat(col[9], 64)
	if err == nil {
		if connector > 0.4 {
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
