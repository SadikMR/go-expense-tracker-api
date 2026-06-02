package controllers

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/SadikMR/go-expense-tracker-api/services"
	"github.com/SadikMR/go-expense-tracker-api/utils"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// ExpenseController handles expense CRUD operations.
type ExpenseController struct {
	beego.Controller
}

type expenseInput struct {
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
}

type expenseUpdateInput struct {
	Title       *string  `json:"title"`
	Amount      *float64 `json:"amount"`
	Category    *string  `json:"category"`
	Note        *string  `json:"note"`
	ExpenseDate *string  `json:"expense_date"`
}

// userIDFromHeader reads and parses the X-User-ID header.
func userIDFromHeader(c *beego.Controller) (int, bool) {
	return utils.ParseID(c.Ctx.Input.Header("X-User-ID"))
}

// handleServiceError writes the correct HTTP response based on error type.
// Validation errors (not wrapped) return 400 with the error message.
// Infrastructure errors (wrapped with %w) return 500.
func handleServiceError(ctx *beego.Controller, tag string, err error) {
	if errors.Unwrap(err) != nil {
		logs.Error("%s infrastructure error: %v", tag, err)
		utils.Error(ctx.Ctx, 500, "Internal server error")
	} else {
		logs.Warn("%s validation error: %v", tag, err)
		utils.Error(ctx.Ctx, 400, err.Error())
	}
}

// Post creates a new expense for the authenticated user.
// @Summary      Create expense
// @Description  Creates a new expense for the authenticated user
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header    int                     true  "User ID"
// @Param        body       body      expenseInput            true  "Expense payload"
// @Success      201        {object}  map[string]interface{}  "Expense created successfully"
// @Failure      400        {object}  map[string]interface{}  "Validation error"
// @Failure      401        {object}  map[string]interface{}  "Unauthorized"
// @Router       /api/v1/expenses [post]
func (c *ExpenseController) Post() {
	userID, ok := userIDFromHeader(&c.Controller)
	if !ok {
		utils.Error(c.Ctx, 401, "Unauthorized")
		return
	}

	body, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		utils.Error(c.Ctx, 400, "Failed to read request body")
		return
	}

	var input expenseInput
	if err := json.Unmarshal(body, &input); err != nil {
		utils.Error(c.Ctx, 400, "Invalid request body")
		return
	}

	expense, err := services.CreateExpense(userID, services.CreateExpenseInput{
		Title:       input.Title,
		Amount:      input.Amount,
		Category:    input.Category,
		Note:        input.Note,
		ExpenseDate: input.ExpenseDate,
	})
	if err != nil {
		handleServiceError(&c.Controller, "[Expense] Post", err)
		return
	}

	logs.Info("[Expense] Created: id=%d user_id=%d", expense.ID, userID)
	utils.Success(c.Ctx, 201, "Expense created successfully", expense)
}

// Get returns a list of expenses or a single expense.
// @Summary      List expenses
// @Description  Returns expenses for the authenticated user with optional filters and sorting
// @Tags         expenses
// @Produce      json
// @Param        X-User-ID   header    int                     true   "User ID"
// @Param        date_from   query     string                  false  "Filter from date (YYYY-MM-DD)"
// @Param        date_to     query     string                  false  "Filter to date (YYYY-MM-DD)"
// @Param        sort_by     query     string                  false  "Sort field: amount or expense_date"
// @Param        sort_order  query     string                  false  "Sort direction: asc or desc (default: desc)"
// @Param        limit       query     int                     false  "Max number of results"
// @Success      200         {object}  map[string]interface{}  "Expenses retrieved"
// @Failure      400         {object}  map[string]interface{}  "Validation error"
// @Failure      401         {object}  map[string]interface{}  "Unauthorized"
// @Router       /api/v1/expenses [get]
func (c *ExpenseController) Get() {
	userID, ok := userIDFromHeader(&c.Controller)
	if !ok {
		utils.Error(c.Ctx, 401, "Unauthorized")
		return
	}

	if idParam := c.Ctx.Input.Param(":id"); idParam != "" {
		c.getOne(userID, idParam)
		return
	}

	dateFrom := c.GetString("date_from")
	dateTo := c.GetString("date_to")
	sortBy := c.GetString("sort_by")
	sortOrder := c.GetString("sort_order")

	if dateFrom != "" && !utils.ValidateDate(dateFrom) {
		utils.Error(c.Ctx, 400, "date_from must be in YYYY-MM-DD format")
		return
	}
	if dateTo != "" && !utils.ValidateDate(dateTo) {
		utils.Error(c.Ctx, 400, "date_to must be in YYYY-MM-DD format")
		return
	}
	if sortBy != "" && sortBy != "amount" && sortBy != "expense_date" {
		utils.Error(c.Ctx, 400, "sort_by must be amount or expense_date")
		return
	}
	if sortOrder != "" && sortOrder != "asc" && sortOrder != "desc" {
		utils.Error(c.Ctx, 400, "sort_order must be asc or desc")
		return
	}

	expenses, err := services.ListExpenses(userID, services.ListExpensesInput{
		DateFrom:  dateFrom,
		DateTo:    dateTo,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Limit:     utils.ParseLimit(c.GetString("limit"), 0),
	})
	if err != nil {
		logs.Error("[Expense] List failed: user_id=%d err=%v", userID, err)
		utils.Error(c.Ctx, 500, "Internal server error")
		return
	}

	logs.Info("[Expense] List: user_id=%d count=%d", userID, len(expenses))
	utils.Success(c.Ctx, 200, "Expenses retrieved", expenses)
}

