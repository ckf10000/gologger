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
	"testing"

	"github.com/ckf10000/gologger/v1/log"
)

func TestFileLoggerDemo01(t *testing.T) {
	projectPath := log.GetProjectAbsPath()
	// projectPath := log.GetExecuteFilePath()
	// projectPath := "./"
	log := log.NewLogger("info", projectPath, "app.log", "standard", 50*1024*1024, true, true, true)
	for i := 0; i < 10000; i++ {
		log.Debug("这是一次打印：debug")
		log.Trace("这是一次打印：trace")
		log.Info("这是一次打印：info")
		log.Warning("这是一次打印：warning")
		log.Error("这是一次打印：error")
		log.Fatal("这是一次打印：fatal")
	}
}
