package router

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	"github.com/kooroshh/fiber-boostrap/pkg/jwt_token"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"go.elastic.co/apm"
)

// MiddlewareValidateAuth is a middleware that validates the authorization header
// on each request. If the header is empty, it returns a 401 Unauthorized response.
// If the header is not empty, it attempts to retrieve the corresponding user session
// from the database. If the retrieval is successful, it validates the JWT token
// using the ValidateToken function. If the validation is successful, it sets the
// username and full_name locals on the request context and calls the next handler.
// If the validation fails, it returns a 401 Unauthorized response.
func MiddlewareValidateAuth(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "MiddlewareValidateAuth", "middleware")
	defer span.End()

	auth := ctx.Get("authorization")
	if auth == "" {
		log.Println("authorization empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	_, err := repository.GetUserSessionByToken(spanCtx, auth)
	if err != nil {
		log.Println("failed to get user session on DB: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	claim, err := jwt_token.ValidateToken(spanCtx, auth)
	if err != nil {
		log.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("jwt token is expired: ", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}

// MiddlewareRefreshToken is a middleware that validates the authorization header
// on each request. If the header is empty, it returns a 401 Unauthorized response.
// If the header is not empty, it attempts to retrieve the corresponding user session
// from the database. If the retrieval is successful, it validates the JWT token
// using the ValidateToken function. If the validation is successful, it sets the
// username and full_name locals on the request context and calls the next handler.
// If the validation fails, it returns a 401 Unauthorized response.
func MiddlewareRefreshToken(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "MiddlewareRefreshToken", "middleware")
	defer span.End()

	auth := ctx.Get("authorization")
	if auth == "" {
		log.Println("authorization empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	claim, err := jwt_token.ValidateToken(spanCtx, auth)
	if err != nil {
		log.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("jwt token is expired: ", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}
