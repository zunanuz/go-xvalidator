package main

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"

	xvalidator "github.com/hotfixfirst/go-xvalidator"
)

// Example 1: Thai ID Card Validator
// Example 1: Thai ID Card Validator
// Thai National ID Card uses MOD 11 checksum algorithm
// Structure: 12 digits + 1 checksum digit (total 13 digits)
// Checksum calculation:
//  1. Multiply each of the first 12 digits by (13 - position)
//  2. Sum all products
//  3. Calculate: (11 - (sum % 11)) % 11
//  4. If result is 10, use 0 as checksum
//  5. Last digit must equal the calculated checksum
func validateThaiIDCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()

	// Must be exactly 13 digits
	if len(idCard) != 13 {
		return false
	}

	// All characters must be numeric
	for _, char := range idCard {
		if char < '0' || char > '9' {
			return false
		}
	}

	// Calculate checksum using MOD 11 algorithm
	sum := 0
	for i := 0; i < 12; i++ {
		digit := int(idCard[i] - '0')
		sum += digit * (13 - i) // Weight: 13, 12, 11, ..., 2
	}

	checksum := (11 - (sum % 11)) % 11
	if checksum == 10 {
		checksum = 0 // Replace 10 with 0
	}

	lastDigit := int(idCard[12] - '0')
	return checksum == lastDigit
}

type ThaiCitizen struct {
	Name   string `validate:"required,min=2,max=100"`
	IDCard string `validate:"required,thai_id_card"`
}

// Example 2: Business Hours Validator
func validateBusinessHours(fl validator.FieldLevel) bool {
	timeStr := fl.Field().String()
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return false
	}
	start, _ := time.Parse("15:04", "09:00")
	end, _ := time.Parse("15:04", "17:00")
	return (t.Equal(start) || t.After(start)) && (t.Before(end) || t.Equal(end))
}

type Appointment struct {
	CustomerName string `validate:"required"`
	Time         string `validate:"required,business_hours"`
}

// Example 3: Thai Phone Number Validator
func validateThaiPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if len(phone) != 10 || phone[0] != '0' {
		return false
	}
	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}
	prefix := phone[0:2]
	validPrefixes := []string{"02", "06", "08", "09"}
	for _, valid := range validPrefixes {
		if prefix == valid {
			return true
		}
	}
	return false
}

type Contact struct {
	Name      string `validate:"required"`
	ThaiPhone string `validate:"required,thai_phone"`
}

// Example 4: Future Date Validator
func validateFutureDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}
	now := time.Now()
	return date.After(now)
}

type Event struct {
	Name      string `validate:"required"`
	StartDate string `validate:"required,future_date"`
}

// Example 5: Decimal Range Validator Factory
func validateDecimalRange(min, max float64) validator.Func {
	return func(fl validator.FieldLevel) bool {
		decStr := fl.Field().String()
		dec, err := decimal.NewFromString(decStr)
		if err != nil {
			return false
		}
		minDec := decimal.NewFromFloat(min)
		maxDec := decimal.NewFromFloat(max)
		return dec.GreaterThanOrEqual(minDec) && dec.LessThanOrEqual(maxDec)
	}
}

