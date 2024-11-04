package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSON(c *fiber.Ctx, code int, message string, data interface{}) error {
	json := Response{
		Message: message,
		Data:    data,
	}
	return c.Status(code).JSON(json)
}
