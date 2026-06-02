package utils

import (
	"github.com/beego/beego/v2/server/web/context"
)

// Success writes a standardised success JSON response.
// Pass nil for data when there is no payload to return.
func Success(ctx *context.Context, status int, message string, data interface{}) {
	ctx.Output.SetStatus(status)

	body := map[string]interface{}{
		"success": true,
		"code":    status,
		"message": message,
	}
	if data != nil {
		body["data"] = data
	}

	ctx.Output.JSON(body, false, false)
}

// Error writes a standardised error JSON response.
// Never pass internal error details — only safe client-facing messages.
func Error(ctx *context.Context, status int, message string) {
	ctx.Output.SetStatus(status)
	ctx.Output.JSON(map[string]interface{}{
		"success": false,
		"code":    status,
		"message": message,
	}, false, false)
}
