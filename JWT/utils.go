package jwt

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"os"
	"strings"
	"time"
)

var Time_to_Expire = time.Hour * 6

type JwtPayload struct {
	Sub      int    `json:"sub,string"`
	Username string `json:"username"`
	Iat      int64  `json:"iat"`
	Exp      int64  `json:"exp"`
}

// LoadSecret manually reads the .env file and retrieves JWT_SECRET_KEY
func LoadSecret() string {
	file, err := os.Open("private/.env")
	if err != nil {
		log.Fatal("Error loading private/.env file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "JWT_SECRET_KEY=") {
			return strings.TrimPrefix(line, "JWT_SECRET_KEY=")
		}
	}

	log.Fatal("JWT_SECRET_KEY not found in private/.env file")
	return ""
}

// base64Encode encodes data to a URL-safe base64 string
func base64Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// base64Decode decodes a URL-safe base64 string
func base64Decode(encoded string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(encoded)
}

func signMessage(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	signature := h.Sum(nil)
	return base64Encode(signature)
}
