package handlers

import (
	"log"
	"math/rand"
	"reports-api/db"
	"reports-api/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

var sessions = map[string]string{}

func generateSessionID() string {
	return strconv.FormatInt(rand.Int63(), 20)
}

func generateDummyToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, 32)
	for i := range token {
		token[i] = charset[rand.Intn(len(charset))]
	}
	return string(token)
}

// LoginHandler handles user login
func LoginHandler(c *fiber.Ctx) error {
	var credentials models.Credentials
	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var id int
	var username, password string
	var role string
	err := db.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username = ? AND deleted_at IS NULL", credentials.Username).Scan(&id, &username, &password, &role)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	if credentials.Password != password {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid username or password"})
	}
	sessionID := generateSessionID()
	sessions[sessionID] = username

	c.Set("role", role)
	c.Set("token", generateDummyToken())

	c.Cookie(&fiber.Cookie{
		Name:     "session_cookie",
		Value:    sessionID,
		Path:     "/",
		HTTPOnly: true,
		MaxAge:   3600 * 24,
	})

	log.Printf("User %s logged in successfully", username)
	return c.JSON(models.LoginResponse{
		Message: "Login successful",
		Data:    &models.Data{ID: id, Username: username, Role: role},
	})
}

// RegisterHandler returns a handler for registering a user or admin
func RegisterHandler(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.RegisterUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		var count int
		err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}
		if count > 0 {
			return c.Status(409).JSON(fiber.Map{"error": "Username already exists"})
		}

		_, err = db.DB.Exec(
			"INSERT INTO users (username, password, role, created_by, created_at) VALUES (?, ?, ?, ?, ?)",
			req.Username, req.Password, role, req.CreatedBy, time.Now(),
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to register user"})
		}

		log.Printf("User %s registered successfully as %s", req.Username, role)
		return c.JSON(fiber.Map{"message": "Registered as " + role})
	}
}

// UpdateUserHandler updates user info
func UpdateUserHandler(c *fiber.Ctx) error {
	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ? AND deleted_at IS NULL", req.ID).Scan(&count)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	_, err = db.DB.Exec(
		"UPDATE users SET username = ?, password = ?, updated_by = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		req.Username, req.Password, req.UpdatedBy, time.Now(), req.ID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	log.Printf("Updating user ID: %d with username: %s", req.ID, req.Username)
	return c.JSON(fiber.Map{"message": "User updated"})
}

// DeleteUserHandler deletes a user
func DeleteUserHandler(c *fiber.Ctx) error {
	var req models.DeleteUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_, err := db.DB.Exec(
		"UPDATE users SET deleted_at = ?, deleted_by = ? WHERE id = ? AND deleted_at IS NULL",
		time.Now(), req.DeletedBy, req.ID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	log.Printf("User ID: %d deleted by user ID: %d", req.ID, req.DeletedBy)
	return c.JSON(fiber.Map{"message": "User deleted"})
}

// LogoutHandler logs out a user
func LogoutHandler(c *fiber.Ctx) error {
	sessionID := c.Cookies("session")
	if sessionID != "" {
		delete(sessions, sessionID)
	}
	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
	})

	log.Printf("User logged out successfully")
	return c.JSON(fiber.Map{"message": "Logged out"})
}

func GetresponsHandler(c *fiber.Ctx) error {
	var users []models.User
	rows, err := db.DB.Query("SELECT id, name FROM responsibilities")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}
		users = append(users, user)
	}

	return c.JSON(models.UserResponse{
		Success: true,
		Data:    users,
	})
}

func AddresponsHandler(c *fiber.Ctx) error {
	var req models.ResponseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_, err := db.DB.Exec(
		"INSERT INTO responsibilities (name) VALUES (?)",
		req.Name,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add responsibility"})
	}

	log.Printf("Responsibility %s added successfully", req.Name)
	return c.JSON(fiber.Map{"message": "Responsibility added"})
}

func UpdateResponsHandler(c *fiber.Ctx) error {
	var req models.ResponseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_, err := db.DB.Exec(
		"UPDATE responsibilities SET name = ? WHERE id = ?",
		req.Name, req.ID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update responsibility"})
	}

	log.Printf("Responsibility ID: %d updated successfully", req.ID)
	return c.JSON(fiber.Map{"message": "Responsibility updated"})
}

func DeleteResponsHandler(c *fiber.Ctx) error {
	var req models.ResponseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_, err := db.DB.Exec(
		"DELETE FROM responsibilities WHERE id = ?",
		req.ID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete responsibility"})
	}

	log.Printf("Responsibility ID: %d deleted successfully", req.ID)
	return c.JSON(fiber.Map{"message": "Responsibility deleted"})
}
