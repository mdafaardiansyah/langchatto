package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"log"
)

func GetHistory(ctx *fiber.Ctx) error {
	resp, err := repository.GetAllMessage(ctx.Context())
	if err != nil {
		log.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}
	return response.SendSuccessResponse(ctx, resp)
}
