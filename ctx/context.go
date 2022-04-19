package ctx

import (
	"context"

	"github.com/rs/xid"
)

var (
	debug = false
)

func Config(isDebug bool) {
	debug = isDebug
}

const (
	KEY_LOGID = "_lid"
)

func GetLogID(c context.Context) string {
	if c == nil {
		return ""
	}
	logid, ok := c.Value(KEY_LOGID).(string)
	if !ok {
		return ""
	}
	return logid
}

func NewWithID() context.Context {
	c := context.Background()
	if debug {
		return c
	}
	return context.WithValue(c, KEY_LOGID, xid.New().String())
}

func GetUserID(c context.Context) int64 {
	if c == nil {
		return 0
	}
	userID, ok := c.Value("user_id").(int64)
	if !ok {
		return 0
	}
	return userID
}

func GetUserName(c context.Context) string {
	if c == nil {
		return ""
	}
	userName, ok := c.Value("user_name").(string)
	if !ok {
		return ""
	}
	return userName
}
