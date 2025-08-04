package middleware

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
}

func AuthMiddleware(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := c.Get("Authorization")
		fmt.Println("🔐 Raw Authorization header:", accessToken)

		if accessToken == "" {
			fmt.Println("❌ Authorization header is missing")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization token",
			})
		}

		// Remove "Bearer " prefix if it exists
		if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
			accessToken = accessToken[7:]
		}
		fmt.Println("🔍 Access token to parse:", accessToken)

		token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			alg := token.Method.Alg()
			fmt.Println("🔑 Token signing method:", alg)

			if alg != jwt.SigningMethodHS256.Alg() {
				fmt.Printf("❌ Unexpected signing method: %v\n", alg)
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				fmt.Println("❌ Token has expired")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "token expired",
				})
			}
			fmt.Println("❌ Token parsing error:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		if claims, ok := token.Claims.(*Claims); !ok || !token.Valid {
			fmt.Printf("❌ Invalid token claims or token is not valid: %v | claims: %v\n", token.Valid, token.Claims)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		} else {
			fmt.Println("✅ Token is valid. Claims:", claims)
		}

		return c.Next()
	}
}
