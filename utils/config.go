package utils

import (
	"os"

	beego "github.com/beego/beego/v2/server/web"
)

// AppConfig holds application-level configuration loaded from app.conf.
// Only values that genuinely differ between environments belong here.
var AppConfig = struct {
	AppName         string
	HTTPPort        string
	RunMode         string
	DataDir         string
	UsersCSVPath    string
	ExpensesCSVPath string
	LogDir          string
}{}

// InitConfig loads all values from app.conf into AppConfig.
// Must be called once at the top of main(), before InitLogger.
func InitConfig() {
	AppConfig.AppName = beego.AppConfig.DefaultString("appname", "")
	if AppConfig.AppName == "" {
		AppConfig.AppName = os.Getenv("APP_NAME")
	}
	if AppConfig.AppName == "" {
		AppConfig.AppName = "go-expense-tracker-api"
	}

	AppConfig.HTTPPort = beego.AppConfig.DefaultString("httpport", "")
	if AppConfig.HTTPPort == "" {
		AppConfig.HTTPPort = os.Getenv("PORT")
	}

	AppConfig.RunMode = beego.AppConfig.DefaultString("runmode", "")
	if AppConfig.RunMode == "" {
		AppConfig.RunMode = os.Getenv("RUN_MODE")
	}
	if AppConfig.RunMode == "" {
		AppConfig.RunMode = "dev"
	}

	AppConfig.DataDir = beego.AppConfig.DefaultString("data_dir", "")
	if AppConfig.DataDir == "" {
		AppConfig.DataDir = os.Getenv("DATA_DIR")
	}

	AppConfig.UsersCSVPath = beego.AppConfig.DefaultString("users_csv_path", "")
	if AppConfig.UsersCSVPath == "" {
		AppConfig.UsersCSVPath = os.Getenv("USERS_CSV_PATH")
	}

	AppConfig.ExpensesCSVPath = beego.AppConfig.DefaultString("expenses_csv_path", "")
	if AppConfig.ExpensesCSVPath == "" {
		AppConfig.ExpensesCSVPath = os.Getenv("EXPENSES_CSV_PATH")
	}

	AppConfig.LogDir = beego.AppConfig.DefaultString("log_dir", "")
	if AppConfig.LogDir == "" {
		AppConfig.LogDir = os.Getenv("LOG_DIR")
	}
}
