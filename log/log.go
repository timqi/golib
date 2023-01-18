package log

import (
	"context"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/timqi/golib/ctx"
)

type K map[string]interface{}

type entry struct {
	Key string
	Val interface{}
}

func Entry(key string, val interface{}) entry {
	return entry{Key: key, Val: val}
}

func KV(k string, v interface{}) entry {
	return entry{Key: k, Val: v}
}

const (
	DEBUG = 1 << iota
	INFO
	WARN
	ERROR
)

func BaseConfig(level int, isDebug bool) {
	configQ <- ChanCfg{
		Level: level,
		Debug: isDebug,
	}
}

func SetLogger(level int, logger *log.Logger) {
	setLoggerQ <- ChanSetLogger{
		Level:  level,
		Logger: logger,
	}
}

func msgWithFile(message string) string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		sb := &strings.Builder{}
		sb.WriteString(filepath.Base(file))
		sb.WriteByte(':')
		sb.WriteString(strconv.Itoa(line))
		sb.WriteByte(' ')
		sb.WriteString(message)
		return sb.String()
	}
	return message
}

func DebugK(c context.Context, message string, fields K) {
	logQ <- ChanLog{
		Level:   DEBUG,
		LogID:   ctx.GetLogID(c),
		Message: msgWithFile(message),
		Fields:  fields,
	}
}

func InfoK(c context.Context, message string, fields K) {
	logQ <- ChanLog{
		Level:   INFO,
		LogID:   ctx.GetLogID(c),
		Message: msgWithFile(message),
		Fields:  fields,
	}
}

func WarnK(c context.Context, message string, fields K) {
	logQ <- ChanLog{
		Level:   WARN,
		LogID:   ctx.GetLogID(c),
		Message: msgWithFile(message),
		Fields:  fields,
	}
}

func Debug(c context.Context, message string, entries ...entry) {
	logQ <- ChanLog{
		Level:   DEBUG,
		LogID:   ctx.GetLogID(c),
		Message: msgWithFile(message),
		Entries: entries,
	}
}

func Info(c context.Context, message string, entries ...entry) {
	logQ <- ChanLog{
		Level:   INFO,
		LogID:   ctx.GetLogID(c),
		Message: msgWithFile(message),
		Entries: entries,
	}
}

func Warn(c context.Context, message string, entries ...entry) {
	logQ <- ChanLog{
		Level:   WARN,
		LogID:   ctx.GetLogID(c),
		Message: msgWithFile(message),
		Entries: entries,
	}
}

func Error(c context.Context, err Err) Err {
	logQ <- ChanLog{
		Level: ERROR,
		LogID: ctx.GetLogID(c),
		Err:   &err,
	}
	return err
}
