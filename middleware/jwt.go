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

func GenerateJWT(email, role, id string) (string, error ) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email" : email,
		"role" : role,
		"userID" : id,
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
		userID, userExist := claims["userID"].(string)

		if !emailExist || !roleExist || !userExist{
			c.JSON(http.StatusUnauthorized, gin.H{
				"error" : "Token tidak mengandung informasi pengguna",
			})
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("role", role)
		c.Set("userID", userID)
		
		c.Next()
	}
}

// Middleware untuk memeriksa role pengguna 

func RoleCheck(allowRoles ...string) gin.HandlerFunc {
	return func (c *gin.Context)  {
		// Ambil role dari context
		role, exist := c.Get("role")
		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{"error" : "Role tidak ditemukan"})
			c.Abort()
			return
		}

		roleString, ok := role.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error" : "Role tidak Valid"})
			c.Abort()
			return
		}

		// Cek role apakah di perbolehkan
		allowed :=  false 
		for _, allowedRole := range allowRoles {
			if roleString == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error" : "Anda tidak memiliki akses"})
			c.Abort()
			return
		}

		
	}
}

func AdminOnlyMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		//Ambil Role 
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error" : "only admin can access",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
