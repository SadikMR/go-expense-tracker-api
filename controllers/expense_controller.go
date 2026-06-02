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
// @Description  Creates a new expense for the authenticated user. Category must be one of the allowed values and expense_date must use YYYY-MM-DD.
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        X-User-ID  header    int                    true  "Authenticated user ID (returned after login). Example: 123"
// @Param        body       body      ExpenseCreateRequest  true  "Expense creation payload"
// @Success      201        {object}  ExpenseResponseWrapper  "Expense created successfully"
// @Failure      400        {object}  ErrorResponse            "Validation error"
// @Failure      401        {object}  ErrorResponse            "Unauthorized"
// @Failure      500        {object}  ErrorResponse            "Internal server error"
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

	var input ExpenseCreateRequest
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

// Get returns a list of expenses for the authenticated user.
// @Summary      List expenses
// @Description  Returns expenses for the authenticated user with optional date, category, sort, and limit filters.
// @Tags         expenses
// @Produce      json
// @Security     ApiKeyAuth
// @Param        X-User-ID   header    int     true   "Authenticated user ID (returned after login). Example: 123"
// @Param        date_from   query     string  false  "Filter from date in YYYY-MM-DD format"
// @Param        date_to     query     string  false  "Filter to date in YYYY-MM-DD format"
// @Param        category    query     string  false  "Filter by category. Allowed values: Food, Transport, Housing, Entertainment, Shopping, Healthcare, Education, Utilities, Other"
// @Param        sort_by     query     string  false  "Sort field: amount or expense_date"
// @Param        sort_order  query     string  false  "Sort direction: asc or desc (default: desc)"
// @Param        limit       query     int     false  "Maximum number of results returned"
// @Success      200         {object}  ExpenseListResponse  "Expenses retrieved successfully"
// @Failure      400         {object}  ErrorResponse        "Validation error"
// @Failure      401         {object}  ErrorResponse        "Unauthorized"
// @Failure      500         {object}  ErrorResponse        "Internal server error"
// @Router       /api/v1/expenses [get]
func (c *ExpenseController) Get() {
	userID, ok := userIDFromHeader(&c.Controller)
	if !ok {
		utils.Error(c.Ctx, 401, "Unauthorized")
		return
	}

	category := c.GetString("category")
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
	if category != "" && !utils.ValidateCategory(category) {
		utils.Error(c.Ctx, 400, "Invalid category")
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
		Category:  category,
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

// GetOne returns a single expense owned by the authenticated user.
// @Summary      Get expense
// @Description  Returns a single expense by ID. Ownership is enforced by X-User-ID.
// @Tags         expenses
// @Produce      json
// @Security     ApiKeyAuth
// @Param        X-User-ID  header    int                    true  "Authenticated user ID (returned after login). Example: 123"
// @Param        id         path      int                    true  "Expense ID"
// @Success      200        {object}  ExpenseResponseWrapper  "Expense retrieved successfully"
// @Failure      400        {object}  ErrorResponse            "Invalid expense ID"
// @Failure      401        {object}  ErrorResponse            "Unauthorized"
// @Failure      404        {object}  ErrorResponse            "Expense not found"
// @Failure      500        {object}  ErrorResponse            "Internal server error"
// @Router       /api/v1/expenses/{id} [get]
func (c *ExpenseController) GetOne() {
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
// @Description  Updates an expense owned by the authenticated user. Only fields present in the request are changed.
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        X-User-ID  header    int                     true  "Authenticated user ID (returned after login). Example: 123"
// @Param        id         path      int                     true  "Expense ID"
// @Param        body       body      ExpenseUpdateRequest   true  "Expense update payload"
// @Success      200        {object}  ExpenseResponseWrapper  "Expense updated successfully"
// @Failure      400        {object}  ErrorResponse            "Validation error"
// @Failure      401        {object}  ErrorResponse            "Unauthorized"
// @Failure      404        {object}  ErrorResponse            "Expense not found"
// @Failure      500        {object}  ErrorResponse            "Internal server error"
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

	var input ExpenseUpdateRequest
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
// @Description  Deletes an expense owned by the authenticated user.
// @Tags         expenses
// @Produce      json
// @Security     ApiKeyAuth
// @Param        X-User-ID  header    int                     true  "Authenticated user ID (returned after login). Example: 123"
// @Param        id         path      int                     true  "Expense ID"
// @Success      200        {object}  StandardResponse  "Expense deleted successfully"
// @Failure      400        {object}  ErrorResponse     "Invalid expense ID"
// @Failure      401        {object}  ErrorResponse     "Unauthorized"
// @Failure      404        {object}  ErrorResponse     "Expense not found"
// @Failure      500        {object}  ErrorResponse     "Internal server error"
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
