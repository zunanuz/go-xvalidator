package main

import (
	"fmt"

	"github.com/hotfixfirst/go-xvalidator"
)

// UserRegistration represents a complete user registration form
type UserRegistration struct {
	// Personal Information
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,mobile_e164"`

	// Account Security
	Username        string `json:"username" validate:"required,min=4,max=20,alphanum"`
	Password        string `json:"password" validate:"required,password_strength"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`

	// Profile
	DateOfBirth string `json:"date_of_birth" validate:"required"`
	Gender      string `json:"gender" validate:"required,oneof=male female other prefer_not_to_say"`
	Website     string `json:"website" validate:"omitempty,url"`
	Bio         string `json:"bio" validate:"omitempty,max=500"`

	// Address
	Street     string `json:"street" validate:"required,min=5,max=200"`
	City       string `json:"city" validate:"required,min=2,max=100"`
	State      string `json:"state" validate:"required,len=2"`
	PostalCode string `json:"postal_code" validate:"required,len=5"`
	Country    string `json:"country" validate:"required,len=2"`

	// Terms
	AcceptTerms     bool `json:"accept_terms" validate:"required,eq=true"`
	AcceptMarketing bool `json:"accept_marketing"` // Optional
	Age             int  `json:"age" validate:"required,min=18,max=120"`
}

