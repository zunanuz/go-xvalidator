package main

import (
	"fmt"

	"github.com/hotfixfirst/go-xvalidator"
)

// UserAccount demonstrates password strength validation
type UserAccount struct {
	Username string `json:"username" validate:"required,min=4,max=20,alphanum"`
	Password string `json:"password" validate:"required,password_strength"`
	Email    string `json:"email" validate:"required,email"`
}

func main() {
	fmt.Println("=== Password Strength Validation Examples ===\n")

	// Create validator instance
	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Example 1: Valid strong password
	fmt.Println("Example 1: Valid Strong Password")
	user1 := UserAccount{
		Username: "john2024",
		Password: "MyP@ssw0rd!",
		Email:    "john@example.com",
	}
	if err := v.Validate(user1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 2: Valid password with numbers and symbols
	fmt.Println("Example 2: Valid Password with Numbers and Symbols")
	user2 := UserAccount{
		Username: "alice2024",
		Password: "Secure#Pass123",
		Email:    "alice@example.com",
	}
	if err := v.Validate(user2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 3: Valid password with mixed characters
	fmt.Println("Example 3: Valid Password with Mixed Characters")
	user3 := UserAccount{
		Username: "bob123",
		Password: "Tr0ng$P@ss",
		Email:    "bob@example.com",
	}
	if err := v.Validate(user3); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 4: Invalid - too short
	fmt.Println("Example 4: Invalid - Too Short (< 8 characters)")
	invalid1 := UserAccount{
		Username: "user1",
		Password: "Short1!",
		Email:    "user1@example.com",
	}
	if err := v.Validate(invalid1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 5: Invalid - no uppercase
	fmt.Println("Example 5: Invalid - No Uppercase Letters")
	invalid2 := UserAccount{
		Username: "user2",
		Password: "password123!",
		Email:    "user2@example.com",
	}
	if err := v.Validate(invalid2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 6: Invalid - no lowercase
	fmt.Println("Example 6: Invalid - No Lowercase Letters")
	invalid3 := UserAccount{
		Username: "user3",
		Password: "PASSWORD123!",
		Email:    "user3@example.com",
	}
	if err := v.Validate(invalid3); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 7: Invalid - no digits
	fmt.Println("Example 7: Invalid - No Digits")
	invalid4 := UserAccount{
		Username: "user4",
		Password: "Password!",
		Email:    "user4@example.com",
	}
	if err := v.Validate(invalid4); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 8: Invalid - no special characters
	fmt.Println("Example 8: Invalid - No Special Characters")
	invalid5 := UserAccount{
		Username: "user5",
		Password: "Password123",
		Email:    "user5@example.com",
	}
	if err := v.Validate(invalid5); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 9: Valid password with space
	fmt.Println("Example 9: Valid Password with Space")
	user9 := UserAccount{
		Username: "user9",
		Password: "Test 123!",
		Email:    "user9@example.com",
	}
	if err := v.Validate(user9); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 10: Valid password with multiple special chars
	fmt.Println("Example 10: Valid Password with Multiple Special Characters")
	user10 := UserAccount{
		Username: "secure123",
		Password: "P@$$w0rd!#",
		Email:    "secure@example.com",
	}
	if err := v.Validate(user10); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 11: Single password validation
	fmt.Println("Example 11: Single Password Field Validation")
	password := "MyP@ssw0rd123"
	if err := v.Var(password, "password_strength"); err != nil {
		fmt.Printf("❌ Password validation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Password is strong!\n")
	}
	fmt.Println()

	// Example 12: Single invalid password validation
	fmt.Println("Example 12: Single Invalid Password Validation")
	weakPassword := "weak"
	if err := v.Var(weakPassword, "password_strength"); err != nil {
		fmt.Printf("❌ Password validation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Password is strong!\n")
	}
	fmt.Println()

	// Example 13: Valid long password
	fmt.Println("Example 13: Valid Long Password")
	user13 := UserAccount{
		Username: "power",
		Password: "Th1s!sAV3ryL0ngAndSecur3P@ssw0rd",
		Email:    "power@example.com",
	}
	if err := v.Validate(user13); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 14: Valid password with underscores
	fmt.Println("Example 14: Valid Password with Underscores")
	user14 := UserAccount{
		Username: "developer",
		Password: "Dev_P@ss123",
		Email:    "dev@example.com",
	}
	if err := v.Validate(user14); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
}
