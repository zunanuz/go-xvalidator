package main

import (
	"fmt"

	"github.com/hotfixfirst/go-xvalidator"
)

// PriceData demonstrates decimal comparison validators
type PriceData struct {
	RegularPrice string `json:"regular_price" validate:"required,dgte=0"`
	SalePrice    string `json:"sale_price" validate:"required,dgte=0,dlt=1000000"`
	Discount     string `json:"discount" validate:"required,dgte=0,dlte=100"`
	MinOrder     string `json:"min_order" validate:"required,dgt=0"`
}

// BankAccount demonstrates decimal equality validators
type BankAccount struct {
	Balance     string `json:"balance" validate:"required,dgte=0"`
	MinBalance  string `json:"min_balance" validate:"required,deq=100"`
	FeeAmount   string `json:"fee_amount" validate:"required,dneq=0"`
	MaxTransfer string `json:"max_transfer" validate:"required,dlte=50000"`
}

// Invoice demonstrates precise decimal validation
type Invoice struct {
	Subtotal string `json:"subtotal" validate:"required,decimal=10:2,dgte=0"`
	Tax      string `json:"tax" validate:"required,decimal=10:2,dgte=0"`
	Total    string `json:"total" validate:"required,decimal=10:2,dgte=0"`
	Paid     string `json:"paid" validate:"required,decimal=10:2,dgte=0"`
}

func main() {
	fmt.Println("=== Decimal Validation Examples ===\n")

	// Create validator instance
	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Example 1: Valid price data
	fmt.Println("Example 1: Valid Price Data")
	validPrice := PriceData{
		RegularPrice: "99.99",
		SalePrice:    "79.99",
		Discount:     "20.00",
		MinOrder:     "1.00",
	}

	if err := v.Validate(validPrice); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 2: Invalid - negative regular price
	fmt.Println("Example 2: Invalid - Negative Regular Price")
	invalidPrice1 := PriceData{
		RegularPrice: "-10.00",
		SalePrice:    "79.99",
		Discount:     "20.00",
		MinOrder:     "1.00",
	}

	if err := v.Validate(invalidPrice1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 3: Invalid - discount over 100%
	fmt.Println("Example 3: Invalid - Discount Over 100%")
	invalidPrice2 := PriceData{
		RegularPrice: "99.99",
		SalePrice:    "79.99",
		Discount:     "120.00",
		MinOrder:     "1.00",
	}

	if err := v.Validate(invalidPrice2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 4: Invalid - min order must be greater than 0
	fmt.Println("Example 4: Invalid - Min Order Must Be Greater Than 0")
	invalidPrice3 := PriceData{
		RegularPrice: "99.99",
		SalePrice:    "79.99",
		Discount:     "20.00",
		MinOrder:     "0.00",
	}

	if err := v.Validate(invalidPrice3); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 5: Valid bank account
	fmt.Println("Example 5: Valid Bank Account")
	validAccount := BankAccount{
		Balance:     "5000.00",
		MinBalance:  "100.00",
		FeeAmount:   "5.00",
		MaxTransfer: "50000.00",
	}

	if err := v.Validate(validAccount); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 6: Invalid - min balance not equal to 100
	fmt.Println("Example 6: Invalid - Min Balance Must Equal 100")
	invalidAccount1 := BankAccount{
		Balance:     "5000.00",
		MinBalance:  "50.00",
		FeeAmount:   "5.00",
		MaxTransfer: "50000.00",
	}

	if err := v.Validate(invalidAccount1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 7: Invalid - fee amount cannot be 0
	fmt.Println("Example 7: Invalid - Fee Amount Cannot Be 0")
	invalidAccount2 := BankAccount{
		Balance:     "5000.00",
		MinBalance:  "100.00",
		FeeAmount:   "0.00",
		MaxTransfer: "50000.00",
	}

	if err := v.Validate(invalidAccount2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 8: Valid invoice with precision
	fmt.Println("Example 8: Valid Invoice with Precision")
	validInvoice := Invoice{
		Subtotal: "1000.50",
		Tax:      "70.04",
		Total:    "1070.54",
		Paid:     "1070.54",
	}

	if err := v.Validate(validInvoice); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 9: Invalid - too many decimal places
	fmt.Println("Example 9: Invalid - Too Many Decimal Places")
	invalidInvoice := Invoice{
		Subtotal: "1000.505",
		Tax:      "70.04",
		Total:    "1070.54",
		Paid:     "1070.54",
	}

	if err := v.Validate(invalidInvoice); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 10: Single decimal validation
	fmt.Println("Example 10: Single Decimal Field Validation")
	amount := "250.75"

	if err := v.Var(amount, "dgte=0,dlte=1000"); err != nil {
		fmt.Printf("❌ Amount validation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Amount '%s' is valid!\n", amount)
	}
	fmt.Println()

	// Example 11: Decimal greater than comparison
	fmt.Println("Example 11: Decimal Greater Than")
	value1 := "100.50"

	if err := v.Var(value1, "dgt=100"); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Printf("✅ '%s' is greater than 100\n", value1)
	}
	fmt.Println()

	// Example 12: Decimal less than comparison
	fmt.Println("Example 12: Decimal Less Than")
	value2 := "99.99"

	if err := v.Var(value2, "dlt=100"); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Printf("✅ '%s' is less than 100\n", value2)
	}
}
