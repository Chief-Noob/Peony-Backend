package jwt_middleware

import (
	"Peony/Peony_backend/models/entity"
	"Peony/config"
	"strings"

	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwt_secret = []byte(config.GetSecretKey())

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth_header := c.Request.Header.Get("Authentication")
		if len(auth_header) == 0 {
			c.JSON(400, gin.H{
				"error": "NO AUTHENTICATION HEADER.",
			})
			c.Abort()
			return
		}
		mat, err := regexp.MatchString(`token.*`, auth_header)

		if err != nil || mat == false {
			c.JSON(401, gin.H{
				"error": "TOKEN IS MALFORMED.",
			})
			c.Abort()
			return
		}
		token := strings.Split(auth_header, " ")[1]

		token_claims, err := jwt.ParseWithClaims(token, &entity.Claims{}, func(token *jwt.Token) (i interface{}, err error) {
			return jwt_secret, nil
		})

		if err != nil {
			var message string
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					message = "TOKEN IS MALFORMED."
				} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
					message = "TOKEN NOT VERIFIED."
				} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					message = "SIGNATURE VALIDATION FAILED."
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					message = "TOKEN IS EXPIRED."
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					message = "TOKEN IS NOT YET VALID."
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
		if _, ok := token_claims.Claims.(*entity.Claims); ok && token_claims.Valid {
			c.Next()
		} else {
			c.Abort()
			return
		}
	}
}
