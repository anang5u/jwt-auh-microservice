package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	signingKey *rsa.PrivateKey
	verifyKey  *rsa.PublicKey
)

func init() {
	// Generate RSA key pair
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to generate RSA key pair: %v", err)
	}
	signingKey = key
	verifyKey = &key.PublicKey
}

func main() {
	r := gin.Default()

	// Endpoint untuk login dan menghasilkan token JWT
	r.POST("login", func(c *gin.Context) {
		// Buat token JWT
		token := jwt.New(jwt.SigningMethodRS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["user_id"] = "11223344"
		claims["email"] = "userdemo@example.com"
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token kadaluarsa dalam 24 jam

		// Tandatangani token menggunakan kunci privat
		tokenString, err := token.SignedString(signingKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
			return
		}

		// Kirim token sebagai respons
		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
			"user": map[string]string{
				"user_id":     claims["user_id"].(string),
				"first_name":  "Example",
				"middle_name": "Demo",
				"last_name":   "User",
				"email":       claims["email"].(string),
			},
		})
	})

	// Endpoint untuk JWKS URI
	r.GET("jwks", func(c *gin.Context) {
		// Menghasilkan JWKS response
		jwks := map[string]interface{}{
			"keys": []map[string]interface{}{
				{
					"kid": "1",
					"use": "sig",
					"kty": "RSA",
					"alg": "RS256",
					"n":   fmt.Sprintf("%x", verifyKey.N),
					"e":   "AQAB",
					"exp": time.Now().Add(time.Hour * 24).Unix(), // Contoh: menetapkan masa berlaku token selama 24 jam
				},
			},
		}

		// Kirim respons dengan JWKS URI
		c.JSON(http.StatusOK, jwks)
	})

	// Menjalankan server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
