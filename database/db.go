// package db

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/adityjoshi/Dosahostel/models"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// var (
// 	bigBoysDB *gorm.DB
// 	boysOneDB *gorm.DB
// )
// var err error

// func InitDB() {

// 	bigBoysDBUser := os.Getenv("DB_BIG_BOYS_USER")
// 	bigBoysDBPassword := os.Getenv("DB_BIG_BOYS_PASSWORD")
// 	bigBoysDBHost := os.Getenv("DB_BIG_BOYS_HOST")
// 	bigBoysDBPort := os.Getenv("DB_BIG_BOYS_PORT")
// 	bigBoysDBName := os.Getenv("DB_BIG_BOYS_NAME")

// 	// Boys One Database
// 	boysOneDBUser := os.Getenv("DB_BOYS_ONE_USER")
// 	boysOneDBPassword := os.Getenv("DB_BOYS_ONE_PASSWORD")
// 	boysOneDBHost := os.Getenv("DB_BOYS_ONE_HOST")
// 	boysOneDBPort := os.Getenv("DB_BOYS_ONE_PORT")
// 	boysOneDBName := os.Getenv("DB_BOYS_ONE_NAME")

// 	// Data Source Names (DSNs)
// 	bigBoysDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", bigBoysDBHost, bigBoysDBUser, bigBoysDBPassword, bigBoysDBName, bigBoysDBPort)
// 	boysOneDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", boysOneDBHost, boysOneDBUser, boysOneDBPassword, boysOneDBName, boysOneDBPort)

// 	// Connecting to Big Boys Database
// 	bigBoysDB, err = gorm.Open(postgres.Open(bigBoysDSN), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Error connecting to the Big Boys database: %v", err)
// 	}

// 	// Connecting to Boys One Database
// 	boysOneDB, err = gorm.Open(postgres.Open(boysOneDSN), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Error connecting to the Boys One database: %v", err)
// 	}

// 	// Ensure the connections are established
// 	sqlBigBoysDB, err := bigBoysDB.DB()
// 	if err != nil {
// 		log.Fatalf("Error getting the Big Boys database object: %v", err)
// 	}
// 	err = sqlBigBoysDB.Ping()
// 	if err != nil {
// 		log.Fatalf("Error pinging the Big Boys database: %v", err)
// 	}

// 	sqlBoysOneDB, err := boysOneDB.DB()
// 	if err != nil {
// 		log.Fatalf("Error getting the Boys One database object: %v", err)
// 	}
// 	err = sqlBoysOneDB.Ping()
// 	if err != nil {
// 		log.Fatalf("Error pinging the Boys One database: %v", err)
// 	}

// 	fmt.Println("Big Boys and Boys One Database connections successful")
// 	err = bigBoysDB.AutoMigrate(&models.Block{})
// 	if err != nil {
// 		log.Fatalf("error migrating models (block): %v", err)
// 	}

// 	// Now migrate the students, wardens, and then users
// 	err = bigBoysDB.AutoMigrate(&models.Student{}, &models.Warden{})
// 	if err != nil {
// 		log.Fatalf("error migrating models (students, wardens): %v", err)
// 	}

// 	err = bigBoysDB.AutoMigrate(&models.User{})
// 	if err != nil {
// 		log.Fatalf("error migrating models (user): %v", err)
// 	}

// 	// Repeat for Boys One DB
// 	err = boysOneDB.AutoMigrate(&models.Block{})
// 	if err != nil {
// 		log.Fatalf("error migrating models (block): %v", err)
// 	}

// 	err = boysOneDB.AutoMigrate(&models.Student{}, &models.Warden{})
// 	if err != nil {
// 		log.Fatalf("error migrating models (students, wardens): %v", err)
// 	}

// 	err = boysOneDB.AutoMigrate(&models.User{})
// 	if err != nil {
// 		log.Fatalf("error migrating models (user): %v", err)
// 	}

// 	fmt.Println("Big Boys and Boys One Database connections successful and tables migrated")
// }

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
	// First, migrate the Block table
	err = bigBoysDB.AutoMigrate(&models.Block{})
	if err != nil {
		log.Fatalf("Error migrating Block table: %v", err)
	}

	// Then, migrate other tables after Block table is confirmed to exist
	err = bigBoysDB.AutoMigrate(&models.Student{}, &models.Warden{})
	if err != nil {
		log.Fatalf("Error migrating models (students, wardens): %v", err)
	}

	err = bigBoysDB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Error migrating models (user): %v", err)
	}

	fmt.Println("Big Boys Database migrations completed successfully")
}

// Boys One DB migrations
func migrateBoysOneDB() {
	err = boysOneDB.AutoMigrate(&models.Block{})
	if err != nil {
		log.Fatalf("Error migrating models (block): %v", err)
	}

	err = boysOneDB.AutoMigrate(&models.Student{}, &models.Warden{}, &models.User{})
	if err != nil {
		log.Fatalf("Error migrating models (students, wardens, user): %v", err)
	}

	fmt.Println("Boys One Database migrations completed successfully")
}
