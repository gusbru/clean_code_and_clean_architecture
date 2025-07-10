package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignupEndpoint(t *testing.T) {
    // Given
	input := map[string]string{
		"name": "Gustavo B",
	}
	inputJson, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Input JSON:", string(inputJson))

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
	assert.Equal(t, "User created successfully", response["message"])
}