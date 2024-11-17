package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	"github.com/kooroshh/fiber-boostrap/pkg/jwt_token"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"go.elastic.co/apm"
	"golang.org/x/crypto/bcrypt"
)

// Register handles the HTTP request to register a new user.
// It parses the request body to create a new user object, validates the user data,
// hashes the user's password, and inserts the new user into the database.
// It responds with a success message and the user data (excluding the password)
// if the registration is successful, or with an error message if any step fails.
func Register(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Register", "controller")
	defer span.End()

	user := new(models.User)

	err := ctx.BodyParser(user)
	if err != nil {
		errResponse := fmt.Errorf("failed to parse body request: %v", err)
		log.Println("Failed to parse body request: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	err = user.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate body request: %v", err)
		log.Println("Failed to validate body request: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		errResponse := fmt.Errorf("failed to hash password: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}
	user.Password = string(hashPassword)

	err = repository.InsertNewUser(spanCtx, user)
	if err != nil {
		errResponse := fmt.Errorf("failed to insert new user: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	resp := user
	resp.Password = ""

	return response.SendSuccessResponse(ctx, resp)
}

// Login handles user authentication by validating credentials provided in the HTTP request.
// It parses the login request, validates the user credentials, and retrieves the user from the database.
// If the credentials are correct, it generates a JWT token and a refresh token.
// It then creates a new user session in the database with these tokens.
// Finally, it returns a success response with the generated tokens or a failure response in case of any errors.
func Login(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Login", "controller")
	defer span.End()

	var (
		loginReq = new(models.LoginRequest)
		resp     = models.LoginResponse{}
		now      = time.Now()
	)

	err := ctx.BodyParser(loginReq)
	if err != nil {
		errResponse := fmt.Errorf("failed to parse body request: %v", err)
		log.Println("Failed to parse body request: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	err = loginReq.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate body request: %v", err)
		log.Println("Failed to validate body request: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	user, err := repository.GetUserByUsername(spanCtx, loginReq.Username)
	if err != nil {
		errResponse := fmt.Errorf("failed to get user by username: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "username/password is wrong", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		errResponse := fmt.Errorf("failed to compare hash password: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "username/password is wrong", nil)
	}

	token, err := jwt_token.GenerateToken(spanCtx, user.Username, user.FullName, "token", now)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	refreshToken, err := jwt_token.GenerateToken(spanCtx, user.Username, user.FullName, "refresh_token", now)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate refresh token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	userSession := &models.UserSession{
		UserID:              user.ID,
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        now.Add(jwt_token.MapTypeToken["token"]),
		RefreshTokenExpired: now.Add(jwt_token.MapTypeToken["refresh_token"]),
	}
	err = repository.InsertNewUserSession(spanCtx, userSession)
	if err != nil {
		errResponse := fmt.Errorf("failed insert user session: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	resp.Username = user.Username
	resp.FullName = user.FullName
	resp.Token = token
	resp.RefreshToken = refreshToken

	return response.SendSuccessResponse(ctx, resp)
}

// Logout handles the HTTP request to log out a user by deleting their session.
// It retrieves the authorization token from the request header and attempts to
// delete the corresponding user session from the database. If the deletion is
// successful, it returns a success response. Otherwise, it returns an error response.
func Logout(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Logout", "controller")
	defer span.End()

	token := ctx.Get("Authorization")
	err := repository.DeleteUserSessionByToken(spanCtx, token)
	if err != nil {
		errResponse := fmt.Errorf("failed delete user session: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}
	return response.SendSuccessResponse(ctx, nil)
}

// RefreshToken handles the HTTP request to refresh a user's token.
// It retrieves the refresh token from the request header and attempts to
// update the corresponding user session in the database with a new token.
// If the update is successful, it returns a success response with the new token.
// Otherwise, it returns an error response.
func RefreshToken(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "RefreshToken", "controller")
	defer span.End()

	now := time.Now()
	refreshToken := ctx.Get("Authorization")
	username := ctx.Locals("username").(string)
	fullName := ctx.Locals("full_name").(string)

	token, err := jwt_token.GenerateToken(spanCtx, username, fullName, "token", now)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	err = repository.UpdateUserSessionToken(spanCtx, token, now.Add(jwt_token.MapTypeToken["token"]), refreshToken)
	if err != nil {
		errResponse := fmt.Errorf("failed to update token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	return response.SendSuccessResponse(ctx, fiber.Map{
		"token": token,
	})
}
