package main

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestValidateName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid name", "John Doe", true},
		{"Single word", "John", false},
		{"Empty string", "", false},
		{"Three words", "John Doe Smith", false},
		{"Only spaces", "   ", false},
		{"Name with spaces at ends", " John Doe ", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateName(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid email", "test@example.com", true},
		{"Valid email with numbers", "test123@example.com", true},
		{"Valid email with subdomain", "test@mail.example.com", true},
		{"Invalid email - no @", "testexample.com", false},
		{"Invalid email - no domain", "test@", false},
		{"Invalid email - no local part", "@example.com", false},
		{"Invalid email - spaces", "test @example.com", false},
		{"Invalid email - invalid characters", "test@exa$mple.com", false},
		{"Empty string", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateEmail(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid password", "Password123", true},
		{"Valid password with special chars", "Password123!", true},
		{"Too short", "Pass1", false},
		{"No numbers", "Password", false},
		{"No lowercase", "PASSWORD123", false},
		{"No uppercase", "password123", false},
		{"Only numbers", "12345678", false},
		{"Empty string", "", false},
		{"Minimum valid", "Aa1bcdef", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidatePassword(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDocumentValidate(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid CPF", "11144477735", true},
		{"Valid CPF with dots and dash", "111.444.777-35", true},
		{"Valid CPF with spaces", "111 444 777 35", true},
		{"Invalid CPF - wrong digits", "12345678901", false},
		{"Invalid CPF - too short", "1234567890", false},
		{"Invalid CPF - too long", "123456789012", false},
		{"Invalid CPF - all same digits", "11111111111", false},
		{"Invalid CPF - all zeros", "00000000000", false},
		{"Empty string", "", false},
		{"Only letters", "abcdefghijk", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := Document{Digits: tc.input}
			result := doc.Validate()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDocumentClean(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"With dots and dash", "111.444.777-35", "11144477735"},
		{"With spaces", "111 444 777 35", "11144477735"},
		{"Mixed formatting", "111.444 777-35", "11144477735"},
		{"Only numbers", "11144477735", "11144477735"},
		{"With letters", "111a444b777c35", "11144477735"},
		{"Empty string", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := Document{Digits: tc.input}
			doc.clean()
			assert.Equal(t, tc.expected, doc.Digits)
		})
	}
}

func TestDocumentAllDigitsSame(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"All ones", "11111111111", true},
		{"All zeros", "00000000000", true},
		{"Mixed digits", "11144477735", false},
		{"Single digit", "1", true},
		{"Empty string", "", false},
		{"Two different digits", "12", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := Document{Digits: tc.input}
			result := doc.allDigitsSame()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDocumentCalculateDigit(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		factor   int
		expected int
	}{
		{"CPF first digit calculation", "111444777", 10, 3},
		{"CPF second digit calculation", "1114447773", 11, 5},
		{"Short input", "123", 10, 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := Document{Digits: tc.input}
			result := doc.calculateDigit(tc.factor)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsQuantityValid(t *testing.T) {
	testCases := []struct {
		name     string
		input    decimal.Decimal
		expected bool
	}{
		{"Positive quantity", decimal.NewFromInt(100), true},
		{"Zero quantity", decimal.Zero, true},
		{"Negative quantity", decimal.NewFromInt(-10), false},
		{"Large positive quantity", decimal.NewFromInt(999999), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isQuantityValid(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsValidUUID(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid UUID v4", "550e8400-e29b-41d4-a716-446655440000", true},
		{"Valid UUID v1", "6ba7b810-9dad-11d1-80b4-00c04fd430c8", true},
		{"Invalid UUID - too short", "550e8400-e29b-41d4-a716", false},
		{"Invalid UUID - wrong format", "550e8400-e29b-41d4-a716-44665544000g", false},
		{"Invalid UUID - missing dashes", "550e8400e29b41d4a716446655440000", false},
		{"Empty string", "", false},
		{"Random string", "not-a-uuid", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidUUID(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}