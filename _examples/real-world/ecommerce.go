package main

import (
	"fmt"

	"github.com/hotfixfirst/go-xvalidator"
)

// Address for shipping and billing
type Address struct {
	Street     string `json:"street" validate:"required,min=5,max=200"`
	City       string `json:"city" validate:"required,min=2,max=100"`
	State      string `json:"state" validate:"required,len=2"`
	PostalCode string `json:"postal_code" validate:"required,len=5"`
	Country    string `json:"country" validate:"required,len=2"`
	Phone      string `json:"phone" validate:"required,mobile_e164"`
}

// CartItem represents a single item in shopping cart
type CartItem struct {
	ProductID   string `json:"product_id" validate:"required"`
	ProductName string `json:"product_name" validate:"required,min=3,max=200"`
	Quantity    int    `json:"quantity" validate:"required,min=1,max=100"`
	UnitPrice   string `json:"unit_price" validate:"required,decimal=10:2,dgt=0"`
	Subtotal    string `json:"subtotal" validate:"required,decimal=10:2,dgt=0"`
}

// EcommerceOrder represents complete e-commerce order
type EcommerceOrder struct {
	OrderID         string     `json:"order_id" validate:"required"`
	CustomerEmail   string     `json:"customer_email" validate:"required,email"`
	CustomerPhone   string     `json:"customer_phone" validate:"required,mobile_e164"`
	ShippingAddress Address    `json:"shipping_address" validate:"required"`
	BillingAddress  *Address   `json:"billing_address" validate:"omitempty"`
	Items           []CartItem `json:"items" validate:"required,min=1,dive"`
	Subtotal        string     `json:"subtotal" validate:"required,decimal=10:2,dgte=0"`
	ShippingFee     string     `json:"shipping_fee" validate:"required,decimal=10:2,dgte=0"`
	Tax             string     `json:"tax" validate:"required,decimal=10:2,dgte=0"`
	Discount        string     `json:"discount" validate:"required,decimal=10:2,dgte=0"`
	Total           string     `json:"total" validate:"required,decimal=10:2,dgt=0"`
	PaymentMethod   string     `json:"payment_method" validate:"required,oneof=credit_card debit_card bank_transfer ewallet cod"`
	Notes           string     `json:"notes" validate:"omitempty,max=500"`
}

