// Package tests
/***********************************************************************************************************************
* ProjectName:  gologger
* FileName:     console_test.go
* Description:  TODO
* Author:       ckf10000
* CreateDate:   2024/03/15 16:07:02
* Copyright ©2011-2024. Hunan xyz Company limited. All rights reserved.
* *********************************************************************************************************************/
package tests

import (
	"gologger/core"
	"testing"
)

// TestConsoleLoggerDemo01 支持级别的控制
func TestConsoleLoggerDemo01(t *testing.T) {
	log := core.NewConsoleLogger("ERROR")
	log.Debug("这是一次打印：debug")
	log.Trace("这是一次打印：trace")
	log.Info("这是一次打印：info")
	log.Warning("这是一次打印：warning")
	log.Error("这是一次打印：error")
	log.Fatal("这是一次打印：fatal")
}

// TestConsoleLoggerDemo02 支持格式化打印
func TestConsoleLoggerDemo02(t *testing.T) {
	log := core.NewConsoleLogger("debug")
	log.Error("这是一次打印：%s, %s", "Hello", "World")
}
