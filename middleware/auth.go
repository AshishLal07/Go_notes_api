package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"notes-api/config"
	"notes-api/models"
	"notes-api/utils"
)

// JWTMiddleware validates JWT tokens and sets user context
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Authorization header is required",
			})
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Authorization header must start with 'Bearer '",
			})
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Token is required",
			})
		}

		// Validate token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid or expired token",
			})
		}

		// Verify user exists in database
		var user models.User
		if err := config.GetDB().First(&user, claims.UserID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "User not found",
			})
		}

		// Set user in context
		c.Locals("user", &user)
		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}

// GetUserFromContext retrieves the authenticated user from context
func GetUserFromContext(c *fiber.Ctx) (*models.User, error) {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "User not found in context")
	}
	return user, nil
}

// GetUserIDFromContext retrieves the authenticated user ID from context
func GetUserIDFromContext(c *fiber.Ctx) (uint, error) {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}
	return userID, nil
}
