package otp

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

var OTPManage *OTPManager

type OTPPurpose string

const (
	WithdrawalOTP        OTPPurpose = "withdrawal"
	EmailVerificationOTP OTPPurpose = "email_verification"
)

// OTP struct represents an OTP with its associated email and expiration time
type OTP struct {
	Email     string
	Token     string
	ExpiresAt time.Time
	Purpose   OTPPurpose
}

// OTPManager manages the generation and verification of OTPs
type OTPManager struct {
	Otps map[string]OTP
	mu   sync.Mutex
}

// NewOTPManager creates a new OTPManager instance
func NewOTPManager() *OTPManager {
	OTPManage = &OTPManager{
		Otps: make(map[string]OTP),
	}
	return OTPManage
}

// GenerateOTP generates a unique OTP for the given email with a specified expiration time
func (m *OTPManager) GenerateOTP(email string, expiration time.Duration) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Check if the OTP exists and has not expired
	if otp, ok := m.Otps[email]; ok && time.Now().Before(otp.ExpiresAt) {
		return otp.Token, nil
	}
	const otpLength = 5
	// Create a character set consisting of alphanumeric capital letters
	charSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charSetLength := len(charSet)

	// Generate a random token of fixed length
	tokenBytes := make([]byte, otpLength)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Map the random bytes to the character set
	for i := 0; i < otpLength; i++ {
		tokenBytes[i] = charSet[tokenBytes[i]%byte(charSetLength)]
	}

	token := string(tokenBytes)

	// Calculate the expiration time
	expiresAt := time.Now().Add(expiration)
	// Store the OTP in the manager
	m.Otps[email] = OTP{
		Email:     email,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return token, nil
}

// VerifyOTP verifies the provided OTP for the given email
func (m *OTPManager) VerifyOTP(email, token string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Retrieve the stored OTP for the email
	otp, ok := m.Otps[email]
	if !ok {
		return false
	}
	// Check if the OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		return false
	}

	// delete the OTP
	if otp.Token == token {
		delete(m.Otps, email)
		return true
	}

	return false
}
