package middlewares

import (
	"log"
	"net/http"
	"raihpeduli/helpers"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthorizeJWT(jwtService helpers.JWTInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				response := helpers.BuildErrorResponse("No Token Found !")
				return c.JSON(http.StatusBadRequest, response)
			}

			tokenString := authHeader[len("Bearer "):]
			token, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				log.Println(err)
				response := helpers.BuildErrorResponse("Token is not valid -" + err.Error())
				return c.JSON(http.StatusUnauthorized, response)
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID, _ := strconv.Atoi(claims["user_id"].(string))
				
				c.Set("user_id", userID)
				c.Set("role_id", claims["role_id"])
				return next(c)
			}

			response := helpers.BuildErrorResponse("Invalid Token Claims")
			return c.JSON(http.StatusUnauthorized, response)
		}
	}
}
