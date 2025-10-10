package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

// getEncryptionKey retrieves the encryption key from environment variable
func getEncryptionKey() []byte {
	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		// Default key for development (should be changed in production)
		key = "myverystrongpasswordo32bitlength"
	}
	return []byte(key)
}

// EncryptPassword encrypts a password using AES-GCM
func EncryptPassword(plaintext string) (string, error) {
	key := getEncryptionKey()

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64 for storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPassword decrypts a password using AES-GCM
func DecryptPassword(ciphertext string) (string, error) {
	key := getEncryptionKey()

	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Check minimum length
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]

	// Decrypt the data
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// HashPasswordForAuth creates a bcrypt hash for authentication (login)
func HashPasswordForAuth(password string) (string, error) {
	// This is still used for login verification
	// We keep bcrypt for authentication security
	return "", nil // Will implement if needed for dual approach
}
