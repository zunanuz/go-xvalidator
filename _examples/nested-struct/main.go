package main

import (
	"fmt"

	"github.com/hotfixfirst/go-xvalidator"
)

// Address demonstrates nested struct validation
type Address struct {
	Street     string `json:"street" validate:"required,min=5,max=200"`
	City       string `json:"city" validate:"required,min=2,max=100"`
	State      string `json:"state" validate:"required,len=2"`
	PostalCode string `json:"postal_code" validate:"required,len=5"`
	Country    string `json:"country" validate:"required,len=2"`
}

// Person demonstrates nested struct with pointer
type Person struct {
	Name    string   `json:"name" validate:"required,min=2,max=100"`
	Email   string   `json:"email" validate:"required,email"`
	Age     int      `json:"age" validate:"required,min=18,max=120"`
	Address Address  `json:"address" validate:"required"`
	Billing *Address `json:"billing" validate:"omitempty"`
}

// Company demonstrates deeply nested structs
type Company struct {
	Name      string   `json:"name" validate:"required,min=3,max=200"`
	TaxID     string   `json:"tax_id" validate:"required,len=10"`
	Address   Address  `json:"address" validate:"required"`
	CEO       Person   `json:"ceo" validate:"required"`
	Employees []Person `json:"employees" validate:"omitempty,dive"`
}

// Order demonstrates nested struct with arrays
type Order struct {
	OrderID     string   `json:"order_id" validate:"required"`
	Customer    Person   `json:"customer" validate:"required"`
	ShipTo      Address  `json:"ship_to" validate:"required"`
	BillTo      *Address `json:"bill_to" validate:"omitempty"`
	TotalAmount string   `json:"total_amount" validate:"required,decimal=10:2,dgte=0"`
}

