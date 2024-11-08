package handler

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/ghulammuzz/backend-parkerin/internal/middleware"
	"github.com/ghulammuzz/backend-parkerin/internal/store/entity"
	"github.com/ghulammuzz/backend-parkerin/internal/store/svc"
	"github.com/ghulammuzz/backend-parkerin/pkg/log"
	"github.com/ghulammuzz/backend-parkerin/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type StoreHandler struct {
	storeService svc.StoreService
}

func NewStoreHandler(storeService svc.StoreService) *StoreHandler {
	return &StoreHandler{storeService: storeService}
}

func (h *StoreHandler) Router(r fiber.Router) {
	r.Get("/stores", h.ListStores)
	r.Get("/store/:id", h.GetStoreDetail)
	r.Get("/store-dashboard", middleware.JWTProtected(), h.DashboardStore)
	r.Put("/store-hiring", middleware.JWTProtected(), h.UpdateIsHiringHandler)
	r.Post("/store-img", middleware.JWTProtected(), h.UplaodStoreIMGHandler)
}

func (h *StoreHandler) ListStores(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10 // default limit
	}

	isHiring, err := strconv.ParseBool(c.Query("isHiring", "false"))
	if err != nil {
		return response.JSON(c, fiber.StatusBadRequest, "Invalid isHiring parameter", nil)
	}
	// log.Debug(strconv.FormatBool(isHiring))

	stores, err := h.storeService.ListStores(page, limit, isHiring)
	if err != nil {
		return response.JSON(c, fiber.StatusInternalServerError, "Failed to retrieve store list", err.Error())
	}

	return response.JSON(c, fiber.StatusOK, "Store list retrieved successfully", stores)
}

func (h *StoreHandler) GetStoreDetail(c *fiber.Ctx) error {
	storeIDStr := c.Params("id")
	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil || storeID < 1 {
		return response.JSON(c, fiber.StatusBadRequest, "Invalid store ID", nil)
	}

	storeDetail, err := h.storeService.GetStoreDetail(storeID)
	if err != nil {
		return response.JSON(c, fiber.StatusInternalServerError, "Failed to retrieve store details", err.Error())
	}

	return response.JSON(c, fiber.StatusOK, "Store details retrieved successfully", storeDetail)
}

func (h *StoreHandler) DashboardStore(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return response.JSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	store, err := h.storeService.DashboardStore(userID)
	if err != nil {
		return response.JSON(c, fiber.StatusInternalServerError, "Failed to fetch store details", err.Error())
	}

	return response.JSON(c, fiber.StatusOK, "Store details", store)
}

func (h *StoreHandler) UpdateIsHiringHandler(c *fiber.Ctx) error {
	var req entity.UpdateIsHiringRequest
	if err := c.BodyParser(&req); err != nil {
		return response.JSON(c, 400, "Payload error", err.Error())

	}
	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return response.JSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	// using get storeID by userID
	storeID, err := h.storeService.GetStoreIDByUserID(userID)
	if err != nil {
		return response.JSON(c, 400, "invalid user ID", nil)
	}
	log.Debug(fmt.Sprint(storeID))

	if err := h.storeService.UpdateIsHiring(req.IsHiring, storeID); err != nil {
		return response.JSON(c, 500, "error svc", err.Error())

	}

	return response.JSON(c, 200, "Success Updated", nil)

}
func (h *StoreHandler) UplaodStoreIMGHandler(c *fiber.Ctx) error {

	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return response.JSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	storeID, err := h.storeService.GetStoreIDByUserID(userID)
	if err != nil {
		return response.JSON(c, 400, "invalid user ID", nil)
	}
	log.Debug(fmt.Sprint(storeID))

	img, err := c.FormFile("img")
	if err != nil {
		log.Error("Error getting img", slog.String("error", err.Error()))
		return response.JSON(c, 400, "img file is required", err.Error())
	}

	const maxSize = 2 * 1024 * 1024
	if img.Size > maxSize {
		return response.JSON(c, 400, "Image size exceeds 2 MB", nil)
	}

	err = h.storeService.UploadStoreIMG(storeID, img)
	if err != nil {
		return response.JSON(c, 500, "error svc upload img", err.Error())
	}

	return response.JSON(c, 200, "img uploaded", nil)
}
