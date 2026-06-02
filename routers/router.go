package routers

import (
	"github.com/SadikMR/go-expense-tracker-api/controllers"
	"github.com/SadikMR/go-expense-tracker-api/middleware"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.InsertFilter("/api/v1/expenses*", beego.BeforeRouter, middleware.AuthMiddleware)

	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/health", &controllers.HealthController{}, "get:Get"),
		beego.NSNamespace("/auth",
			beego.NSRouter("/register", &controllers.AuthController{}, "post:Register"),
			beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
		),
		beego.NSNamespace("/expenses",
			beego.NSRouter("/summary", &controllers.SummaryController{}, "get:Get"),
			beego.NSRouter("", &controllers.ExpenseController{}, "post:Post;get:Get"),
			beego.NSRouter("/:id", &controllers.ExpenseController{}, "get:GetOne;put:Put;delete:Delete"),
		),
	)

	beego.AddNamespace(ns)
}
