# Logrus Wrapper to Log interface <img align="left" width="60px" src="https://avatars0.githubusercontent.com/u/47711035?s=400&u=e8a2891cca67da66972ad478069588deb0299e4b&v=4">

[![Build Status](https://travis-ci.com/go-interfaces/logrusw.svg?branch=master)](https://travis-ci.com/go-interfaces/logrusw) [![codecov](https://codecov.io/gh/go-interfaces/logrusw/branch/master/graph/badge.svg)](https://codecov.io/gh/go-interfaces/logrusw)

[Logrus](https://github.com/sirupsen/logrus) wrapper for [Log interface](https://github.com/go-interfaces/log).


## Usage


Initialize log

```go
import (
    "github.com/go-interfaces/log"
    "github.com/go-interfaces/logwrapper"
    "github.com/sirupsen/logrus"
)	
	
func RegisterLogrusLogger() {
    l := logrus.New()
    // ... configure logrus logger ...

    // verbosity level, most probably taken from flags or environment
    verbosity = 5

    // logger name, if not empty will be added to each log entry
    loggerName = ""

    logger := logrusw.NewLogger(l, verbosity, loggerName)

    // could also use log.SetLogger(name, logger) for non default loggers
    log.SetDefaultLogger(logger)
}

```

Use log
```go
import (
    "github.com/go-interfaces/log"
 )	

// DoSomenthing uses the default logger
func DoSomething(param string) {
    log.Info("about to do something...")

    if err := something(param); err != nil {
        log.Error(err, "error trying to do something", log.KV{"param", param})
    }
}
```
