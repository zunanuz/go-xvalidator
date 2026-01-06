package main

import (
	"fmt"

	"github.com/hotfixfirst/go-xvalidator"
)

// User demonstrates basic validation rules
type User struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,min=18,max=120"`
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Website  string `json:"website" validate:"omitempty,url"`
}

// Product demonstrates numeric and string validations
type Product struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"required,min=10,max=500"`
	SKU         string  `json:"sku" validate:"required,alphanum,len=8"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Stock       int     `json:"stock" validate:"required,min=0"`
}

func main() {
	fmt.Println("=== Basic Validation Examples ===\n")

	// Create validator instance
	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Example 1: Valid user
	fmt.Println("Example 1: Valid User")
	validUser := User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Age:      25,
		Username: "johndoe123",
		Website:  "https://johndoe.com",
	}

	if err := v.Validate(validUser); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 2: Invalid user - missing required field
	fmt.Println("Example 2: Invalid User - Missing Required Field")
	invalidUser1 := User{
		Name:     "",
		Email:    "john@example.com",
		Age:      25,
		Username: "johndoe",
	}

	if err := v.Validate(invalidUser1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 3: Invalid user - invalid email
	fmt.Println("Example 3: Invalid User - Invalid Email")
	invalidUser2 := User{
		Name:     "John Doe",
		Email:    "invalid-email",
		Age:      25,
		Username: "johndoe",
	}

	if err := v.Validate(invalidUser2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 4: Invalid user - age too young
	fmt.Println("Example 4: Invalid User - Age Too Young")
	invalidUser3 := User{
		Name:     "Young User",
		Email:    "young@example.com",
		Age:      16,
		Username: "younguser",
	}

	if err := v.Validate(invalidUser3); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 5: Invalid user - username too short
	fmt.Println("Example 5: Invalid User - Username Too Short")
	invalidUser4 := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      25,
		Username: "ab",
	}

	if err := v.Validate(invalidUser4); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 6: Valid product
	fmt.Println("Example 6: Valid Product")
	validProduct := Product{
		Name:        "Gaming Laptop",
		Description: "High-performance gaming laptop with RTX graphics",
		SKU:         "LAP12345",
		Price:       1299.99,
		Stock:       50,
	}

	if err := v.Validate(validProduct); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 7: Invalid product - price negative
	fmt.Println("Example 7: Invalid Product - Negative Price")
	invalidProduct := Product{
		Name:        "Gaming Laptop",
		Description: "High-performance gaming laptop",
		SKU:         "LAP12345",
		Price:       -100,
		Stock:       50,
	}

	if err := v.Validate(invalidProduct); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 8: Using Var for single field validation
	fmt.Println("Example 8: Single Field Validation with Var")
	email := "test@example.com"

	if err := v.Var(email, "required,email"); err != nil {
		fmt.Printf("❌ Email validation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Email '%s' is valid!\n", email)
	}
	fmt.Println()

	// Example 9: Single field validation - invalid
	fmt.Println("Example 9: Single Field Validation - Invalid Email")
	invalidEmail := "not-an-email"

	if err := v.Var(invalidEmail, "required,email"); err != nil {
		fmt.Printf("❌ Email validation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Email '%s' is valid!\n", invalidEmail)
	}
	fmt.Println()

	// Example 10: Optional field validation
	fmt.Println("Example 10: Optional Field - Empty Website")
	userWithoutWebsite := User{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Age:      30,
		Username: "janedoe",
		Website:  "", // Empty optional field should pass
	}

	if err := v.Validate(userWithoutWebsite); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed! (empty optional field is ok)")
	}
}
