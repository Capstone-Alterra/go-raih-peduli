package middlewares

import (
	"log"
	"net/http"
	"raihpeduli/helpers"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthorizeJWT(jwtService helpers.JWTInterface, role int, secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authHeader := ctx.Request().Header.Get("Authorization")
			if authHeader == "" {
				if role == -1 {
					return next(ctx)
				}
				response := helpers.BuildErrorResponse("no token found")
				return ctx.JSON(http.StatusBadRequest, response)
			}

			tokenString := authHeader[len("Bearer "):]
			token, err := jwtService.ValidateToken(tokenString, secret)
			if err != nil {
				log.Println(err)
				response := helpers.BuildErrorResponse("token is not valid - " + err.Error())
				return ctx.JSON(http.StatusUnauthorized, response)
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID, _ := strconv.Atoi(claims["user_id"].(string))
				roleID, _ := strconv.Atoi(claims["role_id"].(string))

				ctx.Set("user_id", userID)
				ctx.Set("role_id", roleID)

				if role == 0 || role == -1 || roleID == 3 {
					return next(ctx)
				}

				if roleID != role {
					response := helpers.BuildErrorResponse("this user cannot access this endpoint")
					return ctx.JSON(http.StatusUnauthorized, response)
				}

				return next(ctx)
			}

			response := helpers.BuildErrorResponse("invalid token claims")
			return ctx.JSON(http.StatusUnauthorized, response)
		}
	}
}
