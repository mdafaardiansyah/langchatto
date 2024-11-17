package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"go.elastic.co/apm"
)

// GetHistory handles the HTTP request to retrieve the history of messages.
// It initiates a trace span for monitoring, retrieves all messages from the repository,
// and sends a success response with the messages or a failure response in case of an error.
func GetHistory(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "GetHistory", "controller")
	defer span.End()

	resp, err := repository.GetAllMessage(spanCtx)
	if err != nil {
		log.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}
	return response.SendSuccessResponse(ctx, resp)
}
