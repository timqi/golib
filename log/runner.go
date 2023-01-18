package log

import (
	"log"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/timqi/golib/json"
)

type ChanLog struct {
	Level   int
	LogID   string
	Message string
	Fields  K
	Entries []entry
	Err     *Err
}

type ChanSetLogger struct {
	Level  int
	Logger *log.Logger
}

type ChanCfg struct {
	Level int
	Debug bool
}

var (
	logQ       chan ChanLog       = make(chan ChanLog, 100)
	setLoggerQ chan ChanSetLogger = make(chan ChanSetLogger)
	configQ    chan ChanCfg       = make(chan ChanCfg)

	loggers map[int]*log.Logger
)

func doBasicConfig(level int, debug bool) {
	type cfgArg struct {
		Prefix string
		Flag   int
	}
	m := make(map[int]cfgArg)
	flag := log.Ldate | log.Ltime | log.Lmicroseconds
	m[DEBUG] = cfgArg{"D ", flag}
	m[INFO] = cfgArg{"I ", flag}
	m[WARN] = cfgArg{"W ", flag}
	m[ERROR] = cfgArg{"E ", flag}

	colorMap := map[int]string{
		DEBUG: "", INFO: "", WARN: "", ERROR: "",
	}
	if debug {
		colorMap = map[int]string{
			DEBUG: "\x1B[36m",
			INFO:  "\x1B[32m",
			WARN:  "\x1B[33m",
			ERROR: "\x1B[31m",
		}
	}

	loggers = make(map[int]*log.Logger)
	loggers[ERROR] = log.New(os.Stdout,
		colorMap[ERROR]+m[ERROR].Prefix,
		m[ERROR].Flag,
	)

	if level > WARN {
		return
	}
	loggers[WARN] = log.New(os.Stdout,
		colorMap[WARN]+m[WARN].Prefix,
		m[WARN].Flag,
	)

	if level > INFO {
		return
	}
	loggers[INFO] = log.New(os.Stdout,
		colorMap[INFO]+m[INFO].Prefix,
		m[INFO].Flag,
	)

	if level > DEBUG {
		return
	}
	loggers[DEBUG] = log.New(os.Stdout,
		colorMap[DEBUG]+m[DEBUG].Prefix,
		m[DEBUG].Flag,
	)
}

func buildLogString(
	logID string,
	message string,
	fields K,
	entries []entry,
) string {
	var sb strings.Builder
	if logID != "" {
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
		buildLogKV(&sb, k, fields[k])
	}
	for _, entry := range entries {
		buildLogKV(&sb, entry.Key, entry.Val)
	}

	return sb.String()
}

func buildLogKV(sb *strings.Builder, k string, v interface{}) {
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

func run() {
	doBasicConfig(DEBUG, true)
	for {
		select {
		case cfg := <-configQ:
			doBasicConfig(cfg.Level, cfg.Debug)
		case l := <-setLoggerQ:
			loggers[l.Level] = l.Logger
		case li := <-logQ:
			logger, exist := loggers[li.Level]
			if !exist {
				continue
			}
			if li.Level == ERROR {
				if li.LogID != "" {
					li.LogID = li.LogID + " "
				}
				logger.Println(li.LogID + li.Err.Error())
				for _, pc := range li.Err.PCS() {
					f := runtime.FuncForPC(pc)
					if f == nil {
						break
					}
					file, line := f.FileLine(pc)
					logger.Printf("%s  %s:%d\n", li.LogID, file, line)
				}
			} else {
				logger.Output(0, buildLogString(
					li.LogID,
					li.Message,
					li.Fields,
					li.Entries,
				))
			}
		}
	}
}

func init() {
	go run()
}
