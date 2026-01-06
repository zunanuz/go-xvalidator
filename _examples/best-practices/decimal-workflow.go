package main

import (
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/hotfixfirst/go-xvalidator"
)

// PriceRequest represents an API request with price data
// Best Practice: Use string for monetary values in API layer
type PriceRequest struct {
	ProductName string `json:"product_name" validate:"required,min=3"`
	BasePrice   string `json:"base_price" validate:"required,decimal=10:2,dgt=0"`
	Quantity    int    `json:"quantity" validate:"required,min=1,max=1000"`
	DiscountPct string `json:"discount_pct" validate:"required,decimal=5:2,dgte=0,dlte=100"`
}

// PriceResponse represents calculated pricing information
type PriceResponse struct {
	ProductName   string          `json:"product_name"`
	BasePrice     decimal.Decimal `json:"base_price"`
	Quantity      int             `json:"quantity"`
	DiscountPct   decimal.Decimal `json:"discount_pct"`
	Subtotal      decimal.Decimal `json:"subtotal"`
	DiscountAmt   decimal.Decimal `json:"discount_amount"`
	FinalTotal    decimal.Decimal `json:"final_total"`
	FormattedData FormattedPrice  `json:"formatted"`
}

// FormattedPrice contains human-readable string representations
type FormattedPrice struct {
	BasePrice   string `json:"base_price"`
	Subtotal    string `json:"subtotal"`
	DiscountAmt string `json:"discount_amount"`
	FinalTotal  string `json:"final_total"`
}

// CalculatePrice demonstrates the complete workflow:
// 1. Accept string input (API layer)
// 2. Validate using xvalidator
// 3. Convert to decimal.Decimal for calculations (business logic layer)
// 4. Return formatted results
func CalculatePrice(req PriceRequest, v *xvalidator.Validator) (*PriceResponse, error) {
	// Step 1: Validate input using string validators
	if err := v.StructTranslated(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Step 2: Convert validated strings to decimal.Decimal for accurate calculations
	basePrice, err := decimal.NewFromString(req.BasePrice)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base_price: %w", err)
	}

	discountPct, err := decimal.NewFromString(req.DiscountPct)
	if err != nil {
		return nil, fmt.Errorf("failed to parse discount_pct: %w", err)
	}

	// Step 3: Perform accurate decimal calculations
	quantity := decimal.NewFromInt(int64(req.Quantity))

	// Calculate subtotal: base_price * quantity
	subtotal := basePrice.Mul(quantity)

	// Calculate discount amount: subtotal * (discount_pct / 100)
	discountAmt := subtotal.Mul(discountPct).Div(decimal.NewFromInt(100))

	// Calculate final total: subtotal - discount_amount
	finalTotal := subtotal.Sub(discountAmt)

	// Step 4: Return response with both decimal and formatted string values
	return &PriceResponse{
		ProductName: req.ProductName,
		BasePrice:   basePrice,
		Quantity:    req.Quantity,
		DiscountPct: discountPct,
		Subtotal:    subtotal,
		DiscountAmt: discountAmt,
		FinalTotal:  finalTotal,
		FormattedData: FormattedPrice{
			BasePrice:   basePrice.StringFixed(2),
			Subtotal:    subtotal.StringFixed(2),
			DiscountAmt: discountAmt.StringFixed(2),
			FinalTotal:  finalTotal.StringFixed(2),
		},
	}, nil
}

