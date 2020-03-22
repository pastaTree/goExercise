package mylogger

import (
	"fmt"
	"time"
)

// Logger 日志结构体
type consoleLogger struct {
	Level logLevel
}

// NewConsoleLog 构造函数
func NewConsoleLog(levelStr string) consoleLogger {
	level, err := parseLogString(levelStr)
	if err != nil {
		panic(err)
	}
	return consoleLogger{
		Level: level,
	}
}

func (c consoleLogger) enable(logLevel logLevel) bool {
	return c.Level <= logLevel
}

func (c consoleLogger) log(level logLevel, format string, a ...interface{}) {
	if c.enable(level) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now()
		funcName, fileName, lineNum := getInfo(3)
		levelStr := parseLogLevel(level)
		fmt.Printf("[%s] [%s] [File:%s, Func:%s, Line:%d] %s\n", now.Format("2006-01-01 01:02:03"), levelStr, fileName, funcName, lineNum, msg)
	}
}

func (c consoleLogger) Debug(format string, a ...interface{}) {
	c.log(DEBUG, format, a...)
}

func (c consoleLogger) Trace(format string, a ...interface{}) {
	c.log(TRACE, format, a...)
}

func (c consoleLogger) Info(format string, a ...interface{}) {
	c.log(INFO, format, a...)
}

func (c consoleLogger) Warning(format string, a ...interface{}) {
	c.log(WARNING, format, a...)
}

func (c consoleLogger) Error(format string, a ...interface{}) {
	c.log(ERROR, format, a...)
}

func (c consoleLogger) Fatal(format string, a ...interface{}) {
	c.log(FATAL, format, a...)
}