func main() {
	fmt.Println("=== Nested Struct Validation Examples ===\n")

	// Create validator instance
	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Example 1: Valid person with address
	fmt.Println("Example 1: Valid Person with Address")
	person1 := Person{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
		Address: Address{
			Street:     "123 Main St",
			City:       "New York",
			State:      "NY",
			PostalCode: "10001",
			Country:    "US",
		},
	}
	if err := v.Validate(person1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 2: Valid person with billing address
	fmt.Println("Example 2: Valid Person with Billing Address")
	billingAddr := Address{
		Street:     "456 Oak Ave",
		City:       "Los Angeles",
		State:      "CA",
		PostalCode: "90001",
		Country:    "US",
	}
	person2 := Person{
		Name:  "Jane Smith",
		Email: "jane@example.com",
		Age:   28,
		Address: Address{
			Street:     "789 Pine Rd",
			City:       "Chicago",
			State:      "IL",
			PostalCode: "60601",
			Country:    "US",
		},
		Billing: &billingAddr,
	}
	if err := v.Validate(person2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 3: Invalid - missing required nested field
	fmt.Println("Example 3: Invalid - Missing Required Nested Field")
	invalid1 := Person{
		Name:  "Bob Wilson",
		Email: "bob@example.com",
		Age:   35,
		Address: Address{
			Street:     "",
			City:       "Seattle",
			State:      "WA",
			PostalCode: "98101",
			Country:    "US",
		},
	}
	if err := v.Validate(invalid1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 4: Invalid - nested field validation failed
	fmt.Println("Example 4: Invalid - Nested Field Validation Failed")
	invalid2 := Person{
		Name:  "Alice Brown",
		Email: "alice@example.com",
		Age:   25,
		Address: Address{
			Street:     "321 Elm St",
			City:       "Boston",
			State:      "MA",
			PostalCode: "021", // Too short
			Country:    "US",
		},
	}
	if err := v.Validate(invalid2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 5: Valid company with CEO
	fmt.Println("Example 5: Valid Company with CEO")
	company1 := Company{
		Name:  "Tech Corp",
		TaxID: "1234567890",
		Address: Address{
			Street:     "1000 Corporate Blvd",
			City:       "San Francisco",
			State:      "CA",
			PostalCode: "94105",
			Country:    "US",
		},
		CEO: Person{
			Name:  "Steve Jobs",
			Email: "steve@techcorp.com",
			Age:   45,
			Address: Address{
				Street:     "2000 CEO Lane",
				City:       "Palo Alto",
				State:      "CA",
				PostalCode: "94301",
				Country:    "US",
			},
		},
	}
	if err := v.Validate(company1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 6: Valid company with employees
	fmt.Println("Example 6: Valid Company with Employees")
	company2 := Company{
		Name:  "Startup Inc",
		TaxID: "9876543210",
		Address: Address{
			Street:     "500 Innovation Dr",
			City:       "Austin",
			State:      "TX",
			PostalCode: "73301",
			Country:    "US",
		},
		CEO: Person{
			Name:  "Elon Musk",
			Email: "elon@startup.com",
			Age:   50,
			Address: Address{
				Street:     "3000 Space Ave",
				City:       "Austin",
				State:      "TX",
				PostalCode: "73301",
				Country:    "US",
			},
		},
		Employees: []Person{
			{
				Name:  "Employee One",
				Email: "emp1@startup.com",
				Age:   30,
				Address: Address{
					Street:     "100 Worker St",
					City:       "Austin",
					State:      "TX",
					PostalCode: "73301",
					Country:    "US",
				},
			},
			{
				Name:  "Employee Two",
				Email: "emp2@startup.com",
				Age:   28,
				Address: Address{
					Street:     "200 Employee Rd",
					City:       "Austin",
					State:      "TX",
					PostalCode: "73301",
					Country:    "US",
				},
			},
		},
	}
	if err := v.Validate(company2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 7: Invalid - employee validation failed
	fmt.Println("Example 7: Invalid - Employee Validation Failed")
	invalidCompany := Company{
		Name:  "Bad Corp",
		TaxID: "1111111111",
		Address: Address{
			Street:     "999 Error St",
			City:       "Denver",
			State:      "CO",
			PostalCode: "80201",
			Country:    "US",
		},
		CEO: Person{
			Name:  "CEO Person",
			Email: "ceo@badcorp.com",
			Age:   40,
			Address: Address{
				Street:     "800 Boss Blvd",
				City:       "Denver",
				State:      "CO",
				PostalCode: "80201",
				Country:    "US",
			},
		},
		Employees: []Person{
			{
				Name:  "Valid Employee",
				Email: "valid@badcorp.com",
				Age:   25,
				Address: Address{
					Street:     "100 Good St",
					City:       "Denver",
					State:      "CO",
					PostalCode: "80201",
					Country:    "US",
				},
			},
			{
				Name:  "Invalid Employee",
				Email: "invalid-email", // Invalid email
				Age:   22,
				Address: Address{
					Street:     "200 Bad St",
					City:       "Denver",
					State:      "CO",
					PostalCode: "80201",
					Country:    "US",
				},
			},
		},
	}
	if err := v.Validate(invalidCompany); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 8: Valid order
	fmt.Println("Example 8: Valid Order")
	order1 := Order{
		OrderID: "ORD-12345",
		Customer: Person{
			Name:  "Customer Name",
			Email: "customer@example.com",
			Age:   35,
			Address: Address{
				Street:     "700 Customer St",
				City:       "Miami",
				State:      "FL",
				PostalCode: "33101",
				Country:    "US",
			},
		},
		ShipTo: Address{
			Street:     "800 Shipping Ln",
			City:       "Miami",
			State:      "FL",
			PostalCode: "33102",
			Country:    "US",
		},
		TotalAmount: "299.99",
	}
	if err := v.Validate(order1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 9: Valid order with separate billing
	fmt.Println("Example 9: Valid Order with Separate Billing")
	billAddr := Address{
		Street:     "900 Billing Ave",
		City:       "Tampa",
		State:      "FL",
		PostalCode: "33601",
		Country:    "US",
	}
	order2 := Order{
		OrderID: "ORD-67890",
		Customer: Person{
			Name:  "Another Customer",
			Email: "another@example.com",
			Age:   42,
			Address: Address{
				Street:     "1100 Main Pkwy",
				City:       "Orlando",
				State:      "FL",
				PostalCode: "32801",
				Country:    "US",
			},
		},
		ShipTo: Address{
			Street:     "1200 Ship St",
			City:       "Orlando",
			State:      "FL",
			PostalCode: "32802",
			Country:    "US",
		},
		BillTo:      &billAddr,
		TotalAmount: "1299.50",
	}
	if err := v.Validate(order2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 10: Invalid - deeply nested validation error
	fmt.Println("Example 10: Invalid - Deeply Nested Validation Error")
	invalidOrder := Order{
		OrderID: "ORD-99999",
		Customer: Person{
			Name:  "Problem Customer",
			Email: "problem@example.com",
			Age:   30,
			Address: Address{
				Street:     "1300 Issue Rd",
				City:       "Houston",
				State:      "TX",
				PostalCode: "77001",
				Country:    "US",
			},
		},
		ShipTo: Address{
			Street:     "Short", // Too short
			City:       "Houston",
			State:      "TX",
			PostalCode: "77002",
			Country:    "US",
		},
		TotalAmount: "599.99",
	}
	if err := v.Validate(invalidOrder); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
}
