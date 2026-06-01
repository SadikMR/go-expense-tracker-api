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
// Both date_from and date_to are optional query parameters.
//
// @Summary      Expense summary
// @Description  Returns total amount, count and breakdown by category with optional date range filter
// @Tags         expenses
// @Produce      json
// @Param        X-User-ID  header    int     true   "User ID"
// @Param        date_from  query     string  false  "Filter from date (YYYY-MM-DD)"
// @Param        date_to    query     string  false  "Filter to date (YYYY-MM-DD)"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
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
