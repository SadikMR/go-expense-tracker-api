package utils

import (
	beego "github.com/beego/beego/v2/server/web"
)

// AppConfig holds application-level configuration loaded from app.conf.
// Only values that genuinely differ between environments belong here.
var AppConfig = struct {
	AppName  string
	HTTPPort string
	RunMode  string
}{}

// InitConfig loads all values from app.conf into AppConfig.
// Must be called once at the top of main(), before InitLogger.
func InitConfig() {
	AppConfig.AppName = beego.AppConfig.DefaultString("appname", "go-expense-tracker-api")
	AppConfig.HTTPPort = beego.AppConfig.DefaultString("httpport", "8080")
	AppConfig.RunMode = beego.AppConfig.DefaultString("runmode", "dev")
}
