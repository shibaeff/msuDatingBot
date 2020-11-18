package timelogger

import (
	"os"
	"runtime"
)

type TimeLogger interface {
	Start(name string)
	End()
}

type tlogger struct {
	file *os.File
}

func (t *tlogger) Start(name string) {
	panic("implement me")
}

func (t *tlogger) End() {
	panic("implement me")
}

func NewTimeLogger(filename string) (t *tlogger) {
	logFile, _ := os.OpenFile("time.csv", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	t = &tlogger{
		file: logFile,
	}
	runtime.SetFinalizer(t, func(t *tlogger) {
		t.file.Close()
	})
	return
}
