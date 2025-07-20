package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Account represents the account data structure
type Account struct {
	AccountID string `json:"account_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Document  string `json:"document"`
	Password  string `json:"password"`
}

// IAccountDAO defines the interface for account data access operations
type IAccountDAO interface {
	Save(account *Account) error
	GetByID(accountID string) (*Account, error)
	GetByEmail(email string) (*Account, error)
}

// AccountDAODatabase implements IAccountDAO using PostgreSQL database
type AccountDAODatabase struct{}

func NewAccountDAODatabase() *AccountDAODatabase {
	return &AccountDAODatabase{}
}

func (dao *AccountDAODatabase) getDB() (*sql.DB, error) {
	connStr := "host=db port=5432 user=postgres password=postgres dbname=app sslmode=disable"
	return sql.Open("postgres", connStr)
}

func (dao *AccountDAODatabase) Save(account *Account) error {
	db, err := dao.getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "INSERT INTO ccca.account (account_id, name, email, document, password) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.Exec(query, account.AccountID, account.Name, account.Email, account.Document, account.Password)
	return err
}

func (dao *AccountDAODatabase) GetByID(accountID string) (*Account, error) {
	db, err := dao.getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT account_id, name, email, document, password FROM ccca.account WHERE account_id = $1"
	row := db.QueryRow(query, accountID)

	account := &Account{}
	err = row.Scan(&account.AccountID, &account.Name, &account.Email, &account.Document, &account.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return account, nil
}

func (dao *AccountDAODatabase) GetByEmail(email string) (*Account, error) {
	db, err := dao.getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT account_id, name, email, document, password FROM ccca.account WHERE email = $1"
	row := db.QueryRow(query, email)

	account := &Account{}
	err = row.Scan(&account.AccountID, &account.Name, &account.Email, &account.Document, &account.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return account, nil
}

// AccountDAOMemory implements IAccountDAO using in-memory storage
type AccountDAOMemory struct {
	accounts   map[string]*Account
	emailIndex map[string]string
}

func NewAccountDAOMemory() *AccountDAOMemory {
	return &AccountDAOMemory{
		accounts:   make(map[string]*Account),
		emailIndex: make(map[string]string),
	}
}

func (dao *AccountDAOMemory) Save(account *Account) error {
	dao.accounts[account.AccountID] = account
	dao.emailIndex[account.Email] = account.AccountID
	return nil
}

func (dao *AccountDAOMemory) GetByID(accountID string) (*Account, error) {
	account, exists := dao.accounts[accountID]
	if !exists {
		return nil, nil
	}
	return account, nil
}

func (dao *AccountDAOMemory) GetByEmail(email string) (*Account, error) {
	accountID, exists := dao.emailIndex[email]
	if !exists {
		return nil, nil
	}

	account, exists := dao.accounts[accountID]
	if !exists {
		return nil, nil
	}

	return account, nil
}
