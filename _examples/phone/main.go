package main

import (
	"fmt"

	"github.com/hotfixfirst/go-xvalidator"
)

// Contact demonstrates E.164 phone number validation
type Contact struct {
	Name            string   `json:"name" validate:"required,min=2,max=100"`
	PrimaryPhone    string   `json:"primary_phone" validate:"required,mobile_e164"`
	SecondaryPhone  string   `json:"secondary_phone" validate:"omitempty,mobile_e164"`
	EmergencyPhones []string `json:"emergency_phones" validate:"omitempty,dive,mobile_e164"`
}

func main() {
	fmt.Println("=== Phone Number Validation Examples ===\n")

	// Create validator instance
	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Example 1: Valid Thai phone number
	fmt.Println("Example 1: Valid Thai Phone Number (+66)")
	contact1 := Contact{
		Name:           "Somchai Dee",
		PrimaryPhone:   "+66812345678",
		SecondaryPhone: "",
	}
	if err := v.Validate(contact1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 2: Valid US phone number
	fmt.Println("Example 2: Valid US Phone Number (+1)")
	contact2 := Contact{
		Name:         "John Smith",
		PrimaryPhone: "+15551234567",
	}
	if err := v.Validate(contact2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 3: Valid UK phone number
	fmt.Println("Example 3: Valid UK Phone Number (+44)")
	contact3 := Contact{
		Name:         "James Bond",
		PrimaryPhone: "+447911123456",
	}
	if err := v.Validate(contact3); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 4: Valid contact with secondary phone
	fmt.Println("Example 4: Valid Contact with Secondary Phone")
	contact4 := Contact{
		Name:           "Alice Wonder",
		PrimaryPhone:   "+66812345678",
		SecondaryPhone: "+66987654321",
	}
	if err := v.Validate(contact4); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 5: Invalid - missing + prefix
	fmt.Println("Example 5: Invalid - Missing + Prefix")
	invalid1 := Contact{
		Name:         "Invalid User",
		PrimaryPhone: "66812345678",
	}
	if err := v.Validate(invalid1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 6: Invalid - too short
	fmt.Println("Example 6: Invalid - Too Short")
	invalid2 := Contact{
		Name:         "Invalid User",
		PrimaryPhone: "+6681",
	}
	if err := v.Validate(invalid2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 7: Invalid - contains spaces
	fmt.Println("Example 7: Invalid - Contains Spaces")
	invalid3 := Contact{
		Name:         "Invalid User",
		PrimaryPhone: "+66 81 234 5678",
	}
	if err := v.Validate(invalid3); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 8: Invalid - contains dashes
	fmt.Println("Example 8: Invalid - Contains Dashes")
	invalid4 := Contact{
		Name:         "Invalid User",
		PrimaryPhone: "+66-81-234-5678",
	}
	if err := v.Validate(invalid4); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 9: Valid array of emergency phones
	fmt.Println("Example 9: Valid Array of Emergency Phones")
	contact9 := Contact{
		Name:         "Bob Builder",
		PrimaryPhone: "+66812345678",
		EmergencyPhones: []string{
			"+66987654321",
			"+66811111111",
			"+66822222222",
		},
	}
	if err := v.Validate(contact9); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 10: Invalid - one invalid phone in array
	fmt.Println("Example 10: Invalid - One Invalid Phone in Array")
	invalid10 := Contact{
		Name:         "Charlie Brown",
		PrimaryPhone: "+66812345678",
		EmergencyPhones: []string{
			"+66987654321",
			"0811111111", // Invalid - no + prefix
			"+66822222222",
		},
	}
	if err := v.Validate(invalid10); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 11: Single phone validation
	fmt.Println("Example 11: Single Phone Field Validation")
	phone := "+66812345678"
	if err := v.Var(phone, "mobile_e164"); err != nil {
		fmt.Printf("❌ Phone '%s' validation failed: %v\n", phone, err)
	} else {
		fmt.Printf("✅ Phone '%s' is valid!\n", phone)
	}
	fmt.Println()

	// Example 12: Single invalid phone validation
	fmt.Println("Example 12: Single Invalid Phone Validation")
	invalidPhone := "0812345678"
	if err := v.Var(invalidPhone, "mobile_e164"); err != nil {
		fmt.Printf("❌ Phone '%s' validation failed: %v\n", invalidPhone, err)
	} else {
		fmt.Printf("✅ Phone '%s' is valid!\n", invalidPhone)
	}
}
