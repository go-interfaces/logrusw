package logrusw

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/go-interfaces/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	logrusErrorLabel   = "error"
	logrusMessageLabel = "msg"
)

func TestLogger(t *testing.T) {

	var testData = []struct {
		testName  string
		verbosity int
		logName   string
		msg       string
		errorMsg  string
		kv        []log.KV
	}{
		{
			testName:  "simple",
			verbosity: 5,
			logName:   "",
			msg:       "log test",
			errorMsg:  "something bad",
		},
		{
			testName:  "with logger name output",
			verbosity: 5,
			logName:   "named-logger",
			msg:       "log test",
			errorMsg:  "something bad",
		},
		{
			testName:  "with level output",
			verbosity: 5,
			logName:   "",
			msg:       "log test",
			errorMsg:  "something bad",
		},
		{
			testName:  "with KV",
			verbosity: 5,
			logName:   "",
			msg:       "log test",
			errorMsg:  "something bad",
			kv:        []log.KV{{K: "fieldA", V: "valueA"}, {K: "fieldB", V: "valueB"}},
		},
	}

	for _, td := range testData {

		lgrus := logrus.New()
		var buffer bytes.Buffer
		lgrus.Out = &buffer

		l := NewLogger(lgrus, td.verbosity, td.logName)

		l.Info(td.msg, td.kv...)
		assert.Contains(t, buffer.String(), logrusMessageLabel+"=\""+td.msg+"\"")
		if len(td.logName) != 0 {
			assert.Contains(t, buffer.String(), loggerLabel+"="+td.logName)
		}
		if len(td.kv) != 0 {
			for _, kv := range td.kv {
				assert.Contains(t, buffer.String(), fmt.Sprintf("%s=%v", kv.K, kv.V))
			}
		}

		buffer.Reset()
		l.Error(errors.New(td.errorMsg), td.msg, td.kv...)
		assert.Contains(t, buffer.String(), logrusMessageLabel+"=\""+td.msg+"\"")
		assert.Contains(t, buffer.String(), logrusErrorLabel+"=\""+td.errorMsg+"\"")
		if len(td.logName) != 0 {
			assert.Contains(t, buffer.String(), loggerLabel+"="+td.logName)
		}
		if len(td.kv) != 0 {
			for _, kv := range td.kv {
				assert.Contains(t, buffer.String(), fmt.Sprintf("%s=%v", kv.K, kv.V))
			}
		}
	}
}
