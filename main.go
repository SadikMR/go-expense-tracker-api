// @title           Go Expense Tracker API
// @version         1.0
// @description     RESTful API for tracking expenses with authentication, expense CRUD, and summary reporting.
// @contact.name    API Support
// @contact.email   support@example.com
// @host            localhost:8080
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-User-ID
//go:generate swag init -g main.go -o docs
package main

import (
	"github.com/SadikMR/go-expense-tracker-api/models"
	_ "github.com/SadikMR/go-expense-tracker-api/routers"
	_ "github.com/SadikMR/go-expense-tracker-api/docs"
	"github.com/SadikMR/go-expense-tracker-api/utils"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// 1. Load config first — all other layers depend on it
	utils.InitConfig()

	// 2. Init logger — based on runmode
	utils.InitLogger(utils.AppConfig.RunMode)

	// 3. Ensure CSV storage is ready before serving any request
	if err := models.EnsureUsersCSV(); err != nil {
		logs.Critical("[Startup] Failed to initialise users.csv: %v", err)
		return
	}
	if err := models.EnsureExpensesCSV(); err != nil {
		logs.Critical("[Startup] Failed to initialise expenses.csv: %v", err)
		return
	}

	logs.Info("[Startup] %s starting on port %s in %s mode",
		utils.AppConfig.AppName,
		utils.AppConfig.HTTPPort,
		utils.AppConfig.RunMode,
	)

	// Serve Swagger UI from generated annotations
	beego.Handler("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	// 4. Start server
	beego.Run()
}
