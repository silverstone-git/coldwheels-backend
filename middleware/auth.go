package middleware

import (
  "fmt"
  "strings"
  "net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string   `json:"userID"`
    jwt.RegisteredClaims
}


func AuthMiddleware(JwtSecret []byte) gin.HandlerFunc {

  // fmt.Println("jwt secret is: ", JwtSecret)

  return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return

		}

    // Check if the header has the Bearer prefix
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
        c.Abort()
        return
    }

    tokenString := parts[1]

    // Parse and validate the token
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return JwtSecret, nil
    })

    if err != nil {
        if err == jwt.ErrTokenExpired {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
        }
        c.Abort()
        return
    }

    if !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unexpected error"})
        c.Abort()
        return
    }

    c.Set("claims", claims)
    c.Set("UserID", claims.UserID)

    fmt.Println("uid during auth middleware: ", claims.UserID)

		c.Next()
	}
}
