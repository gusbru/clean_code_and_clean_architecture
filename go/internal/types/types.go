package types

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Document string `json:"document"`
	Password string `json:"password"`
}

type User struct {
	AccountID uuid.UUID `json:"accountId"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Document  string    `json:"document"`
	Password  string    `json:"password"`
}

type Account struct {
	AccountID uuid.UUID `json:"accountId"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Document  string    `json:"document"`
	Assets    []Asset   `json:"assets"`
}

type Asset struct {
	AssetID  AssetId         `json:"assetId"`
	Quantity decimal.Decimal `json:"quantity"`
}

type DepositRequest struct {
	AccountID string          `json:"accountId"`
	AssetID   AssetId         `json:"assetId"`
	Quantity  decimal.Decimal `json:"quantity"`
}

type AssetId string

const (
	AssetIdBTC AssetId = "BTC"
	AssetIdUSD AssetId = "USD"
)

type WithdrawRequest struct {
	AccountID string          `json:"accountId"`
	AssetID   AssetId         `json:"assetId"`
	Quantity  decimal.Decimal `json:"quantity"`
}

func (a AssetId) String() string {
	switch a {
	case AssetIdBTC:
		return "BTC"
	case AssetIdUSD:
		return "USD"
	default:
		return "Unknown Asset"
	}
}

func (a AssetId) IsValid() bool {
	return a == AssetIdBTC || a == AssetIdUSD
}
