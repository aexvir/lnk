package logging

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	namespace string
	log       *zap.Logger
}

func NewLogger(namespace string) *Logger {
	log, _ := zap.NewProduction(zap.WithCaller(false))
	return &Logger{
		namespace: namespace,
		log:       log,
	}
}

func (l *Logger) Write(event, msg string, args ...any) {
	l.log.Info(
		fmt.Sprintf(msg, args...),
		zap.String("event", fmt.Sprintf("%s.%s", l.namespace, event)),
	)
}

func (l *Logger) Error(msg string, args ...any) {
	l.Write("error", msg, args...)
}
