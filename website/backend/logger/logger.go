package logger

import (
	"log"
	"os"
	"strconv"
	"time"
)

// StartLog setups logger to write log in given file
//
// Do defer file.Close() just after StartLog call to ensure proper log file closing
func StartLog(logfile string) *os.File {
	//create your file with desired read/write permissions
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//set output of logs to f
	log.SetOutput(f)
	return f
}

func LogRequest(req string) string {
	return `request="` + req + `"`
}

func LogResponse(resp int) string {
	return " response=" + strconv.Itoa(resp)
}

func LogInfo(info string) string {
	return ` info="` + info + `"`
}

func LogResponseInfo(info string, resp int) string {
	return LogResponse(resp) + LogInfo(info)
}

func LogService(t time.Time, msg *string) {
	log.Printf("%s service_time=%.3fms\n", *msg, float64(time.Since(t).Nanoseconds())/1000000)
}
