package handler

import (
	"strconv"

	payService "github.com/ghulammuzz/backend-parkerin/internal/payment/svc"
	"github.com/ghulammuzz/backend-parkerin/pkg/log"
	"github.com/ghulammuzz/backend-parkerin/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	payService payService.PaymentService
}

func (h PaymentHandler) Router(r fiber.Router) {
	// r.Post("/pay/:packageID", middleware.JWTProtected(), h.CreateTransaction)
	r.Post("/pay/:packageID", h.CreateTransaction)
}

func (h PaymentHandler) CreateTransaction(c *fiber.Ctx) error {
	// userToken := c.Locals("user").(*jwt.Token)

	// claims, ok := userToken.Claims.(jwt.MapClaims)
	// if !ok || !userToken.Valid {
	// 	return response.JSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	// }

	// userID := int(claims["user_id"].(float64))

	userID := 33

	packageID, err := strconv.Atoi(c.Params("packageID"))
	if err != nil {
		return response.JSON(c, 400, "invalid package ID", nil)
	}

	if packageID != 1 && packageID != 2 {
		return response.JSON(c, 400, "product not valid", nil)
	}

	log.Debug("userId = %d, packageId = %d", userID, packageID)

	transaction, err := h.payService.CreateTransaction(userID, packageID)
	if err != nil {
		return response.JSON(c, 500, "error creating transaction", err.Error())
	}
	return response.JSON(c, 200, "success creating transaction", transaction)

	// return response.JSON(c, 200, "success", userID)
}

func NewPaymentHandler(payService payService.PaymentService) *PaymentHandler {
	return &PaymentHandler{payService: payService}
}
