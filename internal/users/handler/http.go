package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/ghulammuzz/backend-parkerin/internal/middleware"
	"github.com/ghulammuzz/backend-parkerin/internal/users/entity"
	"github.com/ghulammuzz/backend-parkerin/internal/users/svc"
	"github.com/ghulammuzz/backend-parkerin/pkg/form"
	"github.com/ghulammuzz/backend-parkerin/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService svc.UserService
	val         *validator.Validate
}

func NewUserHandler(userService svc.UserService, val *validator.Validate) *UserHandler {
	return &UserHandler{userService, val}
}

func (h *UserHandler) Router(r fiber.Router) {
	r.Post("/user/register", h.RegisterUser)
	r.Post("/user/login", h.LoginUser)
	r.Post("/store/login", h.LoginStore)
	r.Get("/user/dashboard", middleware.JWTProtected(), h.DashboardUser)
	r.Get("/store/dashboard", middleware.JWTProtected(), h.DashboardStore)
}

func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	user := new(entity.UserRegisterRequest)
	if err := c.BodyParser(user); err != nil {
		return response.JSON(c, 400, "invalid payload", err.Error())
	}
	if err := h.val.Struct(user); err != nil {
		validationErrors := form.ValidationErrorResponse(err)
		return response.JSON(c, 400, "Validation failed", validationErrors)
	}

	exists, err := h.userService.IsPhoneNumberExists(user.PhoneNumber)
	if err != nil {
		return response.JSON(c, 500, "Error checking phone number", err.Error())
	}
	if exists {
		return response.JSON(c, 400, "Phone number already registered", nil)
	}
	if user.Role == "store" {
		if user.StoreName == nil || user.Address == nil || user.Latitude == nil || user.Longitude == nil {
			return response.JSON(c, 400, "Validation failed", "store_name, address, latitude, and longitude are required for store role")
		}
	}

	if err := h.userService.RegisterUser(user); err != nil {
		return response.JSON(c, 500, "register svc error", err.Error())
	}

	return response.JSON(c, 201, "User registered successfully", nil)
}

func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	loginRequest := new(entity.UserLoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return response.JSON(c, 400, "Invalid payload", err.Error())
	}

	token, err := h.userService.LoginUser(loginRequest)
	if err != nil {
		return response.JSON(c, 401, "Login failed", err.Error())
	}

	return response.JSON(c, 200, "Login successful", token)
}

func (h *UserHandler) LoginStore(c *fiber.Ctx) error {
	loginRequest := new(entity.UserLoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return response.JSON(c, 400, "Invalid payload", err.Error())
	}

	token, err := h.userService.LoginStore(loginRequest)
	if err != nil {
		return response.JSON(c, 401, "Login failed", err.Error())
	}

	return response.JSON(c, 200, "Login successful", token)
}

func (h *UserHandler) DashboardUser(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return response.JSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	user, err := h.userService.GetUserDetails(userID)
	if err != nil {
		return response.JSON(c, fiber.StatusInternalServerError, "Failed to fetch user details", err.Error())
	}

	return response.JSON(c, fiber.StatusOK, "User details", fiber.Map{
		"id":           user.ID,
		"name":         user.Name,
		"phone_number": user.PhoneNumber,
		"role":         user.Role,
	})
}

func (h *UserHandler) DashboardStore(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return response.JSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	store, err := h.userService.GetStoreDetails(userID)
	if err != nil {
		return response.JSON(c, fiber.StatusInternalServerError, "Failed to fetch store details", err.Error())
	}

	return response.JSON(c, fiber.StatusOK, "Store details", fiber.Map{
		"store_id":     store.ID,
		"name":         store.Name,
		"phone_number": store.PhoneNumber,
		"role":         store.Role,
		"store_name":   store.StoreName,
		"address":      store.Address,
		"longitude":    store.Longitude,
		"latitude":     store.Latitude,
	})
}
