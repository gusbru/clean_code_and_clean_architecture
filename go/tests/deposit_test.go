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
	// When
	inputDeposit := map[string]string{
		"accountId": newAccountID,
		"assetId":   "BTC",
		"quantity":  "10",
	}
	inputDepositJson, err := json.Marshal(inputDeposit)
	if err != nil {
		t.Fatal(err)
	}
	// When
	resp, err := http.Post("http://app:3000/deposit", "application/json", bytes.NewBuffer(inputDepositJson))
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
	assert.Equal(t, newAccountID, accountResponse["accountId"], "Expected accountId to match")
	assert.NotEmpty(t, accountResponse["assets"], "Expected assets to be present")
	assets := accountResponse["assets"].([]interface{})
	assert.Len(t, assets, 1, "Expected one asset in account")
	asset := assets[0].(map[string]interface{})
	assert.Equal(t, "BTC", asset["assetId"], "Expected assetId to be BTC")
	assert.Equal(t, "10", asset["quantity"], "Expected quantity to be 10")
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
				"assetId":   "BTC",
				"quantity":  "10",
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
			name: "Negative quantity",
			input: map[string]string{
				"accountId": accountID,
				"assetId":   "BTC",
				"quantity":  "-10",
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

func TestGetAccountWithNoAssets(t *testing.T) {
	// Given
	newAccountID, err := CreateValidAccount()
	if err != nil {
		t.Fatal(err)
	}
	// When
	resp, err := http.Get(fmt.Sprintf("http://app:3000/accounts/%s", newAccountID))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Print("Response status code: ", resp.StatusCode, "\n")
	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, newAccountID, response["accountId"], "Expected accountId to match")
	assert.Empty(t, response["assets"], "Expected no assets for new account")
}
