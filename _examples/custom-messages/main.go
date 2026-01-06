package main

import (
	"fmt"

	"github.com/hotfixfirst/go-xvalidator"
)

// User demonstrates custom error messages
type User struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,min=18,max=120"`
	Username string `json:"username" validate:"required,min=4,max=20,alphanum"`
	Password string `json:"password" validate:"required,password_strength"`
	Phone    string `json:"phone" validate:"required,mobile_e164"`
	Website  string `json:"website" validate:"omitempty,url"`
}

// Product demonstrates decimal validation messages
type Product struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Price       string `json:"price" validate:"required,decimal=10:2,dgt=0,dlte=1000000"`
	Stock       int    `json:"stock" validate:"required,min=0,max=100000"`
	Description string `json:"description" validate:"required,min=10,max=500"`
}

func main() {
	fmt.Println("=== Custom Error Messages Examples ===\n")

	// Create validator instance with default locale
	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Example 1: Valid user - no errors
	fmt.Println("Example 1: Valid User - No Errors")
	validUser := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      25,
		Username: "johndoe",
		Password: "MyP@ssw0rd123",
		Phone:    "+66812345678",
		Website:  "https://john.example.com",
	}
	if err := v.StructTranslated(validUser); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 2: Missing required field
	fmt.Println("Example 2: Missing Required Field")
	invalid1 := User{
		Name:     "",
		Email:    "test@example.com",
		Age:      30,
		Username: "testuser",
		Password: "Test@123",
		Phone:    "+66812345678",
	}
	if err := v.StructTranslated(invalid1); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 3: Invalid email format
	fmt.Println("Example 3: Invalid Email Format")
	invalid2 := User{
		Name:     "Jane Doe",
		Email:    "not-an-email",
		Age:      28,
		Username: "janedoe",
		Password: "Jane@Pass123",
		Phone:    "+66812345678",
	}
	if err := v.StructTranslated(invalid2); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 4: Age below minimum
	fmt.Println("Example 4: Age Below Minimum")
	invalid3 := User{
		Name:     "Young User",
		Email:    "young@example.com",
		Age:      16,
		Username: "younguser",
		Password: "Young@Pass123",
		Phone:    "+66812345678",
	}
	if err := v.StructTranslated(invalid3); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 5: Username too short
	fmt.Println("Example 5: Username Too Short")
	invalid4 := User{
		Name:     "Bob Smith",
		Email:    "bob@example.com",
		Age:      35,
		Username: "ab",
		Password: "Bob@Pass123",
		Phone:    "+66812345678",
	}
	if err := v.StructTranslated(invalid4); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 6: Weak password
	fmt.Println("Example 6: Weak Password")
	invalid5 := User{
		Name:     "Alice Wonder",
		Email:    "alice@example.com",
		Age:      27,
		Username: "alicew",
		Password: "password",
		Phone:    "+66812345678",
	}
	if err := v.StructTranslated(invalid5); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 7: Invalid phone number
	fmt.Println("Example 7: Invalid Phone Number")
	invalid6 := User{
		Name:     "Charlie Brown",
		Email:    "charlie@example.com",
		Age:      32,
		Username: "charlieb",
		Password: "Charlie@123",
		Phone:    "0812345678",
	}
	if err := v.StructTranslated(invalid6); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 8: Invalid URL format
	fmt.Println("Example 8: Invalid URL Format")
	invalid7 := User{
		Name:     "David Lee",
		Email:    "david@example.com",
		Age:      29,
		Username: "davidlee",
		Password: "David@Pass123",
		Phone:    "+66812345678",
		Website:  "not-a-url",
	}
	if err := v.StructTranslated(invalid7); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 9: Valid product
	fmt.Println("Example 9: Valid Product")
	validProduct := Product{
		Name:        "Gaming Laptop",
		Price:       "45900.50",
		Stock:       25,
		Description: "High-performance gaming laptop with RTX graphics",
	}
	if err := v.StructTranslated(validProduct); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 10: Invalid decimal price
	fmt.Println("Example 10: Invalid Decimal Price")
	invalidProduct := Product{
		Name:        "Smartphone",
		Price:       "-100.00",
		Stock:       50,
		Description: "Latest model smartphone",
	}
	if err := v.StructTranslated(invalidProduct); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 11: Multiple validation errors
	fmt.Println("Example 11: Multiple Validation Errors")
	multiError := User{
		Name:     "A",
		Email:    "bad-email",
		Age:      15,
		Username: "ab",
		Password: "weak",
		Phone:    "123",
	}
	if err := v.StructTranslated(multiError); err != nil {
		fmt.Printf("❌ Errors:\n%v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
}
