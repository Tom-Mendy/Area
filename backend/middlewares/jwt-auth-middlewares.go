package middlewares

import (
	"log"
	"net/http"

	"area/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthorizeJWT validates the token from the http request, returning a 401 if it's not valid.
func AuthorizeJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		tokenString := authHeader[len("Bearer "):]

		token, err := service.NewJWTService().ValidateToken(tokenString)

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claims: ", claims)
			log.Println("Claims[Id]: ", claims["jti"])
			log.Println("Claims[Name]: ", claims["name"])
			log.Println("Claims[Admin]: ", claims["admin"])
			log.Println("Claims[Issuer]: ", claims["iss"])
			log.Println("Claims[IssuedAt]: ", claims["iat"])
			log.Println("Claims[ExpiresAt]: ", claims["exp"])
		} else {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
