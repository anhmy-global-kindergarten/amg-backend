package auth

import (
	"amg-backend/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "User data"
// @Success 200 {object} map[string]string
// @Router /livenest/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.Users{
		Username:     data["username"],
		Email:        data["email"],
		PasswordHash: string(password),
		Role:         "viewer",
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
	}

	return c.JSON(fiber.Map{"message": "user created successfully"})
}

func (h *AuthHandler) Login(c *fiber.Ctx, db *gorm.DB) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	var user models.Users
	db.Where("email = ?", data["email"]).First(&user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "incorrect password"})
	}

	return c.JSON(fiber.Map{"message": "login successful"})
}
