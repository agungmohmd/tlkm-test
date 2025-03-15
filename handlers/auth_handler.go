package handlers

import (
	"fmt"
	"net/http"

	request "agungmohmd/sematin-front-api/server/requests"
	"agungmohmd/sematin-front-api/usecase"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Handler
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	input := new(request.LoginRequest)
	if err := ctx.BodyParser(input); err != nil {
		fmt.Println(err)
		return h.SendResponse(ctx, nil, err, http.StatusBadRequest)
	}

	authUc := usecase.AuthUC{ContractUC: h.ContractUC}
	res, err := authUc.Login(input)
	if err != nil {
		fmt.Println("on auth handler login func", err)
		return h.SendResponse(ctx, nil, err.Error(), http.StatusBadRequest)
	}

	return h.SendResponse(ctx, res, err, 0)
}

func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	input := new(request.RegisterRequest)
	if err := ctx.BodyParser(input); err != nil {
		fmt.Println(err)
		return h.SendResponse(ctx, nil, err, http.StatusBadRequest)
	}
	fmt.Println(input.Name)

	userUc := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUc.Register(input.Name, input.Username, input.Password, input.Address)
	fmt.Println("here res", res)
	fmt.Println("here err", err)
	if err != nil {
		fmt.Println("on auth handler register func", err)
		return h.SendResponse(ctx, nil, err.Error(), http.StatusBadRequest)
	}

	return h.SendResponse(ctx, res, err, 0)
}

func (h *AuthHandler) UpdateProfile(ctx *fiber.Ctx) error {
	userID := fmt.Sprintf("%v", ctx.Locals("claims"))
	input := new(request.UpdateProfileRequest)
	if err := ctx.BodyParser(input); err != nil {
		fmt.Println(err)
		return h.SendResponse(ctx, nil, err, http.StatusBadRequest)
	}

	userUc := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUc.Register(userID, input.Name, input.Username, input.Address)
	if err != nil {
		fmt.Println("on auth handler update profile func", err)
		return h.SendResponse(ctx, nil, err.Error(), http.StatusBadRequest)
	}

	return h.SendResponse(ctx, res, err, 0)
}
