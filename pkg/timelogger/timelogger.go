package timelogger

import (
	"encoding/csv"
	"os"
	"runtime"
	"strconv"
	"time"
)

type TimeLogger interface {
	Start()
	End()
}

type tlogger struct {
	file   *os.File
	writer *csv.Writer
	start  time.Time
	end    time.Time
	name   string
}

func (t *tlogger) Start() {
	t.start = time.Now()
}

func (t *tlogger) End() {
	t.end = time.Now()
	record := []string{t.name, strconv.Itoa(int(t.end.Sub(t.start).Microseconds()))}
	err := t.writer.Write(record)
	if err != nil {
		panic(err)
	}
	t.writer.Flush()
}

func NewTimeLogger(methodname, filename string) (t *tlogger) {
	logFile, _ := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	t = &tlogger{
		file:   logFile,
		writer: csv.NewWriter(logFile),
		name:   methodname,
	}
	runtime.SetFinalizer(t, func(t *tlogger) {
		t.file.Close()
	})
	return
}