func main() {
	fmt.Println("=== User Registration Validation Example ===\n")

	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Example 1: Valid registration
	fmt.Println("Example 1: Valid User Registration")
	validReg := UserRegistration{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Phone:           "+66812345678",
		Username:        "johndoe123",
		Password:        "SecureP@ss123",
		ConfirmPassword: "SecureP@ss123",
		DateOfBirth:     "1990-01-15",
		Gender:          "male",
		Website:         "https://johndoe.com",
		Bio:             "Software developer passionate about technology.",
		Street:          "123 Main Street, Apt 4B",
		City:            "Bangkok",
		State:           "BK",
		PostalCode:      "10110",
		Country:         "TH",
		AcceptTerms:     true,
		AcceptMarketing: true,
		Age:             34,
	}

	if err := v.StructTranslated(validReg); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Registration validation passed!")
		fmt.Printf("   User: %s %s\n", validReg.FirstName, validReg.LastName)
		fmt.Printf("   Email: %s\n", validReg.Email)
		fmt.Printf("   Username: %s\n", validReg.Username)
	}
	fmt.Println()

	// Example 2: Invalid - password mismatch
	fmt.Println("Example 2: Invalid - Password Confirmation Mismatch")
	invalidReg1 := UserRegistration{
		FirstName:       "Jane",
		LastName:        "Smith",
		Email:           "jane.smith@example.com",
		Phone:           "+66823456789",
		Username:        "janesmith",
		Password:        "MyP@ssw0rd123",
		ConfirmPassword: "DifferentP@ss456", // Doesn't match
		DateOfBirth:     "1992-05-20",
		Gender:          "female",
		Street:          "456 Oak Avenue",
		City:            "Bangkok",
		State:           "BK",
		PostalCode:      "10120",
		Country:         "TH",
		AcceptTerms:     true,
		Age:             31,
	}

	if err := v.StructTranslated(invalidReg1); err != nil {
		fmt.Printf("❌ Validation failed:\n")
		fmt.Printf("   %v\n", err)
	}
	fmt.Println()

	// Example 3: Invalid - weak password
	fmt.Println("Example 3: Invalid - Weak Password")
	invalidReg2 := UserRegistration{
		FirstName:       "Bob",
		LastName:        "Johnson",
		Email:           "bob@example.com",
		Phone:           "+66834567890",
		Username:        "bobjohnson",
		Password:        "weak", // Too weak
		ConfirmPassword: "weak",
		DateOfBirth:     "1988-08-10",
		Gender:          "male",
		Street:          "789 Pine Road",
		City:            "Bangkok",
		State:           "BK",
		PostalCode:      "10130",
		Country:         "TH",
		AcceptTerms:     true,
		Age:             36,
	}

	if err := v.StructTranslated(invalidReg2); err != nil {
		fmt.Printf("❌ Validation failed:\n")
		fmt.Printf("   %v\n", err)
	}
	fmt.Println()

	// Example 4: Invalid - underage user
	fmt.Println("Example 4: Invalid - User Under 18")
	invalidReg3 := UserRegistration{
		FirstName:       "Alice",
		LastName:        "Brown",
		Email:           "alice@example.com",
		Phone:           "+66845678901",
		Username:        "alicebrown",
		Password:        "ValidP@ss123",
		ConfirmPassword: "ValidP@ss123",
		DateOfBirth:     "2010-03-15",
		Gender:          "female",
		Street:          "321 Elm Street",
		City:            "Bangkok",
		State:           "BK",
		PostalCode:      "10140",
		Country:         "TH",
		AcceptTerms:     true,
		Age:             15, // Under 18
	}

	if err := v.StructTranslated(invalidReg3); err != nil {
		fmt.Printf("❌ Validation failed:\n")
		fmt.Printf("   %v\n", err)
	}
	fmt.Println()

	// Example 5: Invalid - invalid email and phone
	fmt.Println("Example 5: Invalid - Invalid Email and Phone Format")
	invalidReg4 := UserRegistration{
		FirstName:       "Charlie",
		LastName:        "Wilson",
		Email:           "not-an-email", // Invalid
		Phone:           "0856789012",   // Missing +
		Username:        "charliewilson",
		Password:        "ValidP@ss123",
		ConfirmPassword: "ValidP@ss123",
		DateOfBirth:     "1995-11-25",
		Gender:          "male",
		Street:          "555 Maple Drive",
		City:            "Bangkok",
		State:           "BK",
		PostalCode:      "10150",
		Country:         "TH",
		AcceptTerms:     true,
		Age:             28,
	}

	if err := v.StructTranslated(invalidReg4); err != nil {
		fmt.Printf("❌ Validation failed:\n")
		fmt.Printf("   %v\n", err)
	}
	fmt.Println()

	// Example 6: Invalid - didn't accept terms
	fmt.Println("Example 6: Invalid - Terms Not Accepted")
	invalidReg5 := UserRegistration{
		FirstName:       "David",
		LastName:        "Miller",
		Email:           "david@example.com",
		Phone:           "+66867890123",
		Username:        "davidmiller",
		Password:        "ValidP@ss123",
		ConfirmPassword: "ValidP@ss123",
		DateOfBirth:     "1987-07-30",
		Gender:          "male",
		Street:          "777 Oak Boulevard",
		City:            "Bangkok",
		State:           "BK",
		PostalCode:      "10160",
		Country:         "TH",
		AcceptTerms:     false, // Must be true
		Age:             37,
	}

	if err := v.StructTranslated(invalidReg5); err != nil {
		fmt.Printf("❌ Validation failed:\n")
		fmt.Printf("   %v\n", err)
	}
	fmt.Println()

	// Example 7: Invalid - username too short
	fmt.Println("Example 7: Invalid - Username Too Short")
	invalidReg6 := UserRegistration{
		FirstName:       "Emma",
		LastName:        "Davis",
		Email:           "emma@example.com",
		Phone:           "+66878901234",
		Username:        "em", // Too short (min=4)
		Password:        "ValidP@ss123",
		ConfirmPassword: "ValidP@ss123",
		DateOfBirth:     "1993-04-12",
		Gender:          "female",
		Street:          "888 Pine Avenue",
		City:            "Bangkok",
		State:           "BK",
		PostalCode:      "10170",
		Country:         "TH",
		AcceptTerms:     true,
		Age:             31,
	}

	if err := v.StructTranslated(invalidReg6); err != nil {
		fmt.Printf("❌ Validation failed:\n")
		fmt.Printf("   %v\n", err)
	}
	fmt.Println()

	// Example 8: Valid - minimal optional fields
	fmt.Println("Example 8: Valid - Minimal Optional Fields")
	validReg2 := UserRegistration{
		FirstName:       "Frank",
		LastName:        "Taylor",
		Email:           "frank@example.com",
		Phone:           "+66889012345",
		Username:        "franktaylor",
		Password:        "ValidP@ss123",
		ConfirmPassword: "ValidP@ss123",
		DateOfBirth:     "1991-09-05",
		Gender:          "male",
		Website:         "", // Optional
		Bio:             "", // Optional
		Street:          "999 Cedar Lane",
		City:            "Bangkok",
		State:           "BK",
		PostalCode:      "10180",
		Country:         "TH",
		AcceptTerms:     true,
		AcceptMarketing: false, // Optional
		Age:             33,
	}

	if err := v.StructTranslated(validReg2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Registration validated with minimal optional fields!")
	}
	fmt.Println()

	fmt.Println("=== User Registration Validation Complete ===")
}
