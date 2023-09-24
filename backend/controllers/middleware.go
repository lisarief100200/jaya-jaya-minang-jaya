package controllers

import (
	"api-lisa/pkg/v1/mysql"
	"api-lisa/utils/config"
	"api-lisa/utils/log"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func MiddlewareFuncOverrideAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Periksa apakah tokenString memiliki prefix "Bearer "
		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token not provided"})
			c.Abort()
			return
		}

		// Ambil bagian token dari tokenString
		token := tokenString[len("Bearer "):]

		// Verifikasi token bearer.
		claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.MyConfig.SecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		// Ambil uid dari claims JWT.
		uid, ok := claims.Claims.(jwt.MapClaims)["uid"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		if !IsValidToken(uid, token) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		// Tambahkan UID ke dalam konteks Gin.
		c.Set("uid", uid)

		c.Next()
	}
}

func IsValidToken(uid, tokenString string) bool {
	// Get connection DB
	db, err := mysql.GetConnectionAdmin()
	if err != nil {
		log.Log.Errorln("Error GetConnectionAdmin", err.Error())
		return false
	}

	var token string

	// Ubah query sesuai dengan kebutuhan Anda.
	query := "SELECT token FROM tbl_session WHERE user_id = ? AND status = 1"

	if err := db.QueryRow(query, uid).Scan(&token); err != nil {
		log.Log.Errorln("Error running query GetAdminProp", err.Error())
		return false
	}

	if token != tokenString {
		return false
	}

	return true
}

func GetUid(c *gin.Context) (string, error) {
	uid, valid := c.Get("uid")
	if !valid {
		return "", errors.New("undefined User Id")
	}
	return uid.(string), nil
}
