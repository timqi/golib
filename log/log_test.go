package log_test

import (
	"bytes"
	"context"
	"errors"
	logStdLib "log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/timqi/golib/log"
)

func TestLog(t *testing.T) {
	c := context.Background()
	buf := &bytes.Buffer{}

	log.BaseConfig(log.DEBUG, false)
	log.SetLogger(log.DEBUG, logStdLib.New(buf, "D ", 0))

	log.Debug(c, "msg")
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "D log_test.go:24 msg\n", buf.String())
	buf.Reset()

	log.BaseConfig(log.INFO, false)
	log.SetLogger(log.INFO, logStdLib.New(buf, "I ", 0))
	log.SetLogger(log.ERROR, logStdLib.New(buf, "E ", 0))

	log.Debug(c, "debug msg")
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "", buf.String())
	buf.Reset()

	_ = log.Error(c, log.ErrParam)
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "E params error\n", buf.String())
	buf.Reset()

	_ = log.Error(c, log.ErrParam.Wrap("wrapped param error"))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "E params error: wrapped param error", strings.Split(buf.String(), "\n")[0])
	buf.Reset()

	_ = log.Error(c, log.ErrParam.WrapK("wrapped K", log.K{"field1": "k1"}))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, `E params error: wrapped K field1="k1"`, strings.Split(buf.String(), "\n")[0])
	buf.Reset()

	_ = log.Error(c, log.ErrParam.WrapE(log.ErrDB))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, `E params error: db error`, strings.Split(buf.String(), "\n")[0])
	buf.Reset()

	_ = log.Error(c, log.ErrParam.WrapE(log.ErrDB, log.Entry("entry", "e1")))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, `E params error: db error entry="e1"`, strings.Split(buf.String(), "\n")[0])
	buf.Reset()

	_ = log.Error(c, log.ErrParam.WrapEK(log.ErrDB, log.K{"field1 error": "k1"}))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, `E params error: db error field1 error="k1"`, strings.Split(buf.String(), "\n")[0])
	buf.Reset()

	log.InfoK(c, "debug", log.K{"f1": "field 1", "f2": "field 2"})
	time.Sleep(time.Millisecond * 100)
	{
		res := buf.String()
		assert.True(t, strings.Contains(res, "I log_test.go:68 debug"))
		assert.True(t, strings.Contains(res, `f1="field 1"`))
		assert.True(t, strings.Contains(res, `f2="field 2"`))
		buf.Reset()
	}

	log.Info(c, "debug2")
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, `I log_test.go:78 debug2`, strings.Split(buf.String(), "\n")[0])
	buf.Reset()

	log.Info(c,
		"debug with entries",
		log.Entry("entry1", "value1"),
		log.Entry("entry2", errors.New("error2")),
	)
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, `I log_test.go:83 debug with entries entry1="value1" entry2=error2`, strings.Split(buf.String(), "\n")[0])
	buf.Reset()
}
