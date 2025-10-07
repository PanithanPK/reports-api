package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"reports-api/db"
	"reports-api/models"
	"reports-api/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// SessionData holds session information
type SessionData struct {
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

var sessions = map[string]SessionData{}

func generateSessionID() string {
	// Generate cryptographically secure random bytes
	bytes := make([]byte, 16) // 16 bytes = 128 bits of entropy
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("Error generating session ID: %v", err)
		// Fallback to timestamp-based ID (not ideal but better than predictable)
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(bytes)
}

func generateDummyToken() string {
	// Generate cryptographically secure random bytes for token
	bytes := make([]byte, 32) // 32 bytes = 256 bits of entropy
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("Error generating dummy token: %v", err)
		// Fallback to hex encoding of timestamp (not ideal but better than predictable)
		return hex.EncodeToString([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
	}
	return hex.EncodeToString(bytes)
}

// GetSessionData retrieves session data by session ID for middleware use
// Returns session data and validity status (checks expiration)
func GetSessionData(sessionID string) (SessionData, bool) {
	sessionData, exists := sessions[sessionID]
	if !exists {
		return SessionData{}, false
	}

	// Check if session has expired using Thailand timezone
	thailandTZ, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		// Fallback to UTC+7 fixed offset if timezone loading fails
		thailandTZ = time.FixedZone("ICT", 7*3600) // UTC+7
	}
	now := time.Now().In(thailandTZ)
	if now.After(sessionData.ExpiresAt) {
		// Remove expired session
		delete(sessions, sessionID)
		log.Printf("Session %s expired and removed", sessionID)
		return SessionData{}, false
	}

	return sessionData, true
}

// CleanupExpiredSessions removes all expired sessions from memory
func CleanupExpiredSessions() {
	// Use Thailand timezone for cleanup
	thailandTZ, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		// Fallback to UTC+7 fixed offset if timezone loading fails
		thailandTZ = time.FixedZone("ICT", 7*3600) // UTC+7
	}
	now := time.Now().In(thailandTZ)
	expiredCount := 0

	for sessionID, sessionData := range sessions {
		if now.After(sessionData.ExpiresAt) {
			delete(sessions, sessionID)
			expiredCount++
		}
	}

	if expiredCount > 0 {
		log.Printf("Cleaned up %d expired sessions", expiredCount)
	}
}

// IsSessionValid checks if a session exists and is not expired
func IsSessionValid(sessionID string) bool {
	_, valid := GetSessionData(sessionID)
	return valid
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
	// Load Thailand timezone (UTC+7)
	thailandTZ, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		// Fallback to UTC+7 fixed offset if timezone loading fails
		thailandTZ = time.FixedZone("ICT", 7*3600) // UTC+7
		log.Printf("Warning: Failed to load Asia/Bangkok timezone, using fixed UTC+7: %v", err)
	}
	now := time.Now().In(thailandTZ)
	sessionExpiry := now.Add(24 * time.Hour) // 24 hours expiration

	sessions[sessionID] = SessionData{
		Username:  username,
		Role:      role,
		CreatedAt: now,
		ExpiresAt: sessionExpiry,
	}

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
		Success: true,
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

		// Hash password for authentication (login verification)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
		}

		// Encrypt password for admin viewing (reversible)
		encryptedPassword, err := utils.EncryptPassword(req.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to encrypt password"})
		}

		_, err = db.DB.Exec(
			"INSERT INTO users (username, password, plain_password, role) VALUES (?, ?, ?, ?)",
			req.Username, string(hashedPassword), encryptedPassword, role,
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

	// Hash password for authentication (login verification)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Encrypt password for admin viewing (reversible)
	encryptedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to encrypt password"})
	}

	_, err = db.DB.Exec(
		"UPDATE users SET username = ?, password = ?, plain_password = ?, role = ?, updated_at=CURRENT_TIMESTAMP WHERE id = ?",
		req.Username, string(hashedPassword), encryptedPassword, req.Role, req.ID,
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
func GetResponsibilitiesHandler(c *fiber.Ctx) error {
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
func AddResponsibilityHandler(c *fiber.Ctx) error {
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
func UpdateResponsibilityHandler(c *fiber.Ctx) error {
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
func DeleteResponsibilityHandler(c *fiber.Ctx) error {
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
func GetResponsibilityDetailHandler(c *fiber.Ctx) error {
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
// @Description Get all users with username, decrypted password and role (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/authEntry/users [get]
func GetAllUsersHandler(c *fiber.Ctx) error {
	rows, err := db.DB.Query("SELECT id, username, IFNULL(plain_password, '') as plain_password, role FROM users WHERE deleted_at IS NULL")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()

	var users []models.UsernameResponse
	for rows.Next() {
		var user models.UsernameResponse
		var encryptedPassword string
		if err := rows.Scan(&user.ID, &user.Username, &encryptedPassword, &user.Role); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}

		// Decrypt password for admin viewing
		if encryptedPassword != "" {
			decryptedPassword, err := utils.DecryptPassword(encryptedPassword)
			if err != nil {
				log.Printf("Error decrypting password for user %s: %v", user.Username, err)
				user.PlainPassword = "[Decryption Error]"
			} else {
				user.PlainPassword = decryptedPassword
			}
		}

		users = append(users, user)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    users,
	})
}

// @Summary Get user details
// @Description Get detailed information about a specific user including username, decrypted password and role (Admin only)
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
	var encryptedPassword string
	err = db.DB.QueryRow("SELECT id, username, IFNULL(plain_password, '') as plain_password, role FROM users WHERE id = ? AND deleted_at IS NULL", id).Scan(&user.ID, &user.Username, &encryptedPassword, &user.Role)

	if err != nil {
		log.Printf("Error fetching user details: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Decrypt password for admin viewing
	var plainPassword string
	if encryptedPassword != "" {
		decryptedPassword, err := utils.DecryptPassword(encryptedPassword)
		if err != nil {
			log.Printf("Error decrypting password for user %s: %v", user.Username, err)
			plainPassword = "[Decryption Error]"
		} else {
			plainPassword = decryptedPassword
		}
	}

	log.Printf("Getting user details Success for ID: %d", id)
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":             user.ID,
			"username":       user.Username,
			"plain_password": plainPassword,
			"role":           user.Role,
		},
	})
}
