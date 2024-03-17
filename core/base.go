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
const DATEIDFORMATTER string = "20060102150405.000"
const MSGSTANDARDFORMATTER string = "%s - [PID-%d] - [Thread-%d] - [%s] - %s - <%s> - [Line-%d] - %s\n"

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
	Debug(format string, args ...interface{})
	Trace(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

// LogEntry 日志条目结构
type LogEntry struct {
	Time      time.Time // 日志记录时间
	Message   string    // 日志内容
	Level     LogLevel  // 日志级别
	FileName  string    // 日志文件名
	Line      int       // 代码行数
	FuncName  string    // 函数名/方法名
	ProcessId int       // 进程ID
	ThreadId  int       // 线程ID
}

// ParseLogLevel 将字符串日志级别，转换成 LogLevel 类型
func ParseLogLevel(levelStr string) (LogLevel, error) {
	strTemp := strings.ToUpper(levelStr)
	level, ok := stringToLogLevel[strTemp]
	if !ok {
		return UNKNOWN, errors.New("invalid log level")
	}
	return level, nil
}

// GetLogString 获取LogLevel类型的日志级别对应的字符串
func GetLogString(lv LogLevel) string {
	return logLevelToString[lv]
}

// formatLogEntry 格式化日志条目
func FormatLogEntry(entry *LogEntry) string {
	return fmt.Sprintf(
		MSGSTANDARDFORMATTER, entry.Time.Format(DATEFORMATTER), entry.ProcessId, entry.ThreadId,
		GetLogString(entry.Level), entry.Message, entry.FuncName, entry.Line, entry.FileName)
}
