package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
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

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return &Database{DB: db}
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
	r := gin.Default()
	db := NewDatabase()
	defer db.DB.Close()

	r.POST("/signup", func(c *gin.Context) {
		var req SignupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		if !ValidateName(req.Name) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
			return
		}
		if !ValidateEmail(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
			return
		}
		emailExists, err := CheckDuplicateEmail(db, req.Email)
		if err != nil {
			fmt.Println("Error checking duplicate email:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email"})
			return
		}
		if emailExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
		if !ValidatePassword(req.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
			return
		}
		document := Document{Digits: req.Document}
		if !document.Validate() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document"})
			return
		}
		user := User{
			AccountID: uuid.New(),
			Name:      req.Name,
			Email:     req.Email,
			Document:  document.Digits,
			Password:  req.Password,
		}
		query := `INSERT INTO ccca.account (account_id, name, email, document, password) VALUES ($1, $2, $3, $4, $5)`
		_, err = db.DB.Exec(query, user.AccountID, user.Name, user.Email, user.Document, user.Password)
		if err != nil {
			fmt.Println("Error inserting account:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"accountId": user.AccountID,
		})
	})

	r.GET("/accounts/:account_id", func(c *gin.Context) {
		accountID := c.Param("account_id")
		query := `SELECT account_id, name, email, document, password FROM ccca.account WHERE account_id = $1`
		var user User
		err := db.DB.QueryRow(query, accountID).Scan(&user.AccountID, &user.Name, &user.Email, &user.Document, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
				return
			}
			fmt.Println("Error querying account:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"accountId": user.AccountID,
		})
	})

	if err := r.Run(":3000"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
