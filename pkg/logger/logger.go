package logger

import (
	"log"
)

const (
	SILENT  = 0
	INFO    = 1 << 0
	WARNING = 1 << 1
	ERROR   = 1 << 2
	FATAL   = 1 << 3
	DEBUG   = 1 << 4
	TRACE   = 1 << 5
)

type Logger struct {
	level int
}

func NewLogger(levelName string) Logger {
	level := INFO | WARNING | ERROR | FATAL
	switch levelName {
	case "silent":
		level = SILENT
	case "info":
		level = INFO | WARNING | ERROR | FATAL
	case "debug":
		level = INFO | WARNING | ERROR | FATAL | DEBUG
	case "trace":
		level = INFO | WARNING | ERROR | FATAL | TRACE
	case "all":
		level = INFO | WARNING | ERROR | FATAL | DEBUG | TRACE
	}
	return Logger{
		level: level,
	}
}

func (l *Logger) Infof(format string, v ...any) {
	if (l.level & INFO) == INFO {
		log.Printf(" [INFO]  "+format+"\n", v...)
	}
}

func (l *Logger) Debugf(format string, v ...any) {
	if (l.level & DEBUG) == DEBUG {
		log.Printf(" [DEBUG] "+format+"\n", v)
	}
}

func (l *Logger) Warnf(format string, v ...any) {
	if (l.level & WARNING) == WARNING {
		log.Printf(" [WARN]  "+format+"\n", v)
	}
}

func (l *Logger) Errorf(format string, v ...any) {
	if (l.level & ERROR) == ERROR {
		log.Printf(" [ERROR] "+format+"\n", v)
	}
}

func (l *Logger) Fatalf(format string, v ...any) {
	if (l.level & FATAL) == FATAL {
		log.Printf(" [FATAL] "+format+"\n", v)
	}
}

func (l *Logger) Tracef(format string, v ...any) {
	if (l.level & TRACE) == TRACE {
		log.Printf(" [TRACE] "+format+"\n", v)
	}
}
