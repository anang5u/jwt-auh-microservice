package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	//"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	issuerKey = "your-issuer"
)

type JSONWebKeySet struct {
	Keys []JSONWebKey `json:"keys"`
}
type JSONWebKey struct {
	Kty string `json:"kty"`
	N   string `json:"n"`
	E   string `json:"e"`
	D   string `json:"d"`
	P   string `json:"p"`
	Q   string `json:"q"`
	Dp  string `json:"dp"`
	Dq  string `json:"dq"`
	Qi  string `json:"qi"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
}

func main() {
	r := gin.Default()

	// Endpoint untuk login dan menghasilkan token JWT
	r.POST("login", func(c *gin.Context) {
		// JWK
		jwks, err := getJwks()
		if err != nil {
			fmt.Printf("failed to get jwks, err: %s\n", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get jwks"})
			return
		}
		keys := jwks.Keys[0]

		// jwkByte
		jwkByte, err := json.Marshal(&keys)
		if err != nil {
			fmt.Printf("failed to convert jwks, err: %s\n", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert jwks"})
			return
		}
		jwkStr := string(jwkByte)

		// Build a JWT!
		// create a new jwt
		token := jwt.New()
		token.Set(jwt.SubjectKey, `https://anangsu13.cloud`)
		token.Set(jwt.IssuerKey, issuerKey)
		token.Set(jwt.IssuedAtKey, time.Unix(500, 0))

		/*
			buf, err := json.MarshalIndent(token, "", "  ")
			if err != nil {
				fmt.Printf("failed to generate JSON: %s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate JSON"})
				return
			}
			fmt.Printf("%s\n", buf)
		*/

		if v, ok := token.Get(`privateClaimKey`); ok {
			fmt.Printf("privateClaimKey -> '%s'\n", v)
		}

		//convert jwk in bytes and return a new key
		//jwkey, err := jwk.ParseKey([]byte(jwkStr))
		jwkey, err := jwk.ParseKey([]byte(jwkStr))

		if err != nil {
			fmt.Printf("failed to parse key: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse key"})
			return
		}

		// signed and return a jwt
		signed, _ := jwt.Sign(token, jwt.WithKey(jwa.RS256, jwkey))

		fmt.Println(string(signed[:]))

		/*
			token, err := jwt.NewBuilder().
				Issuer(issuer).
				IssuedAt(time.Now()).
				Build()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to build token: %s", err)})
				return
			}

			// Sign a JWT!
			signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, privateKey))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to sign token: %s\n", err)})
				return
			}

			// Verify a JWT!
			verifiedToken, err := jwt.Parse(signed, jwt.WithKey(jwa.RS256, publicKey))
			if err != nil {
				fmt.Printf("failed to verify JWS: %s\n", err)
				return
			}
			_ = verifiedToken

			fmt.Printf("%s\n", verifiedToken)

			// Buat token JWT
			/*
				token := jwt.New(jwt.SigningMethodRS256)
				claims := token.Claims.(jwt.MapClaims)
				claims["user_id"] = "11223344"
				claims["email"] = "userdemo@example.com"
				claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token kadaluarsa dalam 24 jam
				claims["iss"] = "your-issuer"                         // Ganti dengan nilai issuer yang sesuai

				// Tandatangani token menggunakan kunci privat
				tokenString, err := token.SignedString(signingKey)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
					return
				}
		*/

		// Kirim token sebagai respons
		c.JSON(http.StatusOK, gin.H{
			"token": string(signed),
			"user": map[string]string{
				"user_id": "111", //claims["user_id"].(string),
			},
		})
	})

	// Endpoint untuk JWKS URI
	r.GET("jwks", func(c *gin.Context) {
		jwks, err := getJwks()
		if err != nil {
			fmt.Printf("Generate jwks failed, err: %s\n", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Generate jwks failed"})
			return
		}

		// Kirim respons dengan JWKS URI
		c.JSON(http.StatusOK, jwks)
	})

	// Menjalankan server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func getJwks() (*JSONWebKeySet, error) {
	/*
		// Generate RSA private key
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			fmt.Println("Failed to generate RSA key pair:", err)
			return
		}

		// Marshal the private key to PEM format
		privateKeyPEM := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}
		privateKeyFile, err := os.Create("private_key.pem")
		if err != nil {
			fmt.Println("Failed to create private key file:", err)
			return
		}
		defer privateKeyFile.Close()

		// Write the private key to the file
		if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
			fmt.Println("Failed to write private key to file:", err)
			return
		}
	*/

	key, err := jwk.ReadFile("private_key.pem", jwk.WithPEM(true))
	if err != nil {
		return nil, fmt.Errorf("failed to parse key in PEM format: %s", err.Error())
	}

	// Mengkonversi *jose.JSONWebKey ke JSON
	keyJSON, err := json.Marshal(key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON Web Key: %s", err.Error())
	}

	// Deklarasikan objek jwk.Set
	jwkSet := JSONWebKeySet{}

	// Lakukan Unmarshal JSON ke dalam objek jwk.Set
	if err := json.Unmarshal(keyJSON, &jwkSet); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %s", err.Error())
	}

	return &jwkSet, nil
}
