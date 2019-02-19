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

func TestLeveled(t *testing.T) {

	var testData = []struct {
		testName     string
		verbosity    int
		msg          string
		entryLevel   int
		outputExpect bool
	}{
		{
			testName:     "simple",
			verbosity:    5,
			msg:          "log test",
			entryLevel:   1,
			outputExpect: true,
		},
		{
			testName:     "level equals verbosity",
			verbosity:    5,
			msg:          "log test",
			entryLevel:   5,
			outputExpect: true,
		},
		{
			testName:     "level greater than verbosity",
			verbosity:    5,
			msg:          "log test",
			entryLevel:   6,
			outputExpect: false,
		},
	}

	for _, td := range testData {
		lgrus := logrus.New()
		var buffer bytes.Buffer
		lgrus.Out = &buffer
		l := NewLogger(lgrus, td.verbosity, "")

		l.V(td.entryLevel).Info(td.msg)
		if td.outputExpect {
			assert.Contains(t, buffer.String(), logrusMessageLabel+"=\""+td.msg+"\"")
		} else {
			assert.Zero(t, len(buffer.String()))
		}
	}
}

func TestSetLevel(t *testing.T) {
	lgrus := logrus.New()
	l := NewLogger(lgrus, 1, "")
	assert.EqualValues(t, l.verbosity, 1)

	l.SetLevel(2)
	assert.EqualValues(t, l.verbosity, 2)
}
