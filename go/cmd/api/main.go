package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gusbru/clean_code_and_clean_architecture/internal/types"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() *Database {
	host := "postgres"
	port := "5432"
	user := "postgres"
	password := "postgres"
	dbname := "app"

	logrus.WithFields(logrus.Fields{
		"host":   host,
		"port":   port,
		"dbname": dbname,
	}).Info("Connecting to database")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}
	logrus.Info("Database connection established successfully")
	return &Database{DB: db}
}

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		logrus.WithFields(logrus.Fields{
			"method": c.Method(),
			"path":   c.Path(),
			"ip":     c.IP(),
		}).Info("Incoming request")
		err := c.Next()
		status := c.Response().StatusCode()
		fields := logrus.Fields{
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     status,
			"latency_ms": time.Since(start).Milliseconds(),
			"ip":         c.IP(),
			"size":       len(c.Response().Body()),
		}
		switch {
		case status >= 500:
			logrus.WithFields(fields).Error("Server error")
		case status >= 400:
			logrus.WithFields(fields).Warn("Client error")
		default:
			logrus.WithFields(fields).Info("Request completed")
		}
		return err
	}
}

func ValidateName(name string) bool {
	return len(strings.Split(name, " ")) == 2
}

func ValidateEmail(email string) bool {
	// Regex for basic email validation
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func CheckDuplicateEmail(db *Database, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM ccca.account WHERE email = $1)`
	err := db.DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func ValidatePassword(password string) bool {
	return len(password) >= 8 &&
		strings.ContainsAny(password, "0123456789") &&
		strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") &&
		strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func ValidateAccountExists(db *Database, accountID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM ccca.account WHERE account_id = $1)`
	err := db.DB.QueryRow(query, accountID).Scan(&exists)
	logrus.WithFields(logrus.Fields{
		"accountId": accountID,
		"exists":    exists,
	}).Info("Checking account existence")
	if err != nil {
		logrus.WithError(err).Error("Error checking account existence")
		return false, err
	}
	if !exists {
		logrus.Warn("Account does not exist", logrus.Fields{"accountId": accountID})
		return false, fmt.Errorf("account does not exist")
	}
	return true, nil
}

type Document struct {
	Digits string
}

func (d *Document) Validate() bool {
	if d.Digits == "" {
		return false
	}
	d.clean()
	if len(d.Digits) != 11 {
		return false
	}
	if d.allDigitsSame() {
		return false
	}
	return d.extractDigits() == fmt.Sprintf("%d%d", d.calculateDigit(10), d.calculateDigit(11))
}

func (d *Document) clean() {
	// Remove non-numeric characters
	re := regexp.MustCompile(`\D`)
	d.Digits = re.ReplaceAllString(d.Digits, "")
}

func (d *Document) allDigitsSame() bool {
	if len(d.Digits) == 0 {
		return false
	}
	firstDigit := d.Digits[0]
	for i := 1; i < len(d.Digits); i++ {
		if d.Digits[i] != firstDigit {
			return false
		}
	}
	return true
}

func (d *Document) calculateDigit(factor int) int {
	sum := 0
	digitsToProcess := factor - 1
	for i := 0; i < digitsToProcess && i < len(d.Digits); i++ {
		digit := int(d.Digits[i] - '0')
		sum += digit * (factor - i)
	}
	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}

func (d *Document) extractDigits() string {
	return d.Digits[9:]
}

func isQuantityValid(quantity decimal.Decimal) bool {
	return quantity.GreaterThanOrEqual(decimal.Zero)
}