type Product struct {
	Name  string `validate:"required"`
	Price string `validate:"required,decimal=10:2,product_price"`
}

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════╗")
	fmt.Println("║     go-xvalidator - Advanced Examples             ║")
	fmt.Println("║     Custom Validators & Advanced Techniques       ║")
	fmt.Println("╚═══════════════════════════════════════════════════╝\n")

	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Register all custom validators
	v.GetValidator().RegisterValidation("thai_id_card", validateThaiIDCard)
	v.GetValidator().RegisterValidation("business_hours", validateBusinessHours)
	v.GetValidator().RegisterValidation("thai_phone", validateThaiPhone)
	v.GetValidator().RegisterValidation("future_date", validateFutureDate)
	v.GetValidator().RegisterValidation("product_price", validateDecimalRange(1.00, 1000000.00))

	// Example 1: Thai ID Card
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("Example 1: Custom Validator - Thai ID Card")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("\nℹ️  Thai ID Card Structure: 12 digits + 1 checksum digit (MOD 11)")
	fmt.Println("   The last digit is calculated from the first 12 digits")
	fmt.Println("   Example: 1103700166114")
	fmt.Println("            └────────────┘└─ checksum = 4")
	fmt.Println("            12 digits")

	valid1 := ThaiCitizen{
		Name:   "สมชาย ใจดี",
		IDCard: "1103700166114",
	}
	fmt.Println("\n✅ Valid Thai Citizen:")
	fmt.Printf("   ID Card: %s (checksum digit: 4 ✓)\n", valid1.IDCard)
	if err := v.StructTranslated(valid1); err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✓ Validation passed")
	}

	invalid1 := ThaiCitizen{
		Name:   "สมหญิง ใจดี",
		IDCard: "1103700166115",
	}
	fmt.Println("\n❌ Invalid Thai ID Card (wrong checksum):")
	fmt.Printf("   ID Card: %s (checksum digit: 5 ✗, should be 4)\n", invalid1.IDCard)
	if err := v.StructTranslated(invalid1); err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✓ Validation passed")
	}

	// Example 2: Business Hours
	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("Example 2: Custom Validator - Business Hours")
	fmt.Println("═══════════════════════════════════════════════════")
	valid2 := Appointment{
		CustomerName: "John Doe",
		Time:         "14:30",
	}
	fmt.Println("\n✅ Valid Appointment (Business Hours):")
	fmt.Printf("   Time: %s\n", valid2.Time)
	if err := v.StructTranslated(valid2); err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✓ Validation passed")
	}

	invalid2 := Appointment{
		CustomerName: "Jane Smith",
		Time:         "18:30",
	}
	fmt.Println("\n❌ Invalid Appointment (After Hours):")
	fmt.Printf("   Time: %s\n", invalid2.Time)
	if err := v.StructTranslated(invalid2); err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✓ Validation passed")
	}

	// Example 3: Thai Phone Number
	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("Example 3: Custom Validator - Thai Phone Number")
	fmt.Println("═══════════════════════════════════════════════════")
	valid3 := Contact{
		Name:      "สมชาย รักษ์ดี",
		ThaiPhone: "0812345678",
	}
	fmt.Println("\n✅ Valid Thai Phone:")
	fmt.Printf("   Phone: %s\n", valid3.ThaiPhone)
	if err := v.StructTranslated(valid3); err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✓ Validation passed")
	}

	invalid3 := Contact{
		Name:      "สมหญิง ใจดี",
		ThaiPhone: "0112345678",
	}
	fmt.Println("\n❌ Invalid Thai Phone (Wrong Prefix):")
	fmt.Printf("   Phone: %s\n", invalid3.ThaiPhone)
	if err := v.StructTranslated(invalid3); err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✓ Validation passed")
	}

	// Example 4: Future Date
	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("Example 4: Custom Validator - Future Date")
	fmt.Println("═══════════════════════════════════════════════════")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	valid4 := Event{
		Name:      "Tech Conference 2024",
		StartDate: tomorrow,
	}
	fmt.Println("\n✅ Valid Future Event:")
	fmt.Printf("   Start: %s\n", valid4.StartDate)
	if err := v.StructTranslated(valid4); err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✓ Validation passed")
	}

	invalid4 := Event{
		Name:      "Past Event",
		StartDate: "2020-01-01",
	}
	fmt.Println("\n❌ Invalid Event (Past Date):")
	fmt.Printf("   Start: %s\n", invalid4.StartDate)
	if err := v.StructTranslated(invalid4); err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✓ Validation passed")
	}

	// Example 5: Decimal Range
	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("Example 5: Custom Validator - Decimal Range")
	fmt.Println("═══════════════════════════════════════════════════")
	valid5 := Product{
		Name:  "Laptop Computer",
		Price: "45900.00",
	}
	fmt.Println("\n✅ Valid Product:")
	fmt.Printf("   Price: %s THB (1-1,000,000)\n", valid5.Price)
	if err := v.StructTranslated(valid5); err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✓ Validation passed")
	}

	invalid5 := Product{
		Name:  "Expensive Item",
		Price: "2000000.00",
	}
	fmt.Println("\n❌ Invalid Product (Price > 1,000,000):")
	fmt.Printf("   Price: %s THB\n", invalid5.Price)
	if err := v.StructTranslated(invalid5); err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✓ Validation passed")
	}

	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("✓ All advanced examples completed!")
	fmt.Println("═══════════════════════════════════════════════════")
}
