// Package core
/***********************************************************************************************************************
* ProjectName:  gologger
* FileName:     utils.go
* Description:  工具模块
* Author:       ckf10000
* CreateDate:   2024/03/15 15:27:05
* Copyright ©2011-2024. Hunan xyz Company limited. All rights reserved.
* *********************************************************************************************************************/
package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func GetInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, filePath, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Printf("runtme.Caller() call failed.\n")
		return
	}
	funcNameSlice := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	funcName = funcNameSlice[len(funcNameSlice)-1]
	fileName = path.Base(filePath)
	return
}

// GetProjectAbsPath 获取项目路径
func GetProjectAbsPath() string {
	// 获取当前源码文件的路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get current file path")
	}

	// 将文件路径转换为绝对路径
	absPath, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}
	// 获取当前文件所在目录的路径
	dir := filepath.Dir(absPath)
	// ../project_name/common
	//fmt.Println(dir)
	// 获取项目根路径（假设项目根路径是 src 目录的父目录）
	root := filepath.Join(dir, "..")
	return root
}

// CheckFileSize 校验文件大小
func CheckFileSize(fileObj *os.File, size int64) bool {
	fileInfo, err := fileObj.Stat()
	if err != nil {
		fmt.Println("Get file info failed, err: ", err)
		return false
	}
	return fileInfo.Size() >= size
}

// 切割文件
func RotatingFile(fileObj *os.File, filePath string) (*os.File, error) {
	// 需要切割日志文件
	nowStr := time.Now().Format("20060102150405000")
	fileInfo, err1 := fileObj.Stat()
	if err1 != nil {
		fmt.Println("get file info failed, err: ", err1)
		return nil, err1
	}
	logName := path.Join(filePath, fileInfo.Name())        // 拿到当前日志文件的完整路径
	newLogName := fmt.Sprintf("%s.bak%s", logName, nowStr) // 拼接一个日志文件备份的名字
	// 1. 关闭当前日志文件
	fileObj.Close()
	// 2. 备份一下 rename xx.log --> xx.log.bak2024031617091
	os.Rename(logName, newLogName)
	// 3. 打开一个新的日志文件
	fileObj, err2 := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err2 != nil {
		fmt.Println("open new log file failed, err: ", err2)
		return nil, err2
	}
	return fileObj, nil
}