func main() {
	fmt.Println("=== Best Practice: String Validation + Decimal Calculation Workflow ===\n")

	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Example 1: Valid calculation with discount
	fmt.Println("Example 1: Valid Price Calculation with 15% Discount")
	req1 := PriceRequest{
		ProductName: "Gaming Laptop",
		BasePrice:   "45900.00",
		Quantity:    2,
		DiscountPct: "15.00",
	}

	resp1, err := CalculatePrice(req1, v)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Calculation completed!")
		fmt.Printf("   Product: %s\n", resp1.ProductName)
		fmt.Printf("   Base Price: %s THB x %d\n", resp1.FormattedData.BasePrice, resp1.Quantity)
		fmt.Printf("   Subtotal: %s THB\n", resp1.FormattedData.Subtotal)
		fmt.Printf("   Discount (%.2f%%): -%s THB\n", resp1.DiscountPct, resp1.FormattedData.DiscountAmt)
		fmt.Printf("   Final Total: %s THB\n", resp1.FormattedData.FinalTotal)
	}
	fmt.Println()

	// Example 2: Valid calculation without discount
	fmt.Println("Example 2: Valid Price Calculation with No Discount")
	req2 := PriceRequest{
		ProductName: "Wireless Mouse",
		BasePrice:   "890.00",
		Quantity:    3,
		DiscountPct: "0.00",
	}

	resp2, err := CalculatePrice(req2, v)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Calculation completed!")
		fmt.Printf("   Product: %s\n", resp2.ProductName)
		fmt.Printf("   Base Price: %s THB x %d\n", resp2.FormattedData.BasePrice, resp2.Quantity)
		fmt.Printf("   Final Total: %s THB\n", resp2.FormattedData.FinalTotal)
	}
	fmt.Println()

	// Example 3: Valid calculation with maximum discount
	fmt.Println("Example 3: Valid Price Calculation with 100% Discount (Free)")
	req3 := PriceRequest{
		ProductName: "Promotional Item",
		BasePrice:   "1500.00",
		Quantity:    1,
		DiscountPct: "100.00",
	}

	resp3, err := CalculatePrice(req3, v)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Calculation completed!")
		fmt.Printf("   Product: %s\n", resp3.ProductName)
		fmt.Printf("   Subtotal: %s THB\n", resp3.FormattedData.Subtotal)
		fmt.Printf("   Discount (100%%): -%s THB\n", resp3.FormattedData.DiscountAmt)
		fmt.Printf("   Final Total: %s THB (FREE!)\n", resp3.FormattedData.FinalTotal)
	}
	fmt.Println()

	// Example 4: Invalid - negative price
	fmt.Println("Example 4: Invalid - Negative Base Price")
	req4 := PriceRequest{
		ProductName: "Test Product",
		BasePrice:   "-100.00",
		Quantity:    1,
		DiscountPct: "0.00",
	}

	resp4, err := CalculatePrice(req4, v)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Result: %s THB\n", resp4.FormattedData.FinalTotal)
	}
	fmt.Println()

	// Example 5: Invalid - discount over 100%
	fmt.Println("Example 5: Invalid - Discount Exceeds 100%")
	req5 := PriceRequest{
		ProductName: "Invalid Discount",
		BasePrice:   "1000.00",
		Quantity:    1,
		DiscountPct: "150.00",
	}

	resp5, err := CalculatePrice(req5, v)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Result: %s THB\n", resp5.FormattedData.FinalTotal)
	}
	fmt.Println()

	// Example 6: Invalid - invalid decimal format
	fmt.Println("Example 6: Invalid - Too Many Decimal Places")
	req6 := PriceRequest{
		ProductName: "Precise Product",
		BasePrice:   "99.999", // More than 2 decimal places
		Quantity:    1,
		DiscountPct: "5.00",
	}

	resp6, err := CalculatePrice(req6, v)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Result: %s THB\n", resp6.FormattedData.FinalTotal)
	}
	fmt.Println()

	// Example 7: Complex calculation demonstrating precision
	fmt.Println("Example 7: Precision Test - Complex Calculation")
	fmt.Println("This demonstrates why decimal.Decimal is essential for financial calculations")
	req7 := PriceRequest{
		ProductName: "High-Precision Item",
		BasePrice:   "33.33",
		Quantity:    3,
		DiscountPct: "7.77",
	}

	resp7, err := CalculatePrice(req7, v)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Calculation completed with precision!")
		fmt.Printf("   33.33 x 3 = %s (no rounding errors)\n", resp7.Subtotal.StringFixed(2))
		fmt.Printf("   Discount 7.77%% = %s\n", resp7.FormattedData.DiscountAmt)
		fmt.Printf("   Final Total: %s THB (accurate to 2 decimals)\n", resp7.FormattedData.FinalTotal)

		// Show what would happen with float64
		float64BasePrice := 33.33
		float64Subtotal := float64BasePrice * 3
		float64Discount := float64Subtotal * 0.0777
		float64Total := float64Subtotal - float64Discount
		fmt.Printf("\n   ⚠️  Using float64 would give: %.2f THB\n", float64Total)
		fmt.Printf("   ✅ Using decimal.Decimal gives: %s THB\n", resp7.FormattedData.FinalTotal)
		fmt.Println("   Note: Results may appear similar when formatted, but internal precision differs")
	}

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("1. ✅ Accept string input from API/users")
	fmt.Println("2. ✅ Validate using xvalidator's decimal validators (dgt, dgte, etc.)")
	fmt.Println("3. ✅ Convert to decimal.Decimal for calculations")
	fmt.Println("4. ✅ Perform all monetary calculations using decimal.Decimal")
	fmt.Println("5. ✅ Format output as string for display/API response")
	fmt.Println("\n❌ Never use float64 for monetary calculations!")
	fmt.Println("❌ Never pass decimal.Decimal directly in API structs!")
}
