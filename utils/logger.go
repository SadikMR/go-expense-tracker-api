package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/core/logs"
)

// InitLogger sets up production-grade logging
func InitLogger(runMode string) {
	logDir := loggerDir()

	if runMode == "prod" {
		// File logging in production only
		os.MkdirAll(logDir, 0755)
		logFile := filepath.Join(logDir, "app.log")
		logs.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"%s","daily":true,"maxdays":7}`, logFile))
		logs.SetLevel(logs.LevelWarning)

	} else {
		// Console logging in development
		logs.SetLogger(logs.AdapterConsole)
		logs.SetLevel(logs.LevelDebug)
	}

	// Helps show correct file + line number
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	logs.Info("[Logger] initialized in %s mode", runMode)
}

func loggerDir() string {
	if AppConfig.LogDir != "" {
		return AppConfig.LogDir
	}
	return "logs"
}
