// log package for universal mono projects
// base methods provided should info, debug, error
//
//  1. you should push context in the first param when you have,
//     context will always have trace id like logid within
//  2. message passed should as short as you can, a log.K
//     to represent every fields you want print
//  3. use nil for context if there is no context
package log

import (
	"context"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/timqi/golib/ctx"
	"github.com/timqi/golib/json"
)

func init() {
	BaseConfig("debug", true)
}

var dLogger, iLogger, eLogger *log.Logger

func GetErrorLogger() *log.Logger {
	return eLogger
}

type K map[string]interface{}

func BaseConfig(level string, isDebugEnv bool) {
	if isDebugEnv {
		flag := log.Ltime | log.Lmicroseconds
		dLogger = log.New(os.Stdout, "\x1B[36mD ", flag|log.Lshortfile)
		iLogger = log.New(os.Stdout, "\x1B[32mI ", flag|log.Lshortfile)
		eLogger = log.New(os.Stderr, "\x1B[31mE ", flag)
	} else {
		flag := log.Ldate | log.Ltime | log.Lmicroseconds
		dLogger = log.New(os.Stdout, "D ", flag|log.Lshortfile)
		iLogger = log.New(os.Stdout, "I ", flag|log.Lshortfile)
		eLogger = log.New(os.Stderr, "E ", flag)
	}
	switch strings.ToLower(level) {
	case "info":
		dLogger = nil
	case "error":
		dLogger = nil
		iLogger = nil
	}
}

func SetLogger(level string, logger *log.Logger) {
    switch strings.ToLower(level) {
    case "debug":
        dLogger = logger
    case "info":
        iLogger = logger
    case "error":
        eLogger = logger
    }
}

func ExtConfig(d, i, e *log.Logger) {
	dLogger = d
	iLogger = i
	eLogger = e
}

func _buildLogString(
	c context.Context,
	message string,
	fields K,
	entries ...entry,
) string {
	var sb strings.Builder
	if logID := ctx.GetLogID(c); logID != "" {
		sb.WriteString(logID)
		sb.WriteByte(' ')
	}
	sb.WriteString(message)

	var ff []string
	for k := range fields {
		ff = append(ff, k)
	}
	sort.Strings(ff)

	for _, k := range ff {
		_buildLogKV(&sb, k, fields[k])
	}
	for _, entry := range entries {
		_buildLogKV(&sb, entry.Key, entry.Val)
	}
	return sb.String()
}

func _buildLogKV(sb *strings.Builder, k string, v interface{}) {
	sb.WriteByte(' ')
	sb.WriteString(k)
	sb.WriteByte('=')
	switch t := v.(type) {
	case error:
		sb.WriteString(t.Error())
	default:
		vs, _ := json.Marshal(v)
		sb.Write(vs)
	}
}

func DebugK(c context.Context, message string, fields K) {
	if dLogger != nil {
		dLogger.Output(2, _buildLogString(c, message, fields))
	}
}

func InfoK(c context.Context, message string, fields K) {
	if iLogger != nil {
		iLogger.Output(2, _buildLogString(c, message, fields))
	}
}

func Debug(c context.Context, message string, entries ...entry) {
	if dLogger != nil {
		dLogger.Output(2, _buildLogString(c, message, nil, entries...))
	}
}

func Info(c context.Context, message string, entries ...entry) {
	if iLogger != nil {
		iLogger.Output(2, _buildLogString(c, message, nil, entries...))
	}
}

type pcs interface {
	PCS() []uintptr
}

func Error(c context.Context, err Err) Err {
	logID := ctx.GetLogID(c)
	if logID != "" {
		logID = logID + " "
	}
	eLogger.Println(logID + err.Error())

	for _, pc := range err.PCS() {
		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		file, line := f.FileLine(pc)
		eLogger.Printf("%s  %s:%d\n", logID, file, line)
	}
	return err
}

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

func UnixSec(t int64) string {
	return time.Unix(t, 0).UTC().Format("2006-01-02T15:04:05Z")
}
