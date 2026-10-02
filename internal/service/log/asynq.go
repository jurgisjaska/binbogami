package log

import (
	"fmt"
	"log/slog"
	"os"
)

// Asynq adapts slog.Logger to the asynq.Logger interface.
type Asynq struct {
	l *slog.Logger
}

func CreateAsynq(l *slog.Logger) *Asynq {
	return &Asynq{l: l}
}

func (a *Asynq) Debug(args ...interface{}) { a.l.Debug(fmt.Sprint(args...)) }
func (a *Asynq) Info(args ...interface{})  { a.l.Info(fmt.Sprint(args...)) }
func (a *Asynq) Warn(args ...interface{})  { a.l.Warn(fmt.Sprint(args...)) }
func (a *Asynq) Error(args ...interface{}) { a.l.Error(fmt.Sprint(args...)) }
func (a *Asynq) Fatal(args ...interface{}) {
	a.l.Error(fmt.Sprint(args...))
	os.Exit(1)
}