// getOne handles GET /api/v1/expenses/:id
func (c *ExpenseController) getOne(userID int, idParam string) {
	id, ok := utils.ParseID(idParam)
	if !ok {
		utils.Error(c.Ctx, 400, "Invalid expense ID")
		return
	}

	expense, err := services.GetExpense(id, userID)
	if errors.Is(err, services.ErrExpenseNotFound) {
		utils.Error(c.Ctx, 404, "Expense not found")
		return
	}
	if err != nil {
		logs.Error("[Expense] Get failed: id=%d user_id=%d err=%v", id, userID, err)
		utils.Error(c.Ctx, 500, "Internal server error")
		return
	}

	utils.Success(c.Ctx, 200, "Expense retrieved successfully", expense)
}

// Put updates an existing expense owned by the authenticated user.
// @Summary      Update expense
// @Description  Updates an expense owned by the authenticated user
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header    int                     true  "User ID"
// @Param        id         path      int                     true  "Expense ID"
// @Param        body       body      expenseInput            true  "Expense payload"
// @Success      200        {object}  map[string]interface{}  "Expense updated successfully"
// @Failure      400        {object}  map[string]interface{}  "Validation error"
// @Failure      401        {object}  map[string]interface{}  "Unauthorized"
// @Failure      404        {object}  map[string]interface{}  "Expense not found"
// @Router       /api/v1/expenses/{id} [put]
func (c *ExpenseController) Put() {
	userID, ok := userIDFromHeader(&c.Controller)
	if !ok {
		utils.Error(c.Ctx, 401, "Unauthorized")
		return
	}

	id, ok := utils.ParseID(c.Ctx.Input.Param(":id"))
	if !ok {
		utils.Error(c.Ctx, 400, "Invalid expense ID")
		return
	}

	body, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		utils.Error(c.Ctx, 400, "Failed to read request body")
		return
	}

	var input expenseUpdateInput
	if err := json.Unmarshal(body, &input); err != nil {
		utils.Error(c.Ctx, 400, "Invalid request body")
		return
	}

	expense, err := services.UpdateExpense(id, userID, services.UpdateExpenseInput{
		Title:       input.Title,
		Amount:      input.Amount,
		Category:    input.Category,
		Note:        input.Note,
		ExpenseDate: input.ExpenseDate,
	})
	if errors.Is(err, services.ErrExpenseNotFound) {
		utils.Error(c.Ctx, 404, "Expense not found")
		return
	}
	if err != nil {
		handleServiceError(&c.Controller, "[Expense] Put", err)
		return
	}

	logs.Info("[Expense] Updated: id=%d user_id=%d", id, userID)
	utils.Success(c.Ctx, 200, "Expense updated successfully", expense)
}

// Delete removes an expense owned by the authenticated user.
// @Summary      Delete expense
// @Description  Deletes an expense owned by the authenticated user
// @Tags         expenses
// @Produce      json
// @Param        X-User-ID  header    int                     true  "User ID"
// @Param        id         path      int                     true  "Expense ID"
// @Success      200        {object}  map[string]interface{}  "Expense deleted successfully"
// @Failure      400        {object}  map[string]interface{}  "Invalid ID"
// @Failure      401        {object}  map[string]interface{}  "Unauthorized"
// @Failure      404        {object}  map[string]interface{}  "Expense not found"
// @Router       /api/v1/expenses/{id} [delete]
func (c *ExpenseController) Delete() {
	userID, ok := userIDFromHeader(&c.Controller)
	if !ok {
		utils.Error(c.Ctx, 401, "Unauthorized")
		return
	}

	id, ok := utils.ParseID(c.Ctx.Input.Param(":id"))
	if !ok {
		utils.Error(c.Ctx, 400, "Invalid expense ID")
		return
	}

	err := services.DeleteExpense(id, userID)
	if errors.Is(err, services.ErrExpenseNotFound) {
		utils.Error(c.Ctx, 404, "Expense not found")
		return
	}
	if err != nil {
		logs.Error("[Expense] Delete failed: id=%d user_id=%d err=%v", id, userID, err)
		utils.Error(c.Ctx, 500, "Internal server error")
		return
	}

	logs.Info("[Expense] Deleted: id=%d user_id=%d", id, userID)
	utils.Success(c.Ctx, 200, "Expense deleted successfully", nil)
}
