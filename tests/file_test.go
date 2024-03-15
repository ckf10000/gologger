// Package tests
/***********************************************************************************************************************
* ProjectName:  gologger
* FileName:     file_test.go
* Description:  TODO
* Author:       ckf10000
* CreateDate:   2024/03/16 01:37:16
* Copyright ©2011-2024. Hunan xyz Company limited. All rights reserved.
* *********************************************************************************************************************/
package tests

import (
	"gologger/core"
	"testing"
)

func TestFileLoggerDemo01(t *testing.T) {
	// projectPath := core.GetProjectAbsPath()
	projectPath := "./"
	log := core.NewFileLogger("DEBUG", projectPath, "test.log", 10*1024)
	for i := 0; i < 10; i++ {
		log.Debug("这是一次打印：debug")
		log.Trace("这是一次打印：trace")
		log.Info("这是一次打印：info")
		log.Warning("这是一次打印：warning")
		log.Error("这是一次打印：error")
		log.Fatal("这是一次打印：fatal")
	}
}
