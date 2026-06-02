package controllers

import (
	"github.com/SadikMR/go-expense-tracker-api/services"
	"github.com/SadikMR/go-expense-tracker-api/utils"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// SummaryController handles the expense summary endpoint.
type SummaryController struct {
	beego.Controller
}

// Get returns a summary of expenses for the authenticated user.
// @Summary      Expense summary
// @Description  Returns total amount, total count, and breakdown by category for the authenticated user. Date filters are inclusive.
// @Tags         expenses
// @Produce      json
// @Security     ApiKeyAuth
// @Param        X-User-ID  header    int     true   "Authenticated user ID (returned after login). Example: 123"
// @Param        date_from  query     string  false  "Filter from date in YYYY-MM-DD format"
// @Param        date_to    query     string  false  "Filter to date in YYYY-MM-DD format"
// @Success      200        {object}  SummaryResponse  "Summary generated successfully"
// @Failure      400        {object}  ErrorResponse    "Validation error"
// @Failure      401        {object}  ErrorResponse    "Unauthorized"
// @Failure      500        {object}  ErrorResponse    "Internal server error"
// @Router       /api/v1/expenses/summary [get]
func (c *SummaryController) Get() {
	userID, ok := userIDFromHeader(&c.Controller)
	if !ok {
		utils.Error(c.Ctx, 401, "Unauthorized")
		return
	}

	dateFrom := c.GetString("date_from")
	dateTo := c.GetString("date_to")

	if dateFrom != "" && !utils.ValidateDate(dateFrom) {
		utils.Error(c.Ctx, 400, "date_from must be in YYYY-MM-DD format")
		return
	}
	if dateTo != "" && !utils.ValidateDate(dateTo) {
		utils.Error(c.Ctx, 400, "date_to must be in YYYY-MM-DD format")
		return
	}

	summary, err := services.GetSummary(userID, dateFrom, dateTo)
	if err != nil {
		logs.Error("[Summary] Failed: user_id=%d err=%v", userID, err)
		utils.Error(c.Ctx, 500, "Internal server error")
		return
	}

	logs.Info("[Summary] user_id=%d total_amount=%.2f total_count=%d",
		userID, summary.TotalAmount, summary.TotalCount,
	)
	utils.Success(c.Ctx, 200, "Summary generated", summary)
}
