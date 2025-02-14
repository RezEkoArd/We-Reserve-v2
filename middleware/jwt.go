package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var secretKey = []byte(viper.GetString("JWT_SECRET_KEY"))

func GenerateJWT(email, role string) (string, error ) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email" : email,
		"role" : role,
		"exp" : time.Now().Add(time.Hour * 3).Unix(), // expired 3 day
	})

	return token.SignedString(secretKey)
}


// fungsi middleware untuk mengambil data di auth apakah ada header authentication atau tidak
func JWTAuthMiddleware() gin.HandlerFunc {
	return func (c *gin.Context)  {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error" : "Token tidak ditemukan",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error " : "Token tidak valid",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error" : "Gagal Memproses klaim token",
			})
			c.Abort()
			return
		}

		email, emailExist := claims["email"].(string)
		role, roleExist := claims["role"].(string)

		if !emailExist || !roleExist {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error" : "Token tidak mengandung informasi pengguna",
			})
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("role", role)
		
		c.Next()
	}
}

// Middleware untuk memeriksa role pengguna 

func RoleCheck(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil role dari context
		role, exist := c.Get("role")
		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error" : "Role tidak ditemukan",
			})
			c.Abort()
			return
		}

		// Periksa apakah role pengguna termaksud dalam allow Roles
		roleString, ok := role.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error" : "Role tidak valid",
			})
			c.Abort()
			return
		}

		// Cek apakah role pengguna diizinkan
		allowed := false
		for _, allowedRole := range allowedRoles {
			if roleString == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error" : "Anda tidak memiliki akses",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}