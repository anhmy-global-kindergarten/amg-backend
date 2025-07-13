package auth

import (
	"amg-backend/config"
	"amg-backend/models"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type RegisterRequest struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=20"`
	Name            string `json:"name" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username"`
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
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Password != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Passwords do not match"})
	}

	collection := h.DB.Database(config.DBName).Collection("User")
	count, err := collection.CountDocuments(context.TODO(), bson.M{"username": req.Username})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: req.Username,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     "parent",
		IsActive: true,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login credentials"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /amg/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	collection := h.DB.Database(config.DBName).Collection("User")

	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Tài khoản hoặc mât khẩu không đúng. Vui lòng thử lại."})
	}
	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Tài khoản của bạn đã bị vô hiệu hóa. Vui lòng liên hệ admin."})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Tài khoản hoặc mât khẩu không đúng. Vui lòng thử lại."})
	}

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var jwtSecret = []byte(config.SecretKey)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not login, failed to sign token"})
	}

	cookie := fiber.Cookie{
		Name:     "session_token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}

	c.Cookie(&cookie)

	user.Password = ""
	return c.JSON(user)
}

// GetCurrentUser godoc
// @Summary Get current logged-in user's info
// @Description Get user info from session token cookie
// @Tags auth
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /amg/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	cookie := c.Cookies("session_token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "not authenticated"})
	}
	var jwtSecret = []byte(config.SecretKey)
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
	}

	userInfo := fiber.Map{
		"id":       claims["id"],
		"username": claims["username"],
		"role":     claims["role"],
	}

	return c.JSON(userInfo)
}

// Logout godoc
// @Summary Logout user
// @Description Clear the session token cookie
// @Tags auth
// @Produce json
// @Success 200 {objects} map[string]string
// @Router /amg/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
}
