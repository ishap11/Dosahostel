package main

import (
	"fmt"

	db "github.com/adityjoshi/Dosahostel/database"
)

func main() {
	fmt.Print("jai shree ram \n")
	fmt.Print("go run main.go")

	db.InitDB()
	fmt.Print("jai shree ram \n")
}
