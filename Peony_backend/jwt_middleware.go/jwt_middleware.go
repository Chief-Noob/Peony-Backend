package jwt_middleware

import (
	"Peony/Peony_backend/models/entity"
	"Peony/config"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte(config.GetSecretKey())

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authentication")
		if len(authHeader) == 0 {
			c.JSON(400, gin.H{
				"error": "No authentication header.",
			})
			c.Abort()
			return
		}

		token := strings.Split(authHeader, " ")[1]
		tokenClaims, err := jwt.ParseWithClaims(token, &entity.Claims{}, func(token *jwt.Token) (i interface{}, err error) {
			return jwtSecret, nil
		})
		if err != nil {
			var message string
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					message = "token is malformed"
				} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
					message = "token could not be verified because of signing problems"
				} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					message = "signature validation failed"
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					message = "token is expired"
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					message = "token is not yet valid before sometime"
				} else {
					message = "can not handle this token"
				}
			}
			c.JSON(401, gin.H{
				"error": message,
			})
			c.Abort()
			return
		}
		if _, ok := tokenClaims.Claims.(*entity.Claims); ok && tokenClaims.Valid {
			c.Next()
		} else {
			fmt.Println("========================")
			c.Abort()
			return
		}
	}
}
