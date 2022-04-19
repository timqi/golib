package log_test

import (
	"bytes"
	"context"
	"errors"
	logStdLib "log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/timqi/golib/log"
)

func TestLogDebug(t *testing.T) {
	c := context.Background()
	buf := &bytes.Buffer{}
	log.ExtConfig(logStdLib.New(buf, "D ", 0), nil, logStdLib.New(buf, "E ", 0))
	_ = log.Error(c, log.ErrParam)
	assert.Equal(t, "E params error\n", buf.String())
	buf.Reset()
	_ = log.Error(c, log.ErrParam.Wrap("wrapped param error"))
	assert.Equal(t, "E params error: wrapped param error", strings.Split(buf.String(), "\n")[0])
	buf.Reset()
	_ = log.Error(c, log.ErrParam.WrapK("wrapped K", log.K{"field1": "k1"}))
	assert.Equal(t, `E params error: wrapped K field1="k1"`, strings.Split(buf.String(), "\n")[0])
	buf.Reset()
	_ = log.Error(c, log.ErrParam.WrapE(log.ErrDB))
	assert.Equal(t, `E params error: db error`, strings.Split(buf.String(), "\n")[0])
	buf.Reset()
	_ = log.Error(c, log.ErrParam.WrapE(log.ErrDB, log.Entry("entry", "e1")))
	assert.Equal(t, `E params error: db error entry="e1"`, strings.Split(buf.String(), "\n")[0])
	buf.Reset()
	_ = log.Error(c, log.ErrParam.WrapEK(log.ErrDB, log.K{"field1 error": "k1"}))
	assert.Equal(t, `E params error: db error field1 error="k1"`, strings.Split(buf.String(), "\n")[0])
	buf.Reset()

	log.DebugK(c, "debug", log.K{"f1": "field 1", "f2": "field 2"})
	{
		res := buf.String()
		assert.True(t, strings.Contains(res, "D debug"))
		assert.True(t, strings.Contains(res, `f1="field 1"`))
		assert.True(t, strings.Contains(res, `f2="field 2"`))
		buf.Reset()
	}
	log.Debug(c, "debug2")
	assert.Equal(t, `D debug2`, strings.Split(buf.String(), "\n")[0])
	buf.Reset()
	log.Debug(c,
		"debug with entries",
		log.Entry("entry1", "value1"),
		log.Entry("entry2", errors.New("error2")),
	)
	assert.Equal(t, `D debug with entries entry1="value1" entry2=error2`, strings.Split(buf.String(), "\n")[0])
	buf.Reset()
}
