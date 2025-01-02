package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"

	"area/schemas"
)

type JWTService interface {
	GenerateToken(userID string, name string, admin bool) string
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetUserIdfromJWTToken(tokenString string) (userID uint64, err error)
}

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    "email@example.com",
	}
}

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET is not set")
	}
	return secret
}

func (jwtSrv *jwtService) GenerateToken(userID string, username string, admin bool) string {
	// Set custom and standard claims
	claims := &jwtCustomClaims{
		username,
		admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * schemas.BearerTokenDuration).Unix(),
			Issuer:    jwtSrv.issuer,
			IssuedAt:  time.Now().Unix(),
			Id:        userID,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token using the secret signing key
	t, err := token.SignedString([]byte(jwtSrv.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (jwtSrv *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	result, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Signing method validation
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret signing key
		return []byte(jwtSrv.secretKey), nil
	})
	return result, err
}

func (jwtSrv *jwtService) GetUserIdfromJWTToken(tokenString string) (userID uint64, err error) {
	token, err := jwtSrv.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		if jti, ok := claims["jti"].(string); ok {
			id, err := strconv.ParseUint(jti, 10, 64)
			if err != nil {
				return 0, errors.New("jti claim is not a float64")
			}
			return id, nil
		}
		return 0, errors.New("jti claim is not a float64")
	} else {
		return 0, errors.New("invalid token")
	}
}
