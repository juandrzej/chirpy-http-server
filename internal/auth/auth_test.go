package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"

	hash1, err := HashPassword(password1)
	if err != nil {
		t.Fatal(err)
	}

	hash2, err := HashPassword(password2)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "password does not match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() = %v, want %v", match, tt.matchPassword)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()

	validToken, err := MakeJWT(userID, "secret", time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	expiredToken, err := MakeJWT(userID, "secret", -time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "expired token",
			tokenString: expiredToken,
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
