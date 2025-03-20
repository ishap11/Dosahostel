package db

import (
	"fmt"
	"log"
	"os"

	"github.com/adityjoshi/Dosahostel/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	bigBoysDB *gorm.DB
	boysOneDB *gorm.DB
)
var err error

func InitDB() {

	// Big Boys Database credentials
	bigBoysDBUser := os.Getenv("DB_BIG_BOYS_USER")
	bigBoysDBPassword := os.Getenv("DB_BIG_BOYS_PASSWORD")
	bigBoysDBHost := os.Getenv("DB_BIG_BOYS_HOST")
	bigBoysDBPort := os.Getenv("DB_BIG_BOYS_PORT")
	bigBoysDBName := os.Getenv("DB_BIG_BOYS_NAME")

	// Boys One Database credentials
	boysOneDBUser := os.Getenv("DB_BOYS_ONE_USER")
	boysOneDBPassword := os.Getenv("DB_BOYS_ONE_PASSWORD")
	boysOneDBHost := os.Getenv("DB_BOYS_ONE_HOST")
	boysOneDBPort := os.Getenv("DB_BOYS_ONE_PORT")
	boysOneDBName := os.Getenv("DB_BOYS_ONE_NAME")

	// Data Source Names (DSNs) for both databases
	bigBoysDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", bigBoysDBHost, bigBoysDBUser, bigBoysDBPassword, bigBoysDBName, bigBoysDBPort)
	boysOneDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", boysOneDBHost, boysOneDBUser, boysOneDBPassword, boysOneDBName, boysOneDBPort)

	// Connecting to Big Boys Database
	bigBoysDB, err = gorm.Open(postgres.Open(bigBoysDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the Big Boys database: %v", err)
	}

	// Connecting to Boys One Database
	boysOneDB, err = gorm.Open(postgres.Open(boysOneDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the Boys One database: %v", err)
	}

	// Ensure the connections are established
	sqlBigBoysDB, err := bigBoysDB.DB()
	if err != nil {
		log.Fatalf("Error getting the Big Boys database object: %v", err)
	}
	err = sqlBigBoysDB.Ping()
	if err != nil {
		log.Fatalf("Error pinging the Big Boys database: %v", err)
	}

	sqlBoysOneDB, err := boysOneDB.DB()
	if err != nil {
		log.Fatalf("Error getting the Boys One database object: %v", err)
	}
	err = sqlBoysOneDB.Ping()
	if err != nil {
		log.Fatalf("Error pinging the Boys One database: %v", err)
	}

	// Print successful connection messages
	fmt.Println("Big Boys and Boys One Database connections successful")

	migrateBigBoysDB()
	migrateBoysOneDB()
}

func migrateBigBoysDB() {
	// Manually check if the "blocks" table exists

	// Then migrate other tables
	err = bigBoysDB.AutoMigrate(&models.Student{}, &models.Warden{})
	if err != nil {
		log.Fatalf("Error migrating models (students, wardens): %v", err)
	}

	fmt.Println("Big Boys Database migrations completed successfully")
}

// Boys One DB migrations
func migrateBoysOneDB() {

	err = boysOneDB.AutoMigrate(&models.Student{}, &models.Warden{})
	if err != nil {
		log.Fatalf("Error migrating models (students, wardens, user): %v", err)
	}

	fmt.Println("Boys One Database migrations completed successfully")
}

func GetDB(region string) (*gorm.DB, error) {
	switch region {
	case "BH1", "BH6":
		return bigBoysDB, nil
	case "BH2", "BH3", "BH4", "BH5":
		return boysOneDB, nil
	default:
		return nil, fmt.Errorf("invalid region: %s", region)
	}
}
