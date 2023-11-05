package config

import (
	"testing"
)

func TestCrypto(t *testing.T) {
	// Arrange
	secretKey, err := generateSecretKey()
	if err != nil {
		t.Fatalf("arrange phase - %s", err)
	}

	expectedPlaintext := "simple-plaintext"

	// Act
	ciphertext, err := encrypt(expectedPlaintext, secretKey)
	if err != nil {
		t.Fatalf("act phase - %s", err)
	}

	// Assert
	actualPlaintext, err := decrypt(ciphertext, secretKey)
	if err != nil {
		t.Fatal(err)
	}

	if expectedPlaintext != actualPlaintext {
		t.Fatalf("expected: %s, actual: %s", expectedPlaintext, actualPlaintext)
	}
}
