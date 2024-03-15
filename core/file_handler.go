// Package core
/***********************************************************************************************************************
* ProjectName:  gologger
* FileName:     file_handler.go
* Description:  文件日志记录器
* Author:       ckf10000
* CreateDate:   2024/03/15 14:15:17
* Copyright ©2011-2024. Hunan xyz Company limited. All rights reserved.
* *********************************************************************************************************************/
package core

import (
	"fmt"
	"os"
	"path"
)

type FileLogger struct {
	BaseLogger
	FilePath    string   // 日志文件保存的路径
	FileName    string   // 日志文件名
	FileObj     *os.File // 文件对象
	ErrFileObj  *os.File // 文件对象
	MaxFileSize int64    // 文件大小
}

// NewFileLogger 文件类日志对象构造函数
func NewFileLogger(LevelStr, filePath, fileName string, maxFileSize int64) *FileLogger {
	baseLogger := NewBaseLogger(LevelStr)
	fileLogger := &FileLogger{
		BaseLogger:  baseLogger,
		FilePath:    filePath,
		FileName:    fileName,
		MaxFileSize: maxFileSize,
	}
	err := fileLogger.initFile() // 按照文件路径和文件名将文件打开
	if err != nil {
		panic(err)
	}
	return fileLogger
}

func (f *FileLogger) initFile() error {
	fullFileName := path.Join(f.FilePath, f.FileName)
	fileObj, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open log file failed, err: %v\n", err)
		return err
	}
	errFileObj, err := os.OpenFile(fullFileName+".err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open err log file failed, err: %v\n", err)
		return err
	}
	f.FileObj = fileObj
	f.ErrFileObj = errFileObj
	return nil
}

// closeFile 关闭日志文件
func (f *FileLogger) closeFile() {
	f.FileObj.Close()
	f.ErrFileObj.Close()
}

func (f FileLogger) logHandler(lv LogLevel, format string, a ...interface{}) {
	var msg string = MsgFormatter(lv, format, a...)
	if CheckFileSize(f.FileObj, f.MaxFileSize) {
		// 需要切割日志文件
		fileObj, err := RotatingFile(f.FileObj, f.FilePath)
		if err != nil {
			fmt.Println("open new log file failed, err: ", err)
			return
		}
		// 4. 将打开的新日志文件对象赋值给f.FileObj
		f.FileObj = fileObj
	}
	fmt.Fprint(f.FileObj, msg)
	if lv >= ERROR {
		if CheckFileSize(f.ErrFileObj, f.MaxFileSize) {
			// 需要切割日志文件
			errFileObj, err := RotatingFile(f.ErrFileObj, f.FilePath)
			if err != nil {
				fmt.Println("open new log file failed, err: ", err)
				return
			}
			// 4. 将打开的新日志文件对象赋值给f.ErrFileObj
			f.ErrFileObj = errFileObj
		}
		// 将大于或等于ERROR级别的日志，单独记录到err文件中
		fmt.Fprint(f.ErrFileObj, msg)
	}

}

func (f FileLogger) logWrapper(level LogLevel, format string, a ...interface{}) {
	if f.enable(level) {
		f.logHandler(level, format, a...)
	}
}

// Debug 记录调试级别的日志
func (f FileLogger) Debug(format string, a ...interface{}) {
	f.logWrapper(DEBUG, format, a...)
}

func (f FileLogger) Trace(format string, a ...interface{}) {
	f.logWrapper(TRACE, format, a...)
}

func (f FileLogger) Info(format string, a ...interface{}) {
	f.logWrapper(INFO, format, a...)
}

func (f FileLogger) Warning(format string, a ...interface{}) {
	f.logWrapper(WARNING, format, a...)
}

func (f FileLogger) Error(format string, a ...interface{}) {
	f.logWrapper(ERROR, format, a...)
}

func (f FileLogger) Fatal(format string, a ...interface{}) {
	f.logWrapper(FATAL, format, a...)
}
