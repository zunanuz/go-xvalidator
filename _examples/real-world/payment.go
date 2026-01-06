package main

import (
	"fmt"

	"github.com/hotfixfirst/go-xvalidator"
)

// PaymentRequest represents a payment processing request
type PaymentRequest struct {
	// Transaction Details
	TransactionID string `json:"transaction_id" validate:"required"`
	OrderID       string `json:"order_id" validate:"required"`
	CustomerID    string `json:"customer_id" validate:"required"`

	// Amount Details
	Amount   string `json:"amount" validate:"required,decimal=10:2,dgt=0,dlte=1000000"`
	Currency string `json:"currency" validate:"required,len=3"`
	Tax      string `json:"tax" validate:"required,decimal=10:2,dgte=0"`
	Fee      string `json:"fee" validate:"required,decimal=10:2,dgte=0"`
	Total    string `json:"total" validate:"required,decimal=10:2,dgt=0"`

	// Payment Method
	PaymentMethod string       `json:"payment_method" validate:"required,oneof=credit_card debit_card bank_transfer ewallet"`
	CardDetails   *CardDetails `json:"card_details" validate:"required_if=PaymentMethod credit_card,required_if=PaymentMethod debit_card"`

	// URLs for callbacks
	SuccessURL string `json:"success_url" validate:"required,https_url"`
	FailureURL string `json:"failure_url" validate:"required,https_url"`
	WebhookURL string `json:"webhook_url" validate:"required,https_url"`

	// Additional Info
	Description string   `json:"description" validate:"required,min=10,max=500"`
	Items       []string `json:"items" validate:"omitempty,dive,required"`
}

// CardDetails for credit/debit card payments
type CardDetails struct {
	CardholderName string `json:"cardholder_name" validate:"required,min=3,max=100"`
	CardNumber     string `json:"card_number" validate:"required,len=16"`
	ExpiryMonth    string `json:"expiry_month" validate:"required,len=2"`
	ExpiryYear     string `json:"expiry_year" validate:"required,len=2"`
	CVV            string `json:"cvv" validate:"required,len=3"`
	BillingZip     string `json:"billing_zip" validate:"required,len=5"`
}

