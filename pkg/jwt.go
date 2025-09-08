package pkg

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint16 `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTClaims(uid uint16, role string) *Claims {
	return &Claims{
		UserID: uid,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Minute)),
			Issuer:    os.Getenv("JWT_ISSUER"),
		},
	}
}

func (c *Claims) GenAccessToken() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("no secrets found")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(jwtSecret))
}

func (c *Claims) ValidateToken(token string) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("no secrets found")
	}

	parsedToken, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, jwt.ErrTokenExpired) {
			return errors.New("token already expired")
		}
		return errors.New("unable to parsing access token")
	}

	iss, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return errors.New("unable to get issuer")
	}
	if iss != os.Getenv("JWT_ISSUER") {
		return errors.New("access token mismatch")
	}

	return nil
}
