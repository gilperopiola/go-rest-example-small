package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthI interface {
	GenerateToken(id int, username, email string, role Role) (string, error)
	ValidateToken(role Role, shouldMatchUserID bool) gin.HandlerFunc
}

func NewAuth(secret string, sessionDurationDays int) *Auth {
	return &Auth{
		secret:              secret,
		sessionDurationDays: sessionDurationDays,
	}
}

type Auth struct {
	secret              string
	sessionDurationDays int
}

type Role string

const (
	AnyRole   Role = "any"
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type CustomClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
	jwt.RegisteredClaims
}

func (auth *Auth) GenerateToken(id int, username, email string, role Role) (string, error) {

	var (
		issuedAt  = time.Now()
		expiresAt = time.Now().Add(time.Hour * 24 * time.Duration(auth.sessionDurationDays))
	)

	// Generate claims containing Username, Email, Role and ID
	claims := &CustomClaims{
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        fmt.Sprint(id),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	// Generate token (struct)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token (string)
	return token.SignedString([]byte(auth.secret))
}

// ValidateToken validates a token for a specific role and sets ID and Email in context
func (auth *Auth) ValidateToken(role Role, shouldMatchUserID bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get token string and then convert it to a *jwt.Token
		token, err := auth.getTokenStructFromContext(c)
		if err != nil {
			c.Error(Wrap("auth.getTokenStructFromContext", ErrUnauthorized))
			c.Abort()
			return
		}

		// Get custom claims from token
		customClaims, ok := token.Claims.(*CustomClaims)

		// Check if claims and token and role are valid
		if !ok || !token.Valid || customClaims.Valid() != nil || (role != AnyRole && customClaims.Role != role) {
			c.Error(Wrap("!token.Valid || customClaims.Role != role", ErrUnauthorized))
			c.Abort()
			return
		}

		// Check if user ID in URL matches user ID in token
		if shouldMatchUserID {
			urlUserID, err := getIntFromURLPath(c.Params, pathUserIDKey)
			if err != nil || customClaims.ID != fmt.Sprint(urlUserID) {
				c.Error(Wrap("!shouldMatchUserID", ErrUnauthorized))
				c.Abort()
				return
			}
		}

		// If OK, set UserID, Username and Email inside of context
		userID, _ := strconv.Atoi(customClaims.ID)
		addUserInfoToContext(c, userID, customClaims.Username, customClaims.Email)
	}
}

func addUserInfoToContext(c *gin.Context, id int, username, email string) {
	c.Set("UserID", id)
	c.Set("Username", username)
	c.Set("Email", email)
}

func (auth *Auth) getTokenStructFromContext(c *gin.Context) (*jwt.Token, error) {

	// Get token string from headers
	tokenString := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer ")

	// Decode string into actual *jwt.Token
	token, err := auth.decodeTokenString(tokenString)
	if err != nil {
		return nil, err
	}

	// Token decoded OK
	return token, nil
}

// decodeTokenString decodes a JWT token string into a *jwt.Token
func (auth *Auth) decodeTokenString(tokenString string) (*jwt.Token, error) {

	// Check length
	if len(tokenString) < 40 {
		return &jwt.Token{}, ErrUnauthorized
	}

	// Make key function
	keyFunc := func(token *jwt.Token) (interface{}, error) { return []byte(auth.secret), nil }

	// Parse
	return jwt.ParseWithClaims(tokenString, &CustomClaims{}, keyFunc)
}

func getIntFromURLPath(params gin.Params, key string) (int, error) {
	value, ok := params.Get(key)
	if !ok {
		return 0, fmt.Errorf("error getting %s from URL params", key)
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("error converting %s from string to int", key)
	}

	return valueInt, nil
}
