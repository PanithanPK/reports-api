package handlers

import (
	"log"
	"math/rand"
	"reports-api/db"
	"reports-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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

// @Summary User login
// @Description Authenticate user and create session
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body models.Credentials true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/authEntry/login [post]
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

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(credentials.Password)); err != nil {
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

// @Summary Register user
// @Description Register a new user or admin
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body models.RegisterUserRequest true "User registration data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/authEntry/registerUser [post]
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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
		}

		_, err = db.DB.Exec(
			"INSERT INTO users (username, password, role) VALUES (?, ?, ?)",
			req.Username, string(hashedPassword), role,
		)

		if err != nil {
			log.Printf("Registering user %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to register user"})
		}

		log.Printf("User %s registered successfully as %s", req.Username, role)
		return c.JSON(fiber.Map{
			"message":  "Registered as " + role,
			"username": req.Username,
			"role":     role,
		})
	}
}

// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UpdateUserRequest true "User update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/authEntry/updateUser [put]
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

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	_, err = db.DB.Exec(
		"UPDATE users SET username = ?, password = ?, role = ?, updated_at=CURRENT_TIMESTAMP WHERE id = ?",
		req.Username, string(hashedPassword), req.Role, req.ID,
	)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	log.Printf("Updating user ID: %d with username: %s", req.ID, req.Username)
	return c.JSON(fiber.Map{
		"message":  "User updated",
		"username": req.Username,
	})
}

// @Summary Delete user
// @Description Delete a user (soft delete)
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.DeleteUserRequest true "User delete data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/authEntry/deleteUser [delete]
func DeleteUserHandler(c *fiber.Ctx) error {
	var req models.DeleteUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Hard delete - remove user completely from database
	_, err := db.DB.Exec(
		"DELETE FROM users WHERE id = ?",
		req.ID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	log.Printf("User ID: %d permanently deleted by user ID: %d", req.ID, req.DeletedBy)
	return c.JSON(fiber.Map{"message": "User permanently deleted"})
}

// @Summary User logout
// @Description Log out user and clear session
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/authEntry/logout [post]
func LogoutHandler(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_cookie")
	if sessionID != "" {
		delete(sessions, sessionID)
	}
	c.Cookie(&fiber.Cookie{
		Name:     "session_cookie",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
	})

	log.Printf("User logged out successfully")
	return c.JSON(fiber.Map{"message": "Logged out"})
}

// @Summary Get responsibilities
// @Description Get all responsibilities
// @Tags responsibilities
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/respons/list [get]
func GetresponsHandler(c *fiber.Ctx) error {
	rows, err := db.DB.Query("SELECT id, IFNULL(telegram_username, '') as telegram_username, COALESCE(name, '') as name FROM responsibilities")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()

	var responsibilities []models.ResponseRequest
	for rows.Next() {
		var resp models.ResponseRequest
		if err := rows.Scan(&resp.ID, &resp.TelegramUsername, &resp.Name); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}
		responsibilities = append(responsibilities, resp)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responsibilities,
	})
}

// @Summary Add responsibility
// @Description Add a new responsibility
// @Tags responsibilities
// @Accept json
// @Produce json
// @Param responsibility body models.ResponseRequest true "Responsibility data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/respons/create [post]
func AddresponsHandler(c *fiber.Ctx) error {
	var req models.ResponseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_, err := db.DB.Exec(
		"INSERT INTO responsibilities (name, telegram_username) VALUES (?, ?)",
		req.Name, req.TelegramUsername,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add responsibility"})
	}

	log.Printf("Responsibility %s added successfully", req.Name)
	return c.JSON(fiber.Map{"message": "Responsibility added"})
}

// @Summary Update responsibility
// @Description Update an existing responsibility
// @Tags responsibilities
// @Accept json
// @Produce json
// @Param id path string true "Responsibility ID"
// @Param responsibility body models.ResponseRequest true "Responsibility data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/respons/update/{id} [put]
func UpdateResponsHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.ResponseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_, err := db.DB.Exec(
		"UPDATE responsibilities SET telegram_username = ?, name = ? WHERE id = ?",
		req.TelegramUsername, req.Name, id,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update responsibility"})
	}

	log.Printf("Responsibility ID: %s updated successfully", id)
	return c.JSON(fiber.Map{"message": "Responsibility updated"})
}

// @Summary Delete responsibility
// @Description Delete a responsibility
// @Tags responsibilities
// @Accept json
// @Produce json
// @Param id path string true "Responsibility ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/respons/delete/{id} [delete]
func DeleteResponsHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := db.DB.Exec(
		"DELETE FROM responsibilities WHERE id = ?",
		id,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete responsibility"})
	}

	log.Printf("Responsibility ID: %s deleted successfully", id)
	return c.JSON(fiber.Map{"message": "Responsibility deleted"})
}

// @Summary Get responsibility details
// @Description Get detailed information about a specific responsibility
// @Tags responsibilities
// @Accept json
// @Produce json
// @Param id path string true "Responsibility ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/respons/{id} [get]
func GetResponsDetailHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}
	var user models.ResponseRequest
	err = db.DB.QueryRow("SELECT id, IFNULL(telegram_username, '') as telegram_username, IFNULL(name, '') as name FROM responsibilities WHERE id = ?", id).Scan(&user.ID, &user.TelegramUsername, &user.Name)

	if err != nil {
		log.Printf("Error fetching program details: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Program not found"})
	}

	log.Printf("Getting program details Success for ID: %d", id)
	return c.JSON(fiber.Map{"success": true, "data": user})
}

// @Summary Get all users
// @Description Get all users with username and role
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/authEntry/users [get]
func GetAllUsersHandler(c *fiber.Ctx) error {
	rows, err := db.DB.Query("SELECT id, username, role FROM users WHERE deleted_at IS NULL")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()

	var users []models.UsernameResponse
	for rows.Next() {
		var user models.UsernameResponse
		if err := rows.Scan(&user.ID, &user.Username, &user.Role); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}
		users = append(users, user)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    users,
	})
}

// @Summary Get user details
// @Description Get detailed information about a specific user including username and role
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/authEntry/user/{id} [get]
func GetUserDetailHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var user models.UsernameResponse
	err = db.DB.QueryRow("SELECT id, username, role FROM users WHERE id = ? AND deleted_at IS NULL", id).Scan(&user.ID, &user.Username, &user.Role)

	if err != nil {
		log.Printf("Error fetching user details: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	log.Printf("Getting user details Success for ID: %d", id)
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
