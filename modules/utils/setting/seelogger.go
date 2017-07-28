//Utils is the tools lib
package setting

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lunny/log"
)

//Seelogger is a logger instance;
type Seelogger struct {
	*log.Logger
}

type seelogwriter struct {
	*sync.Mutex
	currentFileName string
	fd              *os.File
}

var (
	//Seelog can be used
	SeeLog *Seelogger = &Seelogger{Logger: log.New(os.Stderr, "", log.Ldefault())}
)

// -----------------------------------------

func (l *Seelogger) Errorf(format string, v ...interface{}) {
	l.Output("", log.Lerror, 2, fmt.Sprintf(format, v...))
	SendCount(1272, 1)
}

func (l *Seelogger) Error(v ...interface{}) {
	l.Output("", log.Lerror, 2, fmt.Sprintln(v...))
	SendCount(1272, 1)
}

// -----------------------------------------

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func (l *Seelogger) Fatal(v ...interface{}) {
	l.Output("", log.Lfatal, 2, fmt.Sprintln(v...))
	SendCount(1272, 1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func (l *Seelogger) Fatalf(format string, v ...interface{}) {
	l.Output("", log.Lfatal, 2, fmt.Sprintf(format, v...))
	SendCount(1272, 1)
}

// -----------------------------------------

// Panic is equivalent to Print() followed by a call to panic().
func (l *Seelogger) Panic(v ...interface{}) {
	l.Output("", log.Lpanic, 2, fmt.Sprintln(v...))
}

// Panicf is equivalent to Printf() followed by a call to panic().
func (l *Seelogger) Panicf(format string, v ...interface{}) {
	l.Output("", log.Lpanic, 2, fmt.Sprintf(format, v...))
}

// -----------------------------------------
// Write is the interface based split log
func (l *seelogwriter) Write(p []byte) (int, error) {

	defer l.Mutex.Unlock()

	timeFormat := time.Now().Format("20060102")

	l.Mutex.Lock()

	if l.fd == nil || !strings.EqualFold(timeFormat, l.currentFileName) {
		l.currentFileName = timeFormat
		f, err := os.OpenFile("logs/"+l.currentFileName+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
		if err != nil {
			panic("create log file failed:" + err.Error())
		}
		if l.fd != nil {
			l.fd.Close()
		}
		l.fd = f
	}

	return l.fd.Write(p)
}

func init() {
	os.MkdirAll("./logs", os.ModePerm)
	f := &seelogwriter{Mutex: new(sync.Mutex)}
	if IsProMode {
		SeeLog.SetOutput(f)
		SeeLog.SetOutputLevel(log.Linfo) //前期保留尽可能多的日志
	} else {
		w := io.MultiWriter(f, os.Stdout)
		SeeLog.SetOutput(w)
		SeeLog.SetOutputLevel(log.Ldebug)
	}
}
