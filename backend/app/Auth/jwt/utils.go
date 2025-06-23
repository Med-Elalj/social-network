package jwt

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strings"
	"time"

	"social-network/server/logs"
)

const Time_to_Expire = time.Hour * 6

type JwtPayload struct {
	Sub      int    `json:"sub,string"`
	Username string `json:"username"`
	Iat      int64  `json:"iat"`
	Exp      int64  `json:"exp"`
}

// LoadSecret manually reads the .env file and retrieves JWT_SECRET_KEY
func LoadSecret() string {
	file, err := os.Open("../private/.env")
	if err != nil {
		logs.Fatal("Error loading private/.env file", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "JWT_SECRET_KEY=") {
			return strings.TrimPrefix(line, "JWT_SECRET_KEY=")
		}
	}

	logs.Fatal("JWT_SECRET_KEY not found in private/.env file")
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

func signMessage(message string) string {
	if secretKey == "" {
		secretKey = LoadSecret()
		if secretKey == "" {
			logs.Fatal("JWT_SECRET_KEY is not set")
		}
	}
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	signature := h.Sum(nil)
	return base64Encode(signature)
}
