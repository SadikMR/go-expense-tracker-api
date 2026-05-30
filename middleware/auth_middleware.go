package middleware

import (
	"strconv"

	"github.com/SadikMR/go-expense-tracker-api/models"
	"github.com/SadikMR/go-expense-tracker-api/utils"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
)

// AuthMiddleware validates the X-User-ID header on all protected routes.
// Rejects the request with 401 if the header is missing, malformed,
// or refers to a non-existent user.
func AuthMiddleware(ctx *context.Context) {
	header := ctx.Input.Header("X-User-ID")
	if header == "" {
		logs.Warn("[Middleware] Missing X-User-ID header: path=%s", ctx.Request.URL.Path)
		utils.Error(ctx, 401, "Unauthorized")
		ctx.ResponseWriter.Started = true
		return
	}

	id, err := strconv.Atoi(header)
	if err != nil || id <= 0 {
		logs.Warn("[Middleware] Invalid X-User-ID value: value=%s path=%s", header, ctx.Request.URL.Path)
		utils.Error(ctx, 401, "Unauthorized")
		ctx.ResponseWriter.Started = true
		return
	}

	user, err := models.GetUserByID(id)
	if err != nil {
		logs.Error("[Middleware] Failed to look up user: user_id=%d err=%v", id, err)
		utils.Error(ctx, 401, "Unauthorized")
		ctx.ResponseWriter.Started = true
		return
	}
	if user == nil {
		logs.Warn("[Middleware] Non-existent user_id: user_id=%d path=%s", id, ctx.Request.URL.Path)
		utils.Error(ctx, 401, "Unauthorized")
		ctx.ResponseWriter.Started = true
		return
	}
}
