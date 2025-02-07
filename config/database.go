package config

import (
	"fmt"
	"log"
	"os"
	"order-inventory/models"
	"io/ioutil"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	// Auto Migrate
	DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{}, &models.PriceHistory{})

	// Apply triggers
	if err := applyTriggers(); err != nil {
		log.Printf("Error applying triggers: %v", err)
	}
}

func applyTriggers() error {
	// Read triggers SQL file
	content, err := ioutil.ReadFile("migrations/triggers.sql")
	if err != nil {
		return fmt.Errorf("error reading triggers file: %v", err)
	}

	// Convert content to string and normalize line endings
	sqlContent := strings.ReplaceAll(string(content), "\r\n", "\n")

	// Split by semicolon but preserve semicolons within function definitions
	var statements []string
	var currentStmt strings.Builder
	lines := strings.Split(sqlContent, "\n")
	inFunction := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		if strings.Contains(line, "CREATE OR REPLACE FUNCTION") {
			inFunction = true
		}

		currentStmt.WriteString(line)
		currentStmt.WriteString("\n")

		if strings.HasSuffix(line, "LANGUAGE plpgsql;") {
			inFunction = false
			statements = append(statements, currentStmt.String())
			currentStmt.Reset()
		} else if !inFunction && strings.HasSuffix(line, ";") {
			statements = append(statements, currentStmt.String())
			currentStmt.Reset()
		}
	}

	// Execute each statement
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if err := DB.Exec(stmt).Error; err != nil {
			log.Printf("Error executing statement: %s", stmt)
			return fmt.Errorf("error executing trigger: %v", err)
		}
	}

	return nil
}
