package controllers

import (
	"encoding/json"
	"errors"

	"github.com/SadikMR/go-expense-tracker-api/services"
	"github.com/SadikMR/go-expense-tracker-api/utils"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// AuthController handles user registration and login.
type AuthController struct {
	beego.Controller
}

type registerInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register creates a new user account.
// @Summary      Register
// @Description  Creates a new user account with name, email, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      object  true  "Registration payload"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}
// @Router       /api/v1/auth/register [post]
func (c *AuthController) Register() {
	var input registerInput
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		utils.Error(c.Ctx, 400, "Invalid request body")
		return
	}
	if !utils.ValidateRequired(input.Name, input.Email, input.Password) {
		utils.Error(c.Ctx, 400, "Name, email, and password are required")
		return
	}
	if !utils.ValidateEmail(input.Email) {
		utils.Error(c.Ctx, 400, "Invalid email format")
		return
	}
	if !utils.ValidateMinLength(input.Password, 6) {
		utils.Error(c.Ctx, 400, "Password must be at least 6 characters")
		return
	}

	err := services.Register(input.Name, input.Email, input.Password)
	if errors.Is(err, services.ErrEmailExists) {
		logs.Warn("[Auth] Register attempt with existing email: %s", input.Email)
		utils.Error(c.Ctx, 409, "Email already exists")
		return
	}
	if err != nil {
		logs.Error("[Auth] Register failed: %v", err)
		utils.Error(c.Ctx, 500, "Internal server error")
		return
	}

	// NOTE: never log the password — log only safe fields
	logs.Info("[Auth] New user registered: name=%s email=%s", input.Name, input.Email)
	utils.Success(c.Ctx, 201, "User registered successfully", nil)
}

// Login authenticates a user and returns their ID.
// @Summary      Login
// @Description  Authenticates a user and returns user_id for use in X-User-ID header
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      object  true  "Login payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /api/v1/auth/login [post]
func (c *AuthController) Login() {
	var input loginInput
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		utils.Error(c.Ctx, 400, "Invalid request body")
		return
	}
	if !utils.ValidateRequired(input.Email, input.Password) {
		utils.Error(c.Ctx, 400, "Email and password are required")
		return
	}

	user, err := services.Login(input.Email, input.Password)
	if errors.Is(err, services.ErrInvalidCreds) {
		// Warn because repeated failures could indicate a brute-force attempt
		logs.Warn("[Auth] Failed login attempt: email=%s", input.Email)
		utils.Error(c.Ctx, 401, "Invalid email or password")
		return
	}
	if err != nil {
		logs.Error("[Auth] Login error: %v", err)
		utils.Error(c.Ctx, 500, "Internal server error")
		return
	}

	// NOTE: never log the password — log only safe fields
	logs.Info("[Auth] Login successful: user_id=%d email=%s", user.ID, user.Email)
	utils.Success(c.Ctx, 200, "Login successful", map[string]interface{}{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	})
}
