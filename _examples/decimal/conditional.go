package main

import (
	"fmt"

	"github.com/hotfixfirst/go-xvalidator"
)

// PaymentMethod demonstrates conditional decimal validation
type PaymentMethod struct {
	Type          string `json:"type" validate:"required,oneof=credit debit cash"`
	Amount        string `json:"amount" validate:"required,decimal_if=2@Type=credit,decimal_if=2@Type=debit,decimal_if=0@Type=cash"`
	MinAmount     string `json:"min_amount" validate:"decimal_if=10:2@Type=credit"`
	ProcessingFee string `json:"processing_fee" validate:"decimal_if=5:2@Type=credit,decimal_if=5:2@Type=debit"`
}

// ProductPricing demonstrates conditional decimal based on product type
type ProductPricing struct {
	ProductType string `json:"product_type" validate:"required,oneof=digital physical"`
	Price       string `json:"price" validate:"required,decimal_if=10:2@ProductType=digital,decimal_if=10:2@ProductType=physical"`
	ShippingFee string `json:"shipping_fee" validate:"decimal_if=0@ProductType=digital,decimal_if=8:2@ProductType=physical"`
	TaxRate     string `json:"tax_rate" validate:"decimal_if=5:2@ProductType=digital,decimal_if=5:2@ProductType=physical"`
}

// OrderStatus demonstrates validation based on status
type OrderStatus struct {
	Status       string `json:"status" validate:"required,oneof=pending paid shipped delivered"`
	PaidAmount   string `json:"paid_amount" validate:"decimal_if=10:2@Status=paid,decimal_if=10:2@Status=shipped,decimal_if=10:2@Status=delivered"`
	RefundAmount string `json:"refund_amount" validate:"decimal_if=10:2@Status=delivered"`
}

func main() {
	fmt.Println("=== Conditional Decimal Validation Examples ===\n")

	// Create validator instance
	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Example 1: Valid credit card payment
	fmt.Println("Example 1: Valid Credit Card Payment")
	creditPayment := PaymentMethod{
		Type:          "credit",
		Amount:        "100.50",
		MinAmount:     "50.00",
		ProcessingFee: "2.50",
	}
	if err := v.Validate(creditPayment); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 2: Valid debit payment
	fmt.Println("Example 2: Valid Debit Payment")
	debitPayment := PaymentMethod{
		Type:          "debit",
		Amount:        "75.25",
		MinAmount:     "",
		ProcessingFee: "1.50",
	}
	if err := v.Validate(debitPayment); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 3: Valid cash payment (integer only)
	fmt.Println("Example 3: Valid Cash Payment (Integer Only)")
	cashPayment := PaymentMethod{
		Type:          "cash",
		Amount:        "100",
		MinAmount:     "",
		ProcessingFee: "",
	}
	if err := v.Validate(cashPayment); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 4: Invalid cash payment (has decimal)
	fmt.Println("Example 4: Invalid Cash Payment (Has Decimal)")
	invalidCashPayment := PaymentMethod{
		Type:          "cash",
		Amount:        "100.50",
		MinAmount:     "",
		ProcessingFee: "",
	}
	if err := v.Validate(invalidCashPayment); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 5: Valid digital product
	fmt.Println("Example 5: Valid Digital Product")
	digitalProduct := ProductPricing{
		ProductType: "digital",
		Price:       "29.99",
		ShippingFee: "0",
		TaxRate:     "7.50",
	}
	if err := v.Validate(digitalProduct); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 6: Valid physical product
	fmt.Println("Example 6: Valid Physical Product")
	physicalProduct := ProductPricing{
		ProductType: "physical",
		Price:       "49.99",
		ShippingFee: "15.00",
		TaxRate:     "7.50",
	}
	if err := v.Validate(physicalProduct); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 7: Invalid digital product (shipping fee not zero)
	fmt.Println("Example 7: Invalid Digital Product (Shipping Fee Not Zero)")
	invalidDigitalProduct := ProductPricing{
		ProductType: "digital",
		Price:       "29.99",
		ShippingFee: "10.00",
		TaxRate:     "7.50",
	}
	if err := v.Validate(invalidDigitalProduct); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 8: Valid pending order
	fmt.Println("Example 8: Valid Pending Order")
	pendingOrder := OrderStatus{
		Status:       "pending",
		PaidAmount:   "",
		RefundAmount: "",
	}
	if err := v.Validate(pendingOrder); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 9: Valid paid order
	fmt.Println("Example 9: Valid Paid Order")
	paidOrder := OrderStatus{
		Status:       "paid",
		PaidAmount:   "150.75",
		RefundAmount: "",
	}
	if err := v.Validate(paidOrder); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 10: Valid delivered order with refund
	fmt.Println("Example 10: Valid Delivered Order with Refund")
	deliveredOrder := OrderStatus{
		Status:       "delivered",
		PaidAmount:   "200.00",
		RefundAmount: "50.00",
	}
	if err := v.Validate(deliveredOrder); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
}
