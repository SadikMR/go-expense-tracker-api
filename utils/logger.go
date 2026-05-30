package utils

import (
	"os"

	"github.com/beego/beego/v2/core/logs"
)

// InitLogger sets up production-grade logging
func InitLogger(runMode string) {

	// Ensure logs folder exists
	os.MkdirAll("logs", 0755)

	if runMode == "prod" {
		// File logging (production)
		logs.SetLogger(logs.AdapterFile, `{
			"filename": "logs/app.log",
			"daily": true,
			"maxdays": 7
		}`)

		logs.SetLevel(logs.LevelWarning)

	} else {
		// Console logging (development)
		logs.SetLogger(logs.AdapterConsole)
		logs.SetLevel(logs.LevelDebug)
	}

	// Helps show correct file + line number
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	logs.Info("[Logger] initialized in %s mode", runMode)
}
