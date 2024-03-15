// Package core
/***********************************************************************************************************************
* ProjectName:  gologger
* FileName:     base.go
* Description:  TODO
* Author:       ckf10000
* CreateDate:   2024/03/15 14:29:35
* Copyright ©2011-2024. Hunan xyz Company limited. All rights reserved.
* *********************************************************************************************************************/
package core

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// LogLevel 日志级别类型
type LogLevel uint16

// 告警级别常量
const (
	UNKNOWN LogLevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

const DATEFORMATTER string = "2006-01-02 15:04:05.000"
const MSGSTANDARDFORMATTER string = "%s - [PID-10] - [Thread-20] - [%s] - %s - <%s> - [Line-%d] - %s\n"

// 告警级别常量映射
var logLevelToString = map[LogLevel]string{
	UNKNOWN: "UNKNOWN",
	DEBUG:   "DEBUG",
	TRACE:   "TRACE",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

var stringToLogLevel = map[string]LogLevel{}

func init() {
	for key, value := range logLevelToString {
		stringToLogLevel[value] = key
	}
}

// Logger 接口定义了日志记录器的行为
type Logger interface {
	Debug(message string, a ...interface{})
	Trace(message string, a ...interface{})
	Info(message string, a ...interface{})
	Warning(message string, a ...interface{})
	Error(message string, a ...interface{})
	Fatal(message string, a ...interface{})
	logHandler(lv LogLevel, format string, a ...interface{})
}

// BaseLogger 通用的日志记录器
type BaseLogger struct {
	Level LogLevel // 日志级别
}

func NewBaseLogger(levelStr string) BaseLogger {
	logLevel, err := ParseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	return BaseLogger{
		Level: logLevel,
	}
}

func ParseLogLevel(levelStr string) (LogLevel, error) {
	strTemp := strings.ToUpper(levelStr)
	level, ok := stringToLogLevel[strTemp]
	if !ok {
		return UNKNOWN, errors.New("无效的日志级别")
	}
	return level, nil
}

func getLogString(lv LogLevel) string {
	return logLevelToString[lv]
}

func (b BaseLogger) enable(logLevel LogLevel) bool {
	return logLevel >= b.Level
}

func MsgFormatter(lv LogLevel, format string, a ...interface{}) string {
	msg := fmt.Sprintf(format, a...)
	now := time.Now()
	funcName, fileName, lineNo := GetInfo(5)
	formattedString := fmt.Sprintf(MSGSTANDARDFORMATTER, now.Format(DATEFORMATTER), getLogString(lv), msg, funcName, lineNo, fileName)
	return formattedString
}