func main() {
	fmt.Println("=== Payment Processing Validation Example ===\n")

	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Example 1: Valid credit card payment
	fmt.Println("Example 1: Valid Credit Card Payment")
	validPayment := PaymentRequest{
		TransactionID: "TXN-2024-001",
		OrderID:       "ORD-2024-001",
		CustomerID:    "CUST-123",
		Amount:        "99.99",
		Currency:      "USD",
		Tax:           "7.00",
		Fee:           "2.50",
		Total:         "109.49",
		PaymentMethod: "credit_card",
		CardDetails: &CardDetails{
			CardholderName: "John Doe",
			CardNumber:     "4532123456789012",
			ExpiryMonth:    "12",
			ExpiryYear:     "25",
			CVV:            "123",
			BillingZip:     "10110",
		},
		SuccessURL:  "https://shop.example.com/payment/success",
		FailureURL:  "https://shop.example.com/payment/failure",
		WebhookURL:  "https://shop.example.com/webhooks/payment",
		Description: "Payment for Order ORD-2024-001 containing 2 items",
		Items:       []string{"PROD-001", "PROD-002"},
	}

	if err := v.StructTranslated(validPayment); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Payment validation passed!")
		fmt.Printf("   Transaction: %s\n", validPayment.TransactionID)
		fmt.Printf("   Amount: %s %s\n", validPayment.Currency, validPayment.Total)
		fmt.Printf("   Method: %s\n", validPayment.PaymentMethod)
	}
	fmt.Println()

	// Example 2: Valid e-wallet payment (no card details needed)
	fmt.Println("Example 2: Valid E-Wallet Payment")
	validEwalletPayment := PaymentRequest{
		TransactionID: "TXN-2024-002",
		OrderID:       "ORD-2024-002",
		CustomerID:    "CUST-456",
		Amount:        "49.99",
		Currency:      "THB",
		Tax:           "3.50",
		Fee:           "1.00",
		Total:         "54.49",
		PaymentMethod: "ewallet",
		CardDetails:   nil, // Not required for ewallet
		SuccessURL:    "https://shop.example.com/payment/success",
		FailureURL:    "https://shop.example.com/payment/failure",
		WebhookURL:    "https://shop.example.com/webhooks/payment",
		Description:   "Payment via e-wallet for digital product",
		Items:         []string{"PROD-DIGITAL-001"},
	}

	if err := v.StructTranslated(validEwalletPayment); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ E-wallet payment validated!")
		fmt.Printf("   Total: %s %s\n", validEwalletPayment.Currency, validEwalletPayment.Total)
	}
	fmt.Println()

	// Example 3: Invalid - amount exceeds maximum
	fmt.Println("Example 3: Invalid - Amount Exceeds Maximum")
	invalidPayment1 := PaymentRequest{
		TransactionID: "TXN-2024-003",
		OrderID:       "ORD-2024-003",
		CustomerID:    "CUST-789",
		Amount:        "2000000.00", // Exceeds max of 1,000,000
		Currency:      "USD",
		Tax:           "140000.00",
		Fee:           "10000.00",
		Total:         "2150000.00",
		PaymentMethod: "credit_card",
		CardDetails: &CardDetails{
			CardholderName: "Jane Smith",
			CardNumber:     "4532987654321098",
			ExpiryMonth:    "06",
			ExpiryYear:     "26",
			CVV:            "456",
			BillingZip:     "10120",
		},
		SuccessURL:  "https://shop.example.com/payment/success",
		FailureURL:  "https://shop.example.com/payment/failure",
		WebhookURL:  "https://shop.example.com/webhooks/payment",
		Description: "Large payment for bulk order",
	}

	if err := v.StructTranslated(invalidPayment1); err != nil {
		fmt.Printf("❌ Validation failed:\n")
		fmt.Printf("   %v\n", err)
	}
	fmt.Println()

	// Example 4: Invalid - non-HTTPS callback URLs
	fmt.Println("Example 4: Invalid - Non-HTTPS Callback URLs")
	invalidPayment2 := PaymentRequest{
		TransactionID: "TXN-2024-004",
		OrderID:       "ORD-2024-004",
		CustomerID:    "CUST-999",
		Amount:        "75.00",
		Currency:      "USD",
		Tax:           "5.25",
		Fee:           "1.50",
		Total:         "81.75",
		PaymentMethod: "credit_card",
		CardDetails: &CardDetails{
			CardholderName: "Bob Johnson",
			CardNumber:     "4532111122223333",
			ExpiryMonth:    "03",
			ExpiryYear:     "27",
			CVV:            "789",
			BillingZip:     "10130",
		},
		SuccessURL:  "http://shop.example.com/payment/success",  // HTTP not allowed
		FailureURL:  "http://shop.example.com/payment/failure",  // HTTP not allowed
		WebhookURL:  "http://shop.example.com/webhooks/payment", // HTTP not allowed
		Description: "Payment with insecure callback URLs",
	}

	if err := v.StructTranslated(invalidPayment2); err != nil {
		fmt.Printf("❌ Validation failed:\n")
		fmt.Printf("   %v\n", err)
	}
	fmt.Println()

	// Example 5: Invalid - missing card details for card payment
	fmt.Println("Example 5: Invalid - Missing Card Details for Card Payment")
	invalidPayment3 := PaymentRequest{
		TransactionID: "TXN-2024-005",
		OrderID:       "ORD-2024-005",
		CustomerID:    "CUST-111",
		Amount:        "125.00",
		Currency:      "USD",
		Tax:           "8.75",
		Fee:           "2.50",
		Total:         "136.25",
		PaymentMethod: "credit_card",
		CardDetails:   nil, // Required for credit_card
		SuccessURL:    "https://shop.example.com/payment/success",
		FailureURL:    "https://shop.example.com/payment/failure",
		WebhookURL:    "https://shop.example.com/webhooks/payment",
		Description:   "Payment without required card details",
	}

	if err := v.StructTranslated(invalidPayment3); err != nil {
		fmt.Printf("❌ Validation failed:\n")
		fmt.Printf("   %v\n", err)
	}
	fmt.Println()

	// Example 6: Invalid - invalid card details
	fmt.Println("Example 6: Invalid - Invalid Card Number Length")
	invalidPayment4 := PaymentRequest{
		TransactionID: "TXN-2024-006",
		OrderID:       "ORD-2024-006",
		CustomerID:    "CUST-222",
		Amount:        "59.99",
		Currency:      "USD",
		Tax:           "4.20",
		Fee:           "1.50",
		Total:         "65.69",
		PaymentMethod: "credit_card",
		CardDetails: &CardDetails{
			CardholderName: "Alice Brown",
			CardNumber:     "45321234", // Too short (must be 16 digits)
			ExpiryMonth:    "08",
			ExpiryYear:     "25",
			CVV:            "123",
			BillingZip:     "10140",
		},
		SuccessURL:  "https://shop.example.com/payment/success",
		FailureURL:  "https://shop.example.com/payment/failure",
		WebhookURL:  "https://shop.example.com/webhooks/payment",
		Description: "Payment with invalid card number",
	}

	if err := v.StructTranslated(invalidPayment4); err != nil {
		fmt.Printf("❌ Validation failed:\n")
		fmt.Printf("   %v\n", err)
	}
	fmt.Println()

	// Example 7: Invalid - negative amount
	fmt.Println("Example 7: Invalid - Negative Amount")
	invalidPayment5 := PaymentRequest{
		TransactionID: "TXN-2024-007",
		OrderID:       "ORD-2024-007",
		CustomerID:    "CUST-333",
		Amount:        "-50.00", // Negative not allowed
		Currency:      "USD",
		Tax:           "0.00",
		Fee:           "0.00",
		Total:         "-50.00",
		PaymentMethod: "bank_transfer",
		SuccessURL:    "https://shop.example.com/payment/success",
		FailureURL:    "https://shop.example.com/payment/failure",
		WebhookURL:    "https://shop.example.com/webhooks/payment",
		Description:   "Refund or negative payment attempt",
	}

	if err := v.StructTranslated(invalidPayment5); err != nil {
		fmt.Printf("❌ Validation failed:\n")
		fmt.Printf("   %v\n", err)
	}
	fmt.Println()

	fmt.Println("=== Payment Processing Validation Complete ===")
}
