package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"agungmohmd/intikm-test-api/helper"
	"agungmohmd/intikm-test-api/usecase"

	ut "github.com/go-playground/universal-translator"
	"github.com/gofiber/fiber/v2"
	jwtFiber "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gopkg.in/go-playground/validator.v9"
)

type Handler struct {
	FiberApp   *fiber.App
	ContractUC *usecase.ContractUC
	Db         *sql.DB
	Validator  *validator.Validate
	Translator ut.Translator
	JwtConfig  jwtFiber.Config
}

type CustomJWT struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func userContext(ctx context.Context, subject, id interface{}) context.Context {
	return context.WithValue(ctx, subject, id)
}

func (h Handler) VerifyJwtToken(next http.Handler) http.Handler {
	var ctx *fiber.Ctx
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := &CustomJWT{}
		tokenAuth := r.Header.Get("Authorization")
		_, err := jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
			if jwt.SigningMethodHS256 != token.Method {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			h.ContractUC.EnvConfig, _ = godotenv.Read("../.env")
			secret := h.ContractUC.EnvConfig["TOKEN_SECRET"]
			return []byte(secret), nil
		})

		if err != nil {
			msg := "token is invalid"
			if mErr, ok := err.(*jwt.ValidationError); ok {
				if mErr.Errors == jwt.ValidationErrorExpired {
					msg = "token is expired"
				}
			}

			h.SendErrorResponse(ctx, msg, 200)
			// c.SendAuthError(w, msg)
			return
		}

		// TODO: should check to redis/db is token expired or not
		cx := userContext(r.Context(), "identifier", map[string]string{
			"user_id": claims.UserID,
		})

		next.ServeHTTP(w, r.WithContext(cx))
	})
}

func (h Handler) SendResponse(ctx *fiber.Ctx, data interface{}, err interface{}, code int) error {
	if code == 0 && err != nil {
		code = http.StatusUnprocessableEntity
		err = err.(error).Error()
	}

	if code != http.StatusOK && err != nil {
		return h.SendErrorResponse(ctx, err, code)
	}

	if code == http.StatusAccepted && code != http.StatusOK && err != nil {
		return h.SendAcceptedResponse(ctx, data, code)
	}

	return h.SendSuccessResponse(ctx, data)
}

//send response if status code 201
func (h Handler) SendAcceptedResponse(ctx *fiber.Ctx, data interface{}, meta interface{}) error {
	response := helper.SuccessResponse(data)

	return ctx.Status(http.StatusAccepted).JSON(response)
}

//send response if status code 200
func (h Handler) SendSuccessResponse(ctx *fiber.Ctx, data interface{}) error {
	response := helper.SuccessResponse(data)

	return ctx.Status(http.StatusOK).JSON(response)
}

//send response if status code != 200
func (h Handler) SendErrorResponse(ctx *fiber.Ctx, err interface{}, code int) error {
	response := helper.ErrorResponse(err)

	return ctx.Status(code).JSON(response)
}

//send response if status code 200
func (h Handler) SendFileResponse(ctx *fiber.Ctx, data, contentType string) error {
	fileRes, err := os.OpenFile(data, os.O_RDWR, 0644)
	if err != nil {
		return h.SendErrorResponse(ctx, "a"+err.Error(), http.StatusBadRequest)
	}

	fi, err := fileRes.Stat()
	if err != nil {
		return h.SendErrorResponse(ctx, "b"+err.Error(), http.StatusBadRequest)
	}

	ctx.Set("Content-Type", contentType)
	return ctx.Status(http.StatusOK).SendStream(fileRes, int(fi.Size()))
}
