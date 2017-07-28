package setting

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lunny/log"
)

//logger is a logger instance;
type logger struct {
	*log.Logger
}

type Logwriter struct {
	*sync.Mutex
	Prefix          string
	BufferSize      int
	currentFileName string
	fd              *os.File
	buffd           *bufio.Writer
}

var (
	//Logger can be used
	Logger        *logger = &logger{Logger: log.New(os.Stderr, "", log.Ldefault())}
	Logwriterfile *Logwriter
)

// -----------------------------------------

func (l *logger) Infof(format string, v ...interface{}) {
	l.Output("", log.Linfo, 2, fmt.Sprintf(format, v...))
	Logwriterfile.Flush()
}

func (l *logger) Info(v ...interface{}) {
	l.Output("", log.Linfo, 2, fmt.Sprintln(v...))
	Logwriterfile.Flush()
}

// -----------------------------------------

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Output("", log.Lerror, 2, fmt.Sprintf(format, v...))
	Logwriterfile.Flush()
}

func (l *logger) Error(v ...interface{}) {
	l.Output("", log.Lerror, 2, fmt.Sprintln(v...))
	Logwriterfile.Flush()
}

// -----------------------------------------

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func (l *logger) Fatal(v ...interface{}) {
	l.Output("", log.Lfatal, 2, fmt.Sprintln(v...))
	Logwriterfile.Flush()
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Output("", log.Lfatal, 2, fmt.Sprintf(format, v...))
	Logwriterfile.Flush()
}

// -----------------------------------------

// Panic is equivalent to Print() followed by a call to panic().
func (l *logger) Panic(v ...interface{}) {
	l.Output("", log.Lpanic, 2, fmt.Sprintln(v...))
	Logwriterfile.Flush()
}

// Panicf is equivalent to Printf() followed by a call to panic().
func (l *logger) Panicf(format string, v ...interface{}) {
	l.Output("", log.Lpanic, 2, fmt.Sprintf(format, v...))
	Logwriterfile.Flush()
}

// -----------------------------------------
// Write is the interface based split log
func (l *Logwriter) Write(p []byte) (int, error) {

	defer l.Mutex.Unlock()

	timeFormat := time.Now().Format("20060102")

	l.Mutex.Lock()

	if l.fd == nil || !strings.EqualFold(l.Prefix+timeFormat, l.currentFileName) {
		l.currentFileName = l.Prefix + timeFormat
		f, err := os.OpenFile("logs/"+l.currentFileName+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
		if err != nil {
			panic("create log file failed:" + err.Error())
		}
		if l.fd != nil {
			l.buffd.Flush()
			l.fd.Close()
		}
		l.fd = f
		l.buffd = bufio.NewWriterSize(f, l.BufferSize)
	}
	//return l.fd.Write(p)
	return l.buffd.Write(p)
}

func (l *Logwriter) Flush() error {
	return l.buffd.Flush()
}

func init() {
	os.MkdirAll("./logs", os.ModePerm)
	Logwriterfile = &Logwriter{Mutex: new(sync.Mutex), Prefix: "error", BufferSize: 1024}
	if IsProMode {
		Logger.SetOutput(Logwriterfile)
		Logger.SetOutputLevel(log.Linfo) //前期保留尽可能多的日志
	} else {
		w := io.MultiWriter(Logwriterfile, os.Stdout)
		Logger.SetOutput(w)
		Logger.SetOutputLevel(log.Ldebug)
	}
}
