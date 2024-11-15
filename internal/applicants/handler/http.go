package handler

import (
	"strconv"

	"github.com/dgrijalva/jwt-go"

	appService "github.com/ghulammuzz/backend-parkerin/internal/applicants/svc"
	"github.com/ghulammuzz/backend-parkerin/internal/middleware"

	storeService "github.com/ghulammuzz/backend-parkerin/internal/store/svc"
	"github.com/ghulammuzz/backend-parkerin/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type ApplicationHandler struct {
	appService   appService.ApplicationService
	storeService storeService.StoreService
}

func NewApplicationHandler(appService appService.ApplicationService, storeService storeService.StoreService) *ApplicationHandler {
	return &ApplicationHandler{appService, storeService}
}

func (h *ApplicationHandler) Router(r fiber.Router) {
	r.Post("/apply-store/:storeID", middleware.JWTProtected(), h.ApplyStore)
	r.Post("/apply-user/:userID", middleware.JWTProtected(), h.ApplyUser)
	r.Get("/application/store", middleware.JWTProtected(), h.ReviewApplicationsStore)
	r.Get("/application/user", middleware.JWTProtected(), h.ReviewApplicationsUser)
	// need test
	r.Put("/status-apply-user/:appID", middleware.JWTProtected(), h.UpdateApplicationUserStatus)
	r.Put("/status-apply-store/:appID", middleware.JWTProtected(), h.UpdateApplicationStoreStatus)
	r.Delete("/applicants/:appID", middleware.JWTProtected(), h.DeleteAppsInUserHandler)
}

// us (using jwt)
func (h *ApplicationHandler) ApplyStore(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return response.JSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))
	roles := string(claims["role"].(string))

	// log.Debug(fmt.Sprint("roles :", roles))
	if roles != "tukang" {
		return response.JSON(c, 401, "Unauthorized", nil)
	}
	// log.Debug(fmt.Sprint("user id :", userID))

	storeID, err := strconv.Atoi(c.Params("storeID"))
	if err != nil {
		return response.JSON(c, 400, "invalid store ID", nil)
	}
	// log.Debug(fmt.Sprint("store id :", storeID))

	err = h.appService.CreateApply(userID, storeID, false)
	if err != nil {
		return response.JSON(c, 500, "error apply", err.Error())
	}

	return response.JSON(c, 200, "Application submitted successfully", nil)
}

// st (using jwt)
func (h *ApplicationHandler) ApplyUser(c *fiber.Ctx) error {

	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return response.JSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userSelfID := int(claims["user_id"].(float64))

	// using get storeID by userID
	storeID, err := h.storeService.GetStoreIDByUserID(userSelfID)
	if err != nil {
		return response.JSON(c, 400, "invalid user ID", nil)
	}
	// log.Debug(fmt.Sprint(storeID))

	userID, err := strconv.Atoi(c.Params("userID"))
	if err != nil {
		return response.JSON(c, 400, "invalid user ID", nil)
	}

	if err := h.appService.CreateApply(userID, storeID, true); err != nil {
		return response.JSON(c, 500, "error svc createapply", err.Error())
	}

	return response.JSON(c, 200, "Application submitted successfully", nil)
}

// st (using jwt)
func (h *ApplicationHandler) ReviewApplicationsStore(c *fiber.Ctx) error {

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
	// log.Debug(fmt.Sprint(storeID))

	applications, err := h.appService.ReviewApplications(storeID)
	if err != nil {
		return response.JSON(c, 500, "review app svc", err.Error())
	}

	return response.JSON(c, 200, "Applications retrieved successfully", applications)
}

// us (using jwt)
// review app by user

func (h *ApplicationHandler) ReviewApplicationsUser(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return response.JSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	isDirectHire, err := strconv.ParseBool(c.Query("is_direct_hire", "false"))
	if err != nil {
		return response.JSON(c, fiber.StatusBadRequest, "Invalid isdirecthire parameter", nil)
	}

	applications, err := h.appService.ReviewApplicationsUser(userID, isDirectHire)
	if err != nil {
		return response.JSON(c, fiber.StatusInternalServerError, "Error reviewing applications", err.Error())
	}

	return response.JSON(c, fiber.StatusOK, "Applications retrieved successfully", applications)
}

// us
func (h *ApplicationHandler) UpdateApplicationUserStatus(c *fiber.Ctx) error {

	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return response.JSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	appID, err := strconv.Atoi(c.Params("appID"))
	if err != nil {
		return response.JSON(c, 400, "invalid app id", nil)
	}

	status := c.Query("update")
	if status != "accepted" && status != "rejected" {
		return response.JSON(c, 400, "invalid status; must be 'accepted' or 'rejected'", nil)
	}

	if status == "accepted" {
		if err := h.appService.AcceptApplicationUser(appID, userID); err != nil {
			return response.JSON(c, 500, "acc user svc", err.Error())
		}
		return response.JSON(c, 200, "Application accepted", nil)
	} else {
		if err := h.appService.RejectApplicationUser(appID, userID); err != nil {
			return response.JSON(c, 500, "reject user svc", err.Error())
		}
		return response.JSON(c, 200, "Application rejected", nil)
	}
}

// st
func (h *ApplicationHandler) UpdateApplicationStoreStatus(c *fiber.Ctx) error {

	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return response.JSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	appID, err := strconv.Atoi(c.Params("appID"))
	if err != nil {
		return response.JSON(c, 400, "invalid app id", nil)
	}

	storeID, err := h.storeService.GetStoreIDByUserID(userID)
	// log.Debug(fmt.Sprint(storeID))
	if err != nil {
		return response.JSON(c, 400, "invalid user ID", nil)
	}

	status := c.Query("update")
	if status != "accepted" && status != "rejected" {
		return response.JSON(c, 400, "invalid status; must be 'accepted' or 'rejected'", nil)
	}
	// log.Debug(fmt.Sprint(status))
	// log.Debug(fmt.Sprint(appID))

	if status == "accepted" {
		if err := h.appService.AcceptApplicationStore(appID, storeID); err != nil {
			return response.JSON(c, 500, "acc store svc", err.Error())
		}
		return response.JSON(c, 200, "Application accepted", nil)
	} else {
		if err := h.appService.RejectApplicationStore(appID, storeID); err != nil {
			return response.JSON(c, 500, "reject store svc", err.Error())
		}
		return response.JSON(c, 200, "Application rejected", nil)
	}
}

func (h *ApplicationHandler) DeleteAppsInUserHandler(c *fiber.Ctx) error {

	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return response.JSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	appID, err := strconv.Atoi(c.Params("appID"))
	if err != nil {
		return response.JSON(c, 400, "invalid app id", nil)
	}

	err = h.appService.DeleteAppsInUser(userID, appID)
	if err != nil {
		return response.JSON(c, 500, "error svc delete apps", nil)
	}

	return response.JSON(c, 200, "success delete apps ", nil)
}
