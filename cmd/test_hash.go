package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func testHashPassword() {
	password := "TestPass@123"

	// Generate hash for the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		return
	}

	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash: %s\n", hashedPassword)

	// Verify the hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		fmt.Printf("Hash verification failed: %v\n", err)
		return
	}

	fmt.Println("Hash verification successful!")
}

func main() {
	testHashPassword()
}