func isDepositValid(depositRequest types.DepositRequest) (bool, error) {
	if depositRequest.AccountID == "" {
		return false, fmt.Errorf("accountId is required")
	}
	if depositRequest.AssetID == "" || !depositRequest.AssetID.IsValid() {
		return false, fmt.Errorf("assetId is required and must be valid")
	}
	if depositRequest.Quantity.IsZero() || !isQuantityValid(depositRequest.Quantity) {
		return false, fmt.Errorf("quantity is required and must be a valid positive integer")
	}
	return true, nil
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func validateSignupRequest(req types.SignupRequest, db *Database) (bool, error) {
	if !ValidateName(req.Name) {
		return false, fmt.Errorf("Invalid name")
	}
	if !ValidateEmail(req.Email) {
		return false, fmt.Errorf("Invalid email")
	}
	emailExists, err := CheckDuplicateEmail(db, req.Email)
	if err != nil {
		logrus.WithError(err).Error("Error checking duplicate email")
		return false, fmt.Errorf("Failed to check email")
	}
	if emailExists {
		return false, fmt.Errorf("Email already exists")
	}
	if !ValidatePassword(req.Password) {
		return false, fmt.Errorf("Invalid password")
	}
	document := Document{Digits: req.Document}
	if !document.Validate() {
		return false, fmt.Errorf("Invalid document")
	}
	return true, nil
}

func handleSignup(c *fiber.Ctx, db *Database) error {
	var req types.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		logrus.WithError(err).Error("Failed to parse request body")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": "Invalid request body"})
	}
	logrus.WithFields(logrus.Fields{
		"email": req.Email,
		"name":  req.Name,
	}).Info("Processing signup")
	if valid, err := validateSignupRequest(req, db); !valid {
		logrus.WithError(err).Error("Invalid signup request")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	document := Document{Digits: req.Document}
	user := types.User{
		AccountID: uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		Document:  document.Digits,
		Password:  req.Password,
	}
	logrus.WithFields(logrus.Fields{
		"accountId": user.AccountID,
		"email":     user.Email,
	}).Info("Creating new account")
	query := `INSERT INTO ccca.account (account_id, name, email, document, password) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.DB.Exec(query, user.AccountID, user.Name, user.Email, user.Document, user.Password)
	if err != nil {
		logrus.WithError(err).Error("Error inserting account")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": "Failed to create account"})
	}
	logrus.WithField("accountId", user.AccountID).Info("Account created successfully")
	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{"accountId": user.AccountID})
}

func handleGetAccount(c *fiber.Ctx, db *Database) error {
	accountID := c.Params("accountId")
	if !isValidUUID(accountID) {
		logrus.Warn("Invalid account ID format", logrus.Fields{"accountId": accountID})
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": "Invalid account ID format"})
	}
	query := `SELECT account_id, name, email, document FROM ccca.account WHERE account_id = $1`
	var account types.Account
	err := db.DB.QueryRow(query, accountID).Scan(&account.AccountID, &account.Name, &account.Email, &account.Document)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{"error": "Account not found"})
		}
		logrus.WithError(err).Error("Error querying account")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": "Failed to retrieve account"})
	}
	account.Assets = []types.Asset{}
	assetQuery := `SELECT asset_id, quantity FROM ccca.account_asset WHERE account_id = $1`
	rows, err := db.DB.Query(assetQuery, account.AccountID)
	if err != nil {
		logrus.WithError(err).Error("Error querying account assets")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": "Failed to retrieve account assets"})
	}
	defer rows.Close()
	for rows.Next() {
		var asset types.Asset
		if err := rows.Scan(&asset.AssetID, &asset.Quantity); err != nil {
			logrus.WithError(err).Error("Error scanning asset")
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"error": "Failed to retrieve account assets"})
		}
		account.Assets = append(account.Assets, asset)
	}
	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error("Error iterating over account assets")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": "Failed to retrieve account assets"})
	}
	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"accountId": account.AccountID,
		"name":      account.Name,
		"email":     account.Email,
		"document":  account.Document,
		"assets":    account.Assets,
	})
}

func handleDeposit(c *fiber.Ctx, db *Database) error {
	var depositRequest types.DepositRequest
	if err := c.BodyParser(&depositRequest); err != nil {
		logrus.WithError(err).Error("Failed to parse deposit request body")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": "Invalid request body"})
	}
	logrus.WithFields(logrus.Fields{
		"accountId": depositRequest.AccountID,
		"assetId":   depositRequest.AssetID,
		"quantity":  depositRequest.Quantity,
	}).Info("Processing deposit")
	if valid, err := isDepositValid(depositRequest); !valid {
		logrus.WithError(err).Error("Invalid deposit request")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	if exists, err := ValidateAccountExists(db, depositRequest.AccountID); !exists || err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": "Account does not exist"})
	}
	query := `INSERT INTO ccca.account_asset (account_id, asset_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (account_id, asset_id) DO UPDATE SET quantity = ccca.account_asset.quantity + EXCLUDED.quantity`
	_, err := db.DB.Exec(query, depositRequest.AccountID, depositRequest.AssetID, depositRequest.Quantity)
	if err != nil {
		logrus.WithError(err).Error("Error inserting deposit")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": "Failed to process deposit"})
	}
	logrus.WithFields(logrus.Fields{
		"accountId": depositRequest.AccountID,
		"assetId":   depositRequest.AssetID,
		"quantity":  depositRequest.Quantity,
	}).Info("Deposit processed successfully")
	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "Deposit completed",
	})
}

func handleWithdraw(c *fiber.Ctx, db *Database) error {
	var withdrawRequest types.WithdrawRequest
	if err := c.BodyParser(&withdrawRequest); err != nil {
		logrus.WithError(err).Error("Failed to parse withdraw request body")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": "Invalid request body"})
	}
	if !withdrawRequest.AssetID.IsValid() {
		logrus.WithFields(logrus.Fields{
			"accountId": withdrawRequest.AccountID,
			"quantity":  withdrawRequest.Quantity,
			"assetId":   withdrawRequest.AssetID,
		}).Warn("Invalid asset ID for withdrawal")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": "Invalid asset ID"})
	}
	if !isQuantityValid(withdrawRequest.Quantity) {
		logrus.WithFields(logrus.Fields{
			"accountId": withdrawRequest.AccountID,
			"quantity":  withdrawRequest.Quantity,
			"assetId":   withdrawRequest.AssetID,
		}).Warn("Invalid quantity for withdrawal")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": "Invalid quantity"})
	}
	if exists, err := ValidateAccountExists(db, withdrawRequest.AccountID); !exists || err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": "Account does not exist"})
	}
	// Get asset details
	query := `SELECT asset_id, quantity FROM ccca.account_asset WHERE account_id = $1 AND asset_id = $2`
	var asset types.Asset
	err := db.DB.QueryRow(query, withdrawRequest.AccountID, withdrawRequest.AssetID).Scan(&asset.AssetID, &asset.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithFields(logrus.Fields{
				"accountId": withdrawRequest.AccountID,
				"assetId":   withdrawRequest.AssetID,
			}).Warn("Asset not found for withdrawal")
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{"error": "Asset not found"})
		}
		logrus.WithError(err).Error("Error querying asset for withdrawal")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": "Failed to retrieve asset"})
	}

	if asset.Quantity.LessThan(withdrawRequest.Quantity) {
		logrus.WithFields(logrus.Fields{
			"accountId": withdrawRequest.AccountID,
			"assetId":   withdrawRequest.AssetID,
			"quantity":  withdrawRequest.Quantity,
			"available": asset.Quantity,
		}).Warn("Insufficient asset quantity for withdrawal")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": "Insufficient asset quantity"})
	}
	// Update asset quantity
	updateQuery := `UPDATE ccca.account_asset SET quantity = quantity - $1 WHERE account_id = $2 AND asset_id = $3`
	_, err = db.DB.Exec(updateQuery, withdrawRequest.Quantity, withdrawRequest.AccountID, withdrawRequest.AssetID)
	if err != nil {
		logrus.WithError(err).Error("Error updating asset quantity for withdrawal")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": "Failed to process withdrawal"})
	}
	logrus.WithFields(logrus.Fields{
		"accountId": withdrawRequest.AccountID,
		"assetId":   withdrawRequest.AssetID,
		"quantity":  withdrawRequest.Quantity,
	}).Info("Withdrawal processed successfully")
	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{})
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("Starting application initialization")
	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})
	app.Use(LoggerMiddleware())
	db := NewDatabase()
	defer db.DB.Close()
	logrus.Info("Application started")

	app.Post("/signup", func(c *fiber.Ctx) error {
		return handleSignup(c, db)
	})

	app.Get("/accounts/:accountId", func(c *fiber.Ctx) error {
		return handleGetAccount(c, db)
	})

	app.Post("/deposit", func(c *fiber.Ctx) error {
		return handleDeposit(c, db)
	})

	app.Post("/withdraw", func(c *fiber.Ctx) error {
		return handleWithdraw(c, db)
	})

	if err := app.Listen(":3000"); err != nil {
		logrus.WithError(err).Error("Error starting server")
	}
}
