// Package log
/***********************************************************************************************************************
* ProjectName:  gologger
* FileName:     file_handler.go
* Description:  文件日志记录器
* Author:       ckf10000
* CreateDate:   2024/03/15 14:15:17
* Copyright ©2011-2024. Hunan xyz Company limited. All rights reserved.
* *********************************************************************************************************************/
package log

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

// FileLogger 文件日志记录器
type FileLogger struct {
	Level           LogLevel   // 日志级别
	FilePath        string     // 日志目录
	FileName        string     // 日志文件名，不包含目录
	MaxFileSize     int64      // 日志文件最大大小
	ErrFileName     string     // 错误日志文件名，不包含目录
	EnableErrOutput bool       // 错误级别以上日志是否独立输出到文档
	ConsoleOutput   bool       // 是否输出到控制台
	FileOutput      bool       // 是否输出到文件
	File            *os.File   // 日志文件句柄
	ErrFile         *os.File   // 错误日志文件句柄
	Lock            sync.Mutex // 互斥锁
	FormatTemplate  string     // 日志输出模版
}

// GetLogFile 获取日志文件名，全路径
func (f *FileLogger) GetLogFileFullName() string {
	return path.Join(f.FilePath, f.FileName)
}

// GetErrLogFileName 获取错误日志文件名
func (f *FileLogger) GetErrLogFileName() string {
	return f.FileName + ".err"
}

// GetErrLogFile 获取错误日志文件名，全路径
func (f *FileLogger) GetErrLogFileFullName() string {
	return path.Join(f.FilePath, f.GetErrLogFileName())
}

// GetLogBakFile 获取备份日志文件名，全路径
func (f *FileLogger) GetLogBakFileFullName(suffix string) string {
	return fmt.Sprintf("%s.bak%s", f.GetLogFileFullName(), suffix)
}

// GetErrLogBakFile 获取备份错误日志文件名，全路径
func (f *FileLogger) GetErrLogBakFileFullName(suffix string) string {
	return fmt.Sprintf("%s.bak%s", f.GetErrLogFileFullName(), suffix)
}

func (f *FileLogger) enable(logLevel LogLevel) bool {
	return logLevel >= f.Level
}

func (f *FileLogger) logWrapper(level LogLevel, format string, args ...interface{}) {
	if f.enable(level) {
		f.LogHandler(level, format, args...)
	}
}

// Debug 记录调试级别的日志
func (f *FileLogger) Debug(format string, args ...interface{}) {
	f.logWrapper(DEBUG, format, args...)
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	f.logWrapper(TRACE, format, args...)
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	f.logWrapper(INFO, format, args...)
}

func (f *FileLogger) Warning(format string, args ...interface{}) {
	f.logWrapper(WARNING, format, args...)
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	f.logWrapper(ERROR, format, args...)
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	f.logWrapper(FATAL, format, args...)
}

// 初始化日志文件
func (f *FileLogger) initLogFile() error {
	// 创建日志文件
	logFileName := f.GetLogFileFullName()
	file, err := CreateFile(logFileName)
	if err != nil {
		return err
	}
	// 设置日志文件句柄
	f.File = file
	if f.EnableErrOutput {
		f.ErrFileName = f.GetErrLogFileName()
		errLogFileName := f.GetErrLogFileFullName()
		errFile, err := CreateFile(errLogFileName)
		if err != nil {
			return err
		}
		// 设置错误日志文件句柄
		f.ErrFile = errFile
	}
	return nil
}

// log 日志记录方法
func (f *FileLogger) LogHandler(level LogLevel, format string, args ...interface{}) {
	now := time.Now()
	funcName, fileName, lineNo := GetInfo(4)
	logEntry := LogEntry{
		Time:      now,
		Message:   fmt.Sprintf(format, args...),
		Level:     level,
		FileName:  fileName,
		Line:      lineNo,
		FuncName:  funcName,
		ProcessId: os.Getpid(),
		ThreadId:  20,
	}

	// 输出到控制台
	if f.ConsoleOutput {
		fmt.Print(FormatLogEntry(f.FormatTemplate, &logEntry))
	}

	// 输出到日志文件
	if f.FileOutput {
		suffix := strings.Replace(now.Format(DATEIDFORMATTER), ".", "", 1)
		f.Lock.Lock()
		defer f.Lock.Unlock()
		if CheckFileSize(f.File, f.MaxFileSize) {
			// 需要切割日志文件
			fileObj, err := RotatingFile(f.File, f.GetLogFileFullName(), f.GetLogBakFileFullName(suffix))
			if err != nil {
				fmt.Printf("open new log file: %s failed, err: %s\n", f.FileName, err)
				return
			}
			// 4. 将打开的新日志文件对象赋值给f.FileObj
			f.File = fileObj
		}
		WriteFile(f.File, FormatLogEntry(f.FormatTemplate, &logEntry))
		if level >= ERROR {
			if CheckFileSize(f.ErrFile, f.MaxFileSize) {
				// 需要切割日志文件
				errFileObj, err := RotatingFile(f.ErrFile, f.GetErrLogFileFullName(), f.GetErrLogBakFileFullName(suffix))
				if err != nil {
					fmt.Printf("open new log file: %s failed, err: %s\n", f.ErrFileName, err)
					return
				}
				// 4. 将打开的新日志文件对象赋值给f.ErrFileObj
				f.ErrFile = errFileObj
			}
			WriteFile(f.ErrFile, FormatLogEntry(f.FormatTemplate, &logEntry))
		}
	}
}

// NewFileLogger 创建文件日志记录器
func NewLogger(levelStr, logFilePath, logFileName, formatTemplate string, maxFileSize int64, enableErrOutput, consoleOutput, fileOutput bool) *FileLogger {
	logLevel, err := ParseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	if logFilePath != "" {
		err := CreateDirectoryIfNotExists(logFilePath)
		if err != nil {
			panic(err)
		}
	}
	if formatTemplate != SIMPLETEMPLATE && formatTemplate != STANDARDTEMPLATE {
		panic(fmt.Sprintf("Only supported [%s, %s]", SIMPLETEMPLATE, STANDARDTEMPLATE))
	}
	fileLogger := &FileLogger{
		Level:           logLevel,
		FilePath:        logFilePath,
		FileName:        logFileName,
		MaxFileSize:     maxFileSize,
		EnableErrOutput: enableErrOutput,
		ConsoleOutput:   consoleOutput,
		FileOutput:      fileOutput,
		FormatTemplate:  formatTemplate,
	}

	// 创建日志文件
	if fileOutput {
		err := fileLogger.initLogFile()
		if err != nil {
			panic(err)
		}
	}

	return fileLogger
}
