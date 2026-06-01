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

	// Protected routes
	beego.InsertFilter("/api/v1/expenses*", beego.BeforeRouter, middleware.AuthMiddleware)

	// summary MUST be first — before :id
	beego.Router("/api/v1/expenses/summary", &controllers.SummaryController{}, "get:Get")
	beego.Router("/api/v1/expenses", &controllers.ExpenseController{}, "post:Post;get:Get")
	beego.Router("/api/v1/expenses/:id", &controllers.ExpenseController{}, "get:Get;put:Put;delete:Delete")
}
