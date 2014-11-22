package freeserver

import (
	"fmt"

	"github.com/koding/logging"
	"github.com/rcrowley/go-tigertonic"
)

type TigerTonicLogger struct{ log logging.Logger }

func NewTigerTonicLogger(log logging.Logger) tigertonic.Logger {
	return &TigerTonicLogger{log: log}
}

func (t *TigerTonicLogger) Print(v ...interface{}) {
	t.Output(2, fmt.Sprint(v...))
}

func (t *TigerTonicLogger) Println(v ...interface{}) {
	t.Output(2, fmt.Sprintln(v...))
}

func (t *TigerTonicLogger) Printf(format string, v ...interface{}) {
	t.Output(2, fmt.Sprintf(format, v...))
}

func (t *TigerTonicLogger) Output(calldepth int, s string) error {
	t.log.Debug(s)
	return nil
}
