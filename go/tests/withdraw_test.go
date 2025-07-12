package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gusbru/clean_code_and_clean_architecture/internal/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestWithdrawEndpoint(t *testing.T) {
	// Given
	newAccountID, err := CreateValidAccount(CreateAccountOptions{
		AddAsset: true,
		Asset:    types.Asset{AssetID: "BTC", Quantity: decimal.NewFromInt(10)},
	})
	if err != nil {
		t.Fatal(err)
	}

	// When
	inputWithdraw := map[string]string{
		"accountId": newAccountID,
		"assetId":   "BTC",
		"quantity":  "5",
	}
	inputWithdrawJson, err := json.Marshal(inputWithdraw)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post("http://app:3000/withdraw", "application/json", bytes.NewBuffer(inputWithdrawJson))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200 for withdraw")

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
	assets := accountResponse["assets"].([]interface{})
	assert.Len(t, assets, 1, "Expected one asset in account after withdrawal")
	asset := assets[0].(map[string]interface{})
	assert.Equal(t, "BTC", asset["assetId"], "Expected assetId to be BTC")
	assert.Equal(t, "5", asset["quantity"], "Expected quantity to be 5 after withdrawal")
}

func TestWithdrawInvalidCases(t *testing.T) {
	newAccountID, err := CreateValidAccount(CreateAccountOptions{
		AddAsset: true,
		Asset:    types.Asset{AssetID: "BTC", Quantity: decimal.NewFromInt(10)},
	})
	if err != nil {
		t.Fatal(err)
	}

	// When
	testCases := []struct {
		name         string
		input        map[string]string
		expectedCode int
	}{
		{
			name: "Withdraw more than available",
			input: map[string]string{
				"accountId": newAccountID,
				"assetId":   "BTC",
				"quantity":  "15",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Withdraw negative quantity",
			input: map[string]string{
				"accountId": newAccountID,
				"assetId":   "BTC",
				"quantity":  "-5",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Withdraw with invalid asset",
			input: map[string]string{
				"accountId": newAccountID,
				"assetId":   "INVALID",
				"quantity":  "5",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Withdraw with invalid accountId",
			input: map[string]string{
				"accountId": "INVALID",
				"assetId":   "BTC",
				"quantity":  "5",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Withdraw with empty accountId",
			input: map[string]string{
				"accountId": "",
				"assetId":   "BTC",
				"quantity":  "5",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Withdraw with empty assetId",
			input: map[string]string{
				"accountId": newAccountID,
				"assetId":   "",
				"quantity":  "5",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Withdraw with empty quantity",
			input: map[string]string{
				"accountId": newAccountID,
				"assetId":   "BTC",
				"quantity":  "",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Withdraw with non-numeric quantity",
			input: map[string]string{
				"accountId": newAccountID,
				"assetId":   "BTC",
				"quantity":  "invalid",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Withdraw with zero quantity",
			input: map[string]string{
				"accountId": newAccountID,
				"assetId":   "BTC",
				"quantity":  "0",
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputJson, err := json.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.Post("http://app:3000/withdraw", "application/json", bytes.NewBuffer(inputJson))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedCode, resp.StatusCode, "Expected status code for "+tc.name)
		})
	}
}
