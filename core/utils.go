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
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
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
		fmt.Println("Get project path failed.")
		return ""
	}

	// 将文件路径转换为绝对路径
	absPath, err := filepath.Abs(filename)
	if err != nil {
		fmt.Println("Get project path failed.")
		return ""
	}
	// 获取当前文件所在目录的路径
	dir := filepath.Dir(absPath)
	// ../project_name/common
	//fmt.Println(dir)
	// 获取项目根路径（假设项目根路径是 src 目录的父目录）
	root := filepath.Join(dir, "..")
	return root
}

// CreateDirectoryIfNotExists 判断目录是否存在，如果不存在，则递归式逐级创建，支持windows和linux
func CreateDirectoryIfNotExists(filePath string) error {
	// 检查目录是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// 如果目录不存在则创建
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
		fmt.Printf("Directory '%s' created successfully.\n", filePath)
	} else if err != nil {
		return err
	}
	return nil
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

// CreateFile 创建文件
func CreateFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		err := fmt.Sprintf("Failed to open log file %s, error: %s\n", fileName, err)
		return nil, errors.New(err)
	}
	return file, nil
}

// WriteFile 保存至文件
func WriteFile(file *os.File, format string) {
	_, err := fmt.Fprint(file, format)
	// _, err := file.WriteString(format)
	if err != nil {
		fmt.Printf("Failed to write log to file, error: %s\n", err)
	}
}

// 切割文件
func RotatingFile(fileObj *os.File, fileName, bakFileName string) (*os.File, error) {
	// 需要切割日志文件
	// 1. 关闭当前日志文件
	fileObj.Close()
	// 2. 备份一下 rename xx.log --> xx.log.bak2024031617091
	os.Rename(fileName, bakFileName)
	// 3. 打开一个新的日志文件
	return CreateFile(fileName)
}
