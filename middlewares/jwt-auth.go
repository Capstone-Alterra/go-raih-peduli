package middlewares

import (
	"log"
	"net/http"
	"raihpeduli/helpers"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthorizeJWT(jwtService helpers.JWTInterface, role int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authHeader := ctx.Request().Header.Get("Authorization")
			if authHeader == "" {
				response := helpers.BuildErrorResponse("no token found")
				return ctx.JSON(http.StatusBadRequest, response)
			}

			tokenString := authHeader[len("Bearer "):]
			token, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				log.Println(err)
				response := helpers.BuildErrorResponse("token is not valid - " + err.Error())
				return ctx.JSON(http.StatusUnauthorized, response)
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID, _ := strconv.Atoi(claims["user_id"].(string))
				
				ctx.Set("user_id", userID)
				ctx.Set("role_id", claims["role_id"])

				if claims["role_id"].(int) == 0 {
					return next(ctx)
				}

				if claims["role_id"].(int) != role {
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
