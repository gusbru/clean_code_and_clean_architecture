package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
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

	logrus.Info("Connecting to database", logrus.Fields{
		"host":   host,
		"port":   port,
		"dbname": dbname,
	})

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

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Document string `json:"document"`
	Password string `json:"password"`
}

type User struct {
	AccountID uuid.UUID `json:"account_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Document  string    `json:"document"`
	Password  string    `json:"password"`
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
		var req SignupRequest
		if err := c.BodyParser(&req); err != nil {
			logrus.WithError(err).Error("Failed to parse request body")
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": "Invalid request body"})
		}
		logrus.Info("Processing signup", logrus.Fields{
			"email": req.Email,
			"name":  req.Name,
		})
		if !ValidateName(req.Name) {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": "Invalid name"})
		}
		if !ValidateEmail(req.Email) {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": "Invalid email"})
		}
		emailExists, err := CheckDuplicateEmail(db, req.Email)
		if err != nil {
			logrus.WithError(err).Error("Error checking duplicate email")
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"error": "Failed to check email"})
		}
		if emailExists {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": "Email already exists"})
		}
		if !ValidatePassword(req.Password) {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": "Invalid password"})
		}
		document := Document{Digits: req.Document}
		if !document.Validate() {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": "Invalid document"})
		}
		user := User{
			AccountID: uuid.New(),
			Name:      req.Name,
			Email:     req.Email,
			Document:  document.Digits,
			Password:  req.Password,
		}
		logrus.Info("Creating new account", logrus.Fields{
			"account_id": user.AccountID,
			"email":      user.Email,
		})
		query := `INSERT INTO ccca.account (account_id, name, email, document, password) VALUES ($1, $2, $3, $4, $5)`
		_, err = db.DB.Exec(query, user.AccountID, user.Name, user.Email, user.Document, user.Password)
		if err != nil {
			logrus.WithError(err).Error("Error inserting account")
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"error": "Failed to create account"})
		}
		logrus.Info("Account created successfully", logrus.Fields{
			"account_id": user.AccountID,
		})
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{"accountId": user.AccountID})
	})

	app.Get("/accounts/:account_id", func(c *fiber.Ctx) error {
		accountID := c.Params("account_id")
		query := `SELECT account_id, name, email, document, password FROM ccca.account WHERE account_id = $1`
		var user User
		err := db.DB.QueryRow(query, accountID).Scan(&user.AccountID, &user.Name, &user.Email, &user.Document, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.JSON(fiber.Map{"error": "Account not found"})
			}
			logrus.WithError(err).Error("Error querying account")
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"error": "Failed to retrieve account"})
		}
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"accountId": user.AccountID,
		})
	})

	if err := app.Listen(":3000"); err != nil {
		logrus.WithError(err).Error("Error starting server")
	}
}
