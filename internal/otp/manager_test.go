package otp

import (
	"testing"
	"time"
)

func TestOTPManager_GenerateOTP(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		expiration time.Duration
	}{
		{
			name:       "Valid OTP",
			email:      "test@example.com",
			expiration: time.Minute,
		},
		{
			name:       "Expired OTP",
			email:      "expired@example.com",
			expiration: time.Nanosecond,
		},
	}

	otpManager := NewOTPManager()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			otpToken, err := otpManager.GenerateOTP(test.email, test.expiration)
			if err != nil {
				t.Errorf("Failed to generate OTP: %v", err)
			}

			valid := otpManager.VerifyOTP(test.email, otpToken)

			if test.name == "Valid OTP" && !valid {
				t.Errorf("OTP verification failed for a valid OTP")
			}

			if test.name == "Expired OTP" && valid {
				t.Errorf("OTP verification succeeded for an expired OTP")
			}
		})
	}
}
func TestGenerateOTP(t *testing.T) {
	t.Log("Generating random OTP")

	otp := NewOTPManager()

	code, _ := otp.GenerateOTP("john@gmail.com", time.Millisecond*1000)

	good := otp.VerifyOTP("john@gmail.com", code)

	if !good {
		t.Errorf("Invalid OTP %s", code)
	}
}

func TestOTPManager_VerifyOTP_Invalid(t *testing.T) {
	otpManager := NewOTPManager()

	email := "test@example.com"
	expiration := time.Minute

	// Generate OTP
	_, err := otpManager.GenerateOTP(email, expiration)
	if err != nil {
		t.Errorf("Failed to generate OTP: %v", err)
	}

	// Verify OTP with an invalid token
	valid := otpManager.VerifyOTP(email, "invalid-token")
	if valid {
		t.Errorf("OTP verification succeeded for an invalid OTP token")
	}
}

func TestOTPManager_VerifyOTP_Expired(t *testing.T) {
	otpManager := NewOTPManager()

	email := "test@example.com"
	expiration := time.Second // Set a very short expiration time for testing

	// Generate OTP
	otpToken, err := otpManager.GenerateOTP(email, expiration)
	if err != nil {
		t.Errorf("Failed to generate OTP: %v", err)
	}

	// Wait for OTP to expire
	time.Sleep(time.Second)

	// Verify OTP
	valid := otpManager.VerifyOTP(email, otpToken)
	if valid {
		t.Errorf("OTP verification succeeded for an expired OTP")
	}
}
