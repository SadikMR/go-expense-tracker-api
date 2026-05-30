package main

import (
	"github.com/SadikMR/go-expense-tracker-api/models"
	_ "github.com/SadikMR/go-expense-tracker-api/routers"
	"github.com/SadikMR/go-expense-tracker-api/utils"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// 1. Load config first — all other layers depend on it
	utils.InitConfig()

	// 2. Init logger second — based on runmode from config
	utils.InitLogger(utils.AppConfig.RunMode)

	// 3. Ensure CSV storage is ready before serving any request
	if err := models.EnsureUsersCSV(); err != nil {
		logs.Critical("[Startup] Failed to initialise users.csv: %v", err)
		return
	}

	logs.Info("[Startup] %s starting on port %s in %s mode",
		utils.AppConfig.AppName,
		utils.AppConfig.HTTPPort,
		utils.AppConfig.RunMode,
	)

	// 4. Start server
	beego.Run()
}
