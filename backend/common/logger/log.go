
package logger

import (
	"fmt"
	"log"
	"time"
)

type severity int

const (
	debug severity = iota
	info
	warn
	err
	fatal
	panic
	emergency
)

type Microservice string

const (
	TFC Microservice = "TFC"
)

type LogMessage struct {
	Microservice Microservice `json:"microservice" bson:"microservice"`
	Severity     severity     `json:"severity" bson:"severity"`
	Message      string       `json:"message" bson:"message"`
	Timestamp    time.Time    `json:"time" bson:"time"`
}

type logger interface {
	Log(severity severity, message string)
}

var l logger

func write(s severity, msg string, args ...any) {
	if l == nil {
		log.Printf("[fallback][%v] %s", s, fmt.Sprintf(msg, args...))
		return
	}

	l.Log(s, fmt.Sprintf(msg, args...))

	if s >= fatal {
		log.Panic("something fatal or more extreme happened")
	}
}

func Debug(msg string, args ...any) {
	go write(debug, msg, args...)
}

func Info(msg string, args ...any) {
	go write(info, msg, args...)
}

func Warn(msg string, args ...any) {
	go write(warn, msg, args...)
}

func Error(msg string, args ...any) {
	go write(err, msg, args...)
}

func Fatal(msg string, args ...any) {
	go write(fatal, msg, args...)
}

func Panic(msg string, args ...any) {
	go write(panic, msg, args...)
}

func Emergency(msg string, args ...any) {
	go write(emergency, msg, args...)
}
