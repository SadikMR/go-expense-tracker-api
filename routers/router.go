package routers

import (
	"github.com/SadikMR/go-expense-tracker-api/controllers"
	"github.com/SadikMR/go-expense-tracker-api/middleware"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Public routes — no middleware
	beego.Router("/api/v1/health", &controllers.HealthController{}, "get:Get")
	beego.Router("/api/v1/auth/register", &controllers.AuthController{}, "post:Register")
	beego.Router("/api/v1/auth/login", &controllers.AuthController{}, "post:Login")

	// Protected routes — auth middleware applied to all /expenses* paths
	// IMPORTANT: /expenses/summary must be registered before /expenses/:id (Day 3)
	beego.InsertFilter("/api/v1/expenses*", beego.BeforeRouter, middleware.AuthMiddleware)
}
