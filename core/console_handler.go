// Package core
/***********************************************************************************************************************
* ProjectName:  gologger
* FileName:     console_handler.go
* Description:  控制台日志记录器
* Author:       ckf10000
* CreateDate:   2024/03/15 14:44:59
* Copyright ©2011-2024. Hunan xyz Company limited. All rights reserved.
* *********************************************************************************************************************/
package core

import "fmt"

// ConsoleLogger 控制台日志对象
type ConsoleLogger struct {
	BaseLogger
}

// newConsoleLogger 控制台日志头构造函数
func NewConsoleLogger(levelStr string) ConsoleLogger {
	baseLogger := NewBaseLogger(levelStr)
	return ConsoleLogger{
		BaseLogger: baseLogger,
	}
}

func (c ConsoleLogger) logHandler(lv LogLevel, format string, a ...interface{}) {
	var msg string = MsgFormatter(lv, format, a...)
	fmt.Print(msg)
}

func (c ConsoleLogger) logWrapper(level LogLevel, format string, a ...interface{}) {
	if c.enable(level) {
		c.logHandler(level, format, a...)
	}
}

// Debug 记录调试级别的日志
func (c ConsoleLogger) Debug(format string, a ...interface{}) {
	c.logWrapper(DEBUG, format, a...)
}

func (c ConsoleLogger) Trace(format string, a ...interface{}) {
	c.logWrapper(TRACE, format, a...)
}

func (c ConsoleLogger) Info(format string, a ...interface{}) {
	c.logWrapper(INFO, format, a...)
}

func (c ConsoleLogger) Warning(format string, a ...interface{}) {
	c.logWrapper(WARNING, format, a...)
}

func (c ConsoleLogger) Error(format string, a ...interface{}) {
	c.logWrapper(ERROR, format, a...)
}

func (c ConsoleLogger) Fatal(format string, a ...interface{}) {
	c.logWrapper(FATAL, format, a...)
}
