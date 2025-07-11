package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignupEndpoint(t *testing.T) {
	// Given
	input := map[string]string{
		"name":     "Gustavo B",
		"email":    fmt.Sprintf("gustavo-%d@example.com", time.Now().UnixNano()),
		"document": "11144477735",
		"password": "SecurePassword1234",
	}
	inputJson, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	// When
	resp, err := http.Post("http://app:3000/signup", "application/json", bytes.NewBuffer(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	newAccountID := response["accountId"]
	assert.NotEmpty(t, newAccountID, "Expected account_id to be present in response")
	resp, err = http.Get(fmt.Sprintf("http://app:3000/accounts/%s", newAccountID))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200 for account check")
	var accountResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&accountResponse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, newAccountID, accountResponse["accountId"], "Expected accountId in response to match the one created")
}

func TestInvalidNameSignup(t *testing.T) {
	testCases := []struct {
		name         string
		inputName    string
		expectedCode int
		expectedErr  string
	}{
		{
			name:         "Single word name",
			inputName:    "Gustavo",
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Invalid name",
		},
		{
			name:         "Empty name",
			inputName:    "",
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Invalid name",
		},
		{
			name:         "Only spaces",
			inputName:    "   ",
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Invalid name",
		},
		{
			name:         "Name with more than two words",
			inputName:    "Gustavo Bruno B",
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Invalid name",
		},
		{
			name:         "Valid name",
			inputName:    "Gustavo B",
			expectedCode: http.StatusOK,
			expectedErr:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			input := map[string]string{
				"name":     tc.inputName,
				"email":    fmt.Sprintf("valid-%d@example.com", time.Now().UnixNano()),
				"document": "11144477735",
				"password": "SecurePassword1234",
			}
			inputJson, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			// When
			resp, err := http.Post("http://app:3000/signup", "application/json", bytes.NewBuffer(inputJson))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			// Then
			assert.Equal(t, tc.expectedCode, resp.StatusCode)

			if tc.expectedCode == http.StatusBadRequest {
				var response map[string]string
				err = json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.expectedErr, response["error"])
			}
		})
	}
}

func TestInvalidEmailSignup(t *testing.T) {
	testCases := []struct {
		name         string
		inputEmail   string
		expectedCode int
		expectedErr  string
	}{
		{
			name:         "Invalid email format",
			inputEmail:   "invalid-email",
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Invalid email",
		},
		{
			name:         "Empty email",
			inputEmail:   "",
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Invalid email",
		},
		{
			name:         "Email with spaces",
			inputEmail:   "test @example.com",
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Invalid email",
		},
		{
			name:         "Email without domain",
			inputEmail:   "test@",
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Invalid email",
		},
		{
			name:         "Email with invalid characters",
			inputEmail:   "test@exa$mple.com",
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Invalid email",
		},
		{
			name:         "Valid email",
			inputEmail:   fmt.Sprintf("test-%d@example.com", time.Now().UnixNano()),
			expectedCode: http.StatusOK,
			expectedErr:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			input := map[string]string{
				"name":     "Gustavo B",
				"email":    tc.inputEmail,
				"document": "11144477735",
				"password": "SecurePassword1234",
			}
			inputJson, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			// When
			resp, err := http.Post("http://app:3000/signup", "application/json", bytes.NewBuffer(inputJson))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			// Then
			assert.Equal(t, tc.expectedCode, resp.StatusCode)

			if tc.expectedCode == http.StatusBadRequest {
				var response map[string]string
				err = json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.expectedErr, response["error"])
			}
		})
	}
}

func TestDuplicatedEmail(t *testing.T) {
	// First signup
	input := map[string]string{
		"name":     "Gustavo B",
		"email":    fmt.Sprintf("gustavo-%d@example.com", time.Now().UnixNano()),
		"document": "11144477735",
		"password": "SecurePassword1234",
	}
	inputJson, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	// When
	resp, err := http.Post("http://app:3000/signup", "application/json", bytes.NewBuffer(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Second signup with the same email
	resp, err = http.Post("http://app:3000/signup", "application/json", bytes.NewBuffer(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Email already exists", response["error"])
}

func TestInvalidPasswordSignup(t *testing.T) {
	testCases := []struct {
		name          string
		inputPassword string
		expectedCode  int
		expectedErr   string
	}{
		{
			name:          "Password too short",
			inputPassword: "short",
			expectedCode:  http.StatusBadRequest,
			expectedErr:   "Invalid password",
		},
		{
			name:          "Password without numbers",
			inputPassword: "NoNumbersHere",
			expectedCode:  http.StatusBadRequest,
			expectedErr:   "Invalid password",
		},
		{
			name:          "Password without uppercase letters",
			inputPassword: "nouppercase123",
			expectedCode:  http.StatusBadRequest,
			expectedErr:   "Invalid password",
		},
		{
			name:          "Valid password",
			inputPassword: "Valid1234",
			expectedCode:  http.StatusOK,
			expectedErr:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			input := map[string]string{
				"name":     "Gustavo B",
				"email":    fmt.Sprintf("valid-%d@example.com", time.Now().UnixNano()),
				"document": "11144477735",
				"password": tc.inputPassword,
			}
			inputJson, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			// When
			resp, err := http.Post("http://app:3000/signup", "application/json", bytes.NewBuffer(inputJson))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			// Then
			assert.Equal(t, tc.expectedCode, resp.StatusCode)

			if tc.expectedCode == http.StatusBadRequest {
				var response map[string]string
				err = json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.expectedErr, response["error"])
			}
		})
	}
}

func TestInvalidDocument(t *testing.T) {
	testCases := []struct {
		name          string
		inputDocument string
		expectedCode  int
		expectedErr   string
	}{
		{
			name:          "Invalid document format",
			inputDocument: "12345678901",
			expectedCode:  http.StatusBadRequest,
			expectedErr:   "Invalid document",
		},
		{
			name:          "Empty document",
			inputDocument: "",
			expectedCode:  http.StatusBadRequest,
			expectedErr:   "Invalid document",
		},
		{
			name:          "Document with spaces",
			inputDocument: "111 325 777 35",
			expectedCode:  http.StatusBadRequest,
			expectedErr:   "Invalid document",
		},
		{
			name:          "Valid document",
			inputDocument: "11144477735",
			expectedCode:  http.StatusOK,
			expectedErr:   "",
		},
		{
			name:          "Valid document",
			inputDocument: "111 444 777 35",
			expectedCode:  http.StatusOK,
			expectedErr:   "",
		},
		{
			name:          "Valid document",
			inputDocument: "111.444.777-35",
			expectedCode:  http.StatusOK,
			expectedErr:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			input := map[string]string{
				"name":     "Gustavo B",
				"email":    fmt.Sprintf("gustavo-%d@example.com", time.Now().UnixNano()),
				"document": tc.inputDocument,
				"password": "SecurePassword1234",
			}
			inputJson, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			// When
			resp, err := http.Post("http://app:3000/signup", "application/json", bytes.NewBuffer(inputJson))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			// Then
			assert.Equal(t, tc.expectedCode, resp.StatusCode)

			if tc.expectedCode == http.StatusBadRequest {
				var response map[string]string
				err = json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.expectedErr, response["error"])
			}
		})
	}
}
