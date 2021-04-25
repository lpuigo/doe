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
	//log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(f)

	return f
}

type Record struct {
	t        time.Time
	source   string
	user     string
	Request  string
	Response int
	Info     string
	Error    string
}

func Entry(source string) *Record {
	return &Record{source: source}
}

func TimedEntry(source string) *Record {
	e := Entry(source)
	e.t = time.Now()
	return e
}

func (e *Record) getMsg() string {
	msg := ""
	if e.source != "" {
		msg += e.source
	}
	if e.user != "" {
		msg += ` user="` + e.user + `"`
	}
	if e.Request != "" {
		msg += ` request="` + e.Request + `"`
	}
	if e.Response != 0 {
		msg += ` response=` + strconv.Itoa(e.Response)
	}
	if e.Info != "" {
		msg += ` info="` + e.Info + `"`
	}
	if e.Error != "" {
		msg += ` error="` + e.Error + `"`
	}
	if !e.t.IsZero() {
		msg += " service_time=" + strconv.FormatFloat(float64(time.Since(e.t).Nanoseconds())/1000000, 'f', 3, 64) + "ms"
	}
	return msg
}

func (e *Record) Log() {
	log.Println(e.getMsg())
}

func (e *Record) LogInfo(info string) {
	e.Info = info
	e.Log()
}

func (e *Record) LogError(err string) {
	e.Error = err
	e.Log()
}

func (e *Record) LogErr(err error) {
	if err == nil {
		return
	}
	e.Error = err.Error()
	e.Log()
}

func (e *Record) Fatal(err error) {
	e.Error = err.Error()
	log.Fatal(e.getMsg())
}

func (e *Record) AddTime() *Record {
	e.t = time.Now()
	return e
}

func (e *Record) AddRequest(req string) *Record {
	e.Request = req
	return e
}

func (e *Record) AddUser(user string) *Record {
	e.user = user
	return e
}

func (e *Record) AddInfoResponse(inf string, code int) *Record {
	e.Info = inf
	e.Response = code
	return e
}

func (e *Record) AddResponse(code int) *Record {
	e.Response = code
	return e
}
