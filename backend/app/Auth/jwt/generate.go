package jwt

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var secretKey string

func CreateJwtPayload(expiration time.Duration, id int, username string, sessionID string) JwtPayload {
	iat := time.Now().Unix()
	exp := iat + int64(expiration.Seconds())

	return JwtPayload{
		Sub:       id,
		Username:  username,
		SessionID: sessionID,
		Iat:       iat,
		Exp:       exp,
	}
}

// Generate creates a JWT with a given payload
func Generate(payload JwtPayload) string {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)

	encodedHeader := base64Encode(headerJSON)
	encodedPayload := base64Encode(payloadJSON)

	signature := signMessage(encodedHeader + "." + encodedPayload)

	return encodedHeader + "." + encodedPayload + "." + signature
}

// Verify checks if the JWT signature is valid
func JWTVerify(token string) (*JwtPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	decodedPayload, err := base64Decode(parts[1])
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}

	if signMessage(parts[0]+"."+parts[1]) != parts[2] {
		return nil, errors.New("invalid signature")
	}

	var payload JwtPayload
	err = json.Unmarshal(decodedPayload, &payload)
	if err != nil {
		return nil, errors.New("invalid payload JSON")
	}
	now := time.Now().Unix()
	if payload.Exp < now {
		return nil, errors.New("token expired")
	}
	return &payload, nil
}