func main() {
	fmt.Println("=== E-Commerce Order Validation Example ===\n")

	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Example 1: Valid complete order
	fmt.Println("Example 1: Valid Complete E-Commerce Order")
	validOrder := EcommerceOrder{
		OrderID:       "ORD-2024-001",
		CustomerEmail: "customer@example.com",
		CustomerPhone: "+66812345678",
		ShippingAddress: Address{
			Street:     "123 Main Street",
			City:       "Bangkok",
			State:      "BK",
			PostalCode: "10110",
			Country:    "TH",
			Phone:      "+66812345678",
		},
		Items: []CartItem{
			{
				ProductID:   "PROD-001",
				ProductName: "Gaming Laptop",
				Quantity:    1,
				UnitPrice:   "45900.00",
				Subtotal:    "45900.00",
			},
			{
				ProductID:   "PROD-002",
				ProductName: "Wireless Mouse",
				Quantity:    2,
				UnitPrice:   "890.00",
				Subtotal:    "1780.00",
			},
		},
		Subtotal:      "47680.00",
		ShippingFee:   "200.00",
		Tax:           "3337.60",
		Discount:      "500.00",
		Total:         "50717.60",
		PaymentMethod: "credit_card",
		Notes:         "Please deliver between 9 AM - 5 PM",
	}

	if err := v.StructTranslated(validOrder); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
		fmt.Printf("   Order ID: %s\n", validOrder.OrderID)
		fmt.Printf("   Items: %d\n", len(validOrder.Items))
		fmt.Printf("   Total: %s THB\n", validOrder.Total)
	}
	fmt.Println()

	// Example 2: Valid order with separate billing address
	fmt.Println("Example 2: Valid Order with Separate Billing Address")
	billingAddr := Address{
		Street:     "456 Office Building",
		City:       "Bangkok",
		State:      "BK",
		PostalCode: "10120",
		Country:    "TH",
		Phone:      "+66812345678",
	}
	orderWithBilling := EcommerceOrder{
		OrderID:       "ORD-2024-002",
		CustomerEmail: "business@example.com",
		CustomerPhone: "+66898765432",
		ShippingAddress: Address{
			Street:     "789 Home Address",
			City:       "Bangkok",
			State:      "BK",
			PostalCode: "10150",
			Country:    "TH",
			Phone:      "+66898765432",
		},
		BillingAddress: &billingAddr,
		Items: []CartItem{
			{
				ProductID:   "PROD-003",
				ProductName: "Office Chair",
				Quantity:    3,
				UnitPrice:   "4500.00",
				Subtotal:    "13500.00",
			},
		},
		Subtotal:      "13500.00",
		ShippingFee:   "500.00",
		Tax:           "980.00",
		Discount:      "0.00",
		Total:         "14980.00",
		PaymentMethod: "bank_transfer",
	}

	if err := v.StructTranslated(orderWithBilling); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
		fmt.Printf("   Separate billing address provided\n")
	}
	fmt.Println()

	// Example 3: Invalid - missing required item field
	fmt.Println("Example 3: Invalid Order - Missing Product Name")
	invalidOrder1 := EcommerceOrder{
		OrderID:       "ORD-2024-003",
		CustomerEmail: "test@example.com",
		CustomerPhone: "+66812345678",
		ShippingAddress: Address{
			Street:     "111 Test Street",
			City:       "Bangkok",
			State:      "BK",
			PostalCode: "10100",
			Country:    "TH",
			Phone:      "+66812345678",
		},
		Items: []CartItem{
			{
				ProductID:   "PROD-004",
				ProductName: "", // Missing
				Quantity:    1,
				UnitPrice:   "1000.00",
				Subtotal:    "1000.00",
			},
		},
		Subtotal:      "1000.00",
		ShippingFee:   "100.00",
		Tax:           "77.00",
		Discount:      "0.00",
		Total:         "1177.00",
		PaymentMethod: "credit_card",
	}

	if err := v.StructTranslated(invalidOrder1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 4: Invalid - negative price
	fmt.Println("Example 4: Invalid Order - Negative Unit Price")
	invalidOrder2 := EcommerceOrder{
		OrderID:       "ORD-2024-004",
		CustomerEmail: "test2@example.com",
		CustomerPhone: "+66812345678",
		ShippingAddress: Address{
			Street:     "222 Test Avenue",
			City:       "Bangkok",
			State:      "BK",
			PostalCode: "10200",
			Country:    "TH",
			Phone:      "+66812345678",
		},
		Items: []CartItem{
			{
				ProductID:   "PROD-005",
				ProductName: "Test Product",
				Quantity:    1,
				UnitPrice:   "-100.00", // Invalid negative
				Subtotal:    "-100.00",
			},
		},
		Subtotal:      "-100.00",
		ShippingFee:   "50.00",
		Tax:           "0.00",
		Discount:      "0.00",
		Total:         "-50.00",
		PaymentMethod: "cod",
	}

	if err := v.StructTranslated(invalidOrder2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 5: Invalid - invalid email
	fmt.Println("Example 5: Invalid Order - Invalid Email Format")
	invalidOrder3 := EcommerceOrder{
		OrderID:       "ORD-2024-005",
		CustomerEmail: "not-an-email", // Invalid
		CustomerPhone: "+66812345678",
		ShippingAddress: Address{
			Street:     "333 Test Road",
			City:       "Bangkok",
			State:      "BK",
			PostalCode: "10300",
			Country:    "TH",
			Phone:      "+66812345678",
		},
		Items: []CartItem{
			{
				ProductID:   "PROD-006",
				ProductName: "Sample Product",
				Quantity:    1,
				UnitPrice:   "500.00",
				Subtotal:    "500.00",
			},
		},
		Subtotal:      "500.00",
		ShippingFee:   "50.00",
		Tax:           "38.50",
		Discount:      "0.00",
		Total:         "588.50",
		PaymentMethod: "ewallet",
	}

	if err := v.StructTranslated(invalidOrder3); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 6: Invalid - invalid phone number (not E.164)
	fmt.Println("Example 6: Invalid Order - Invalid Phone Number")
	invalidOrder4 := EcommerceOrder{
		OrderID:       "ORD-2024-006",
		CustomerEmail: "valid@example.com",
		CustomerPhone: "0812345678", // Invalid - not E.164 format
		ShippingAddress: Address{
			Street:     "444 Test Lane",
			City:       "Bangkok",
			State:      "BK",
			PostalCode: "10400",
			Country:    "TH",
			Phone:      "+66812345678",
		},
		Items: []CartItem{
			{
				ProductID:   "PROD-007",
				ProductName: "Another Product",
				Quantity:    1,
				UnitPrice:   "750.00",
				Subtotal:    "750.00",
			},
		},
		Subtotal:      "750.00",
		ShippingFee:   "100.00",
		Tax:           "59.50",
		Discount:      "0.00",
		Total:         "909.50",
		PaymentMethod: "credit_card",
	}

	if err := v.StructTranslated(invalidOrder4); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 7: Invalid - empty items array
	fmt.Println("Example 7: Invalid Order - No Items")
	invalidOrder5 := EcommerceOrder{
		OrderID:       "ORD-2024-007",
		CustomerEmail: "empty@example.com",
		CustomerPhone: "+66812345678",
		ShippingAddress: Address{
			Street:     "555 Empty Cart St",
			City:       "Bangkok",
			State:      "BK",
			PostalCode: "10500",
			Country:    "TH",
			Phone:      "+66812345678",
		},
		Items:         []CartItem{}, // Empty
		Subtotal:      "0.00",
		ShippingFee:   "0.00",
		Tax:           "0.00",
		Discount:      "0.00",
		Total:         "0.00",
		PaymentMethod: "credit_card",
	}

	if err := v.StructTranslated(invalidOrder5); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}

	fmt.Println("\n=== E-Commerce Order Validation Complete ===")
}
