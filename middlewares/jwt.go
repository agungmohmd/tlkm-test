package middlewares

import (
	"agungmohmd/sematin-front-api/server/handlers"
	"agungmohmd/sematin-front-api/usecase"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

// JwtMiddleware ...
type JwtMiddleware struct {
	ContractUC *usecase.ContractUC
	Roles      []string
	PhoneValid bool
}

// CustomClaim ...
type CustomClaim struct {
	jwt.StandardClaims
}

// New jwt middleware
func (jwtMiddleware JwtMiddleware) New(ctx *fiber.Ctx) (err error) {
	claims := &CustomClaim{}
	handler := handlers.Handler{ContractUC: jwtMiddleware.ContractUC}

	//check header is present or not
	header := ctx.Get("Authorization")
	if !strings.Contains(header, "Bearer") {
		return handler.SendResponse(ctx, nil, "Header Not Present", http.StatusUnauthorized)
	}

	//check claims and signing method
	token := strings.Replace(header, "Bearer ", "", -1)
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		jwtMiddleware.ContractUC.EnvConfig, _ = godotenv.Read("../.env")
		secret := []byte(jwtMiddleware.ContractUC.EnvConfig["TOKEN_SECRET"])
		return secret, nil
	})
	if err != nil {
		return handler.SendResponse(ctx, nil, "unexpected claim methods", http.StatusUnauthorized)
	}

	//check token live time
	if claims.ExpiresAt < time.Now().Unix() {
		return handler.SendResponse(ctx, nil, "token expired", http.StatusUnauthorized)
	}

	// claims.Id =
	//set id to uce case contract
	// claims.Id = fmt.Sprintf("%v", jweRes["id"])
	ctx.Locals("claims", claims.Id)

	jwtMiddleware.ContractUC.UserID = claims.Id

	return ctx.Next()
}
