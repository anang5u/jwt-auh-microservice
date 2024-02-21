package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = []byte("your-secret-key")

func main() {
	r := gin.Default()

	// Endpoint untuk login, tidak memerlukan token
	r.POST("/login", handleLogin)

	// Grouping untuk endpoint yang memerlukan token
	auth := r.Group("/")
	auth.Use(authMiddleware())
	{
		// Endpoint-endpoint yang memerlukan token
		auth.GET("/secured", handleSecuredEndpoint)
	}

	r.Run(":8080")
}

// Middleware untuk melakukan verifikasi token JWT
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header tidak ditemukan"})
			c.Abort()
			return
		}

		// Parsing token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi algortima
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode penandatanganan tidak valid")
			}
			return secretKey, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		// Jika token valid, lanjutkan ke handler
		if token.Valid {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
		}
	}
}

// Handler untuk endpoint login
func handleLogin(c *gin.Context) {
	// Proses autentikasi, biasanya dengan username dan password
	// Setelah berhasil, kita akan membuat token JWT
	token := jwt.New(jwt.SigningMethodHS256)

	// Menambahkan klaim ke token JWT, misalnya role atau scope
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = "exampleuser"
	claims["role"] = "admin" // Contoh klaim role admin

	// Set waktu kadaluarsa token
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Menandatangani token dengan secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Mengirimkan token sebagai respons
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Handler untuk endpoint yang memerlukan token
func handleSecuredEndpoint(c *gin.Context) {
	// Mengambil klaim dari token untuk digunakan dalam aksi tertentu
	claims := c.MustGet("claims").(jwt.MapClaims)
	username := claims["username"].(string)
	role := claims["role"].(string)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Endpoint yang memerlukan token",
		"username": username,
		"role":     role,
	})
}
