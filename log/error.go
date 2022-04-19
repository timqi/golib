package log

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/timqi/golib/json"
)

var (
	// error code 400~600 used for http specific
	ErrBadReq   = ErrStatus(http.StatusBadRequest)
	ErrInternal = ErrStatus(http.StatusInternalServerError)
	ErrTooMany  = ErrStatus(http.StatusTooManyRequests)

	// error code 1000~2000 used for infra errors
	Panic = Err{Code: 1000, Msg: "Panic"}

	// error code above 5000 used for business errors
	ErrGeneral          = Err{Code: 5000, Msg: "general error"}
	ErrParam            = Err{Code: 5001, Msg: "params error"}
	ErrInvalid          = Err{Code: 5002, Msg: "invalid error type"}
	ErrDB               = Err{Code: 5003, Msg: "db error"}
	ErrRedis            = Err{Code: 5004, Msg: "redis error"}
	ErrIPLocation       = Err{Code: 5005, Msg: "IP location error"}
	ErrNetwork          = Err{Code: 5006, Msg: "network error"}
	ErrJSON             = Err{Code: 5007, Msg: "json error"}
	ErrCipher           = Err{Code: 5008, Msg: "cipher error"}
	ErrMachineID        = Err{Code: 5009, Msg: "machineID error"}
	ErrPermissionDenied = Err{Code: 5010, Msg: "permission denied"}
	ErrCreateReq        = Err{Code: 5011, Msg: "create request error"}

	ErrSuperminerAdmin = Err{Code: 6000, Msg: "admin error"}

	ErrNoAuth = Err{Code: 10000, Msg: "no auth"}
)

type Err struct {
	Code int32
	Msg  string

	Pcs     []uintptr
	Fields  K
	Entries []entry
}

func (e Err) Error() string {
	var sb strings.Builder
	sb.WriteString(e.Msg)
	for k, v := range e.Fields {
		sb.WriteByte(' ')
		sb.WriteString(k)
		sb.WriteByte('=')
		vs, _ := json.MarshalToString(v)
		sb.WriteString(vs)
	}
	for _, entry := range e.Entries {
		sb.WriteByte(' ')
		sb.WriteString(entry.Key)
		sb.WriteByte('=')
		vs, _ := json.Marshal(entry.Val)
		sb.Write(vs)
	}
	return sb.String()
}

func (e *Err) WrapE(er error, entries ...entry) Err {
	return _wrap(e, er.Error(), nil, entries...)
}

func (e *Err) WrapEK(er error, fields K) Err {
	return _wrap(e, er.Error(), fields)
}

func (e *Err) Wrap(message string, entries ...entry) Err {
	return _wrap(e, message, nil, entries...)
}

func (e *Err) WrapK(message string, fields K) Err {
	return _wrap(e, message, fields)
}

func _wrap(e *Err, message string, fields K, entries ...entry) Err {
	fs := make(map[string]interface{})
	for k, v := range e.Fields {
		fs[k] = v
	}
	for k, v := range fields {
		fs[k] = v
	}
	pcs := make([]uintptr, 5)
	runtime.Callers(3, pcs)
	return Err{
		Code:    e.Code,
		Msg:     e.Msg + ": " + message,
		Pcs:     pcs,
		Fields:  fs,
		Entries: entries,
	}
}

func (e Err) PCS() []uintptr {
	return e.Pcs
}

func ErrStatus(code int) Err {
	return Err{Code: int32(code), Msg: http.StatusText(code)}
}
