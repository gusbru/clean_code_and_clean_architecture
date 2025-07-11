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

func CreateValidAccount() (string, error) {
	inputNewAccount := map[string]string{
		"name":     "Gustavo B",
		"email":    fmt.Sprintf("gustavo-%d@example.com", time.Now().UnixNano()),
		"document": "11144477735",
		"password": "SecurePassword1234",
	}
	inputNewAccountJson, err := json.Marshal(inputNewAccount)
	if err != nil {
		return "", err
	}
	resp, err := http.Post("http://app:3000/signup", "application/json", bytes.NewBuffer(inputNewAccountJson))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var responseNewAccount map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseNewAccount)
	if err != nil {
		return "", err
	}
	newAccountID := responseNewAccount["accountId"]
	return newAccountID, nil
}

func TestDepositEndpoint(t *testing.T) {
	// Given
	newAccountID, err := CreateValidAccount()
	if err != nil {
		t.Fatal(err)
	}
	input := map[string]string{
		"accountId": newAccountID,
		"assetId":   "BTC",
		"quantity":  "10",
	}
	inputJson, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	// When
	input = map[string]string{
		"accountId": newAccountID,
		"assetId":   "BTC",
		"quantity":  "10",
	}
	inputJson, err = json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	// When
	resp, err := http.Post("http://app:3000/deposit", "application/json", bytes.NewBuffer(inputJson))
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
	assert.Equal(t, "Deposit completed", response["message"], "Expected message to indicate deposit completed")
}

func TestInvalidDepositRequest(t *testing.T) {
	accountID, err := CreateValidAccount()
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name         string
		input        map[string]string
		expectedCode int
		expectedErr  string
	}{
		{
			name: "Missing accountId",
			input: map[string]string{
				"assetId":  "BTC",
				"quantity": "10",
			},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "accountId is required",
		},
		{
			name: "Missing accountId",
			input: map[string]string{
				"accountId": "",
				"assetId":  "BTC",
				"quantity": "10",
			},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "accountId is required",
		},
		{
			name: "Missing assetId",
			input: map[string]string{
				"accountId": accountID,
				"quantity":  "10",
			},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "assetId is required and must be valid",
		},
		{
			name: "Missing quantity",
			input: map[string]string{
				"accountId": accountID,
				"assetId":   "BTC",
			},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "quantity is required and must be a valid positive integer",
		},
		{
			name: "Invalid Asset ID",
			input: map[string]string{
				"accountId": accountID,
				"assetId":   "INVALID",
				"quantity":  "10",
			},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "assetId is required and must be valid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputJson, err := json.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := http.Post("http://app:3000/deposit", "application/json", bytes.NewBuffer(inputJson))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			assert.Equal(t, tc.expectedCode, resp.StatusCode)
			var response map[string]string
			err = json.NewDecoder(resp.Body).Decode(&response)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.expectedErr, response["error"], "Expected error message to match")
		})
	}
}