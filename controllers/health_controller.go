package controllers

import (
	"github.com/SadikMR/go-expense-tracker-api/utils"

	beego "github.com/beego/beego/v2/server/web"
)

// HealthController handles the health check endpoint.
type HealthController struct {
	beego.Controller
}

// Get returns 200 if the server is running.
// @Summary      Health check
// @Description  Returns 200 if the server is running
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /api/v1/health [get]
func (c *HealthController) Get() {
	utils.Success(c.Ctx, 200, "Server is running", nil)
}
