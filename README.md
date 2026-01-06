# go-xvalidator

A powerful extension package for go-playground/validator that provides enhanced validation capabilities including decimal validation, phone number validation (E.164 format), URL validation, password strength validation, and localized error messages for Go applications.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.25-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Features

- 🔢 **Decimal Validation**: Precise decimal comparison operators (greater than, less than, equal, etc.)
- 📱 **Phone Validation**: International phone number validation with country-specific rules (E.164 format)
- 🔗 **URL Validation**: Standard and HTTPS-only URL validation
- 🔒 **Password Strength**: Comprehensive password strength validation
- 🌍 **Localization**: Built-in English translations for validation errors
- ✨ **Easy Integration**: Drop-in replacement for go-playground/validator

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Available Validators](#available-validators)
  - [Decimal Validators](#decimal-validators)
  - [Conditional Decimal Validators](#conditional-decimal-validators)
  - [Phone Number Validators](#phone-number-validators)
  - [URL Validators](#url-validators)
  - [Password Strength Validator](#password-strength-validator)
- [Examples](#examples)
- [Best Practices](#best-practices)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Installation

```bash
go get github.com/hotfixfirst/go-xvalidator
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/hotfixfirst/go-xvalidator"
)

func main() {
    // Create validator instance
    v, err := xvalidator.NewValidator()
    if err != nil {
        panic(err)
    }

    // Define your struct with validation tags
    type User struct {
        Name     string `json:"name" validate:"required,min=3"`
        Email    string `json:"email" validate:"required,email"`
        Phone    string `json:"phone" validate:"mobile_e164=US"`
        Balance  string `json:"balance" validate:"dgte=0"`
        Password string `json:"password" validate:"password_strength"`
    }

    // Validate
    user := User{
        Name:     "John Doe",
        Email:    "john.doe@example.com",
        Phone:    "+447912345678",
        Balance:  "1000.50",
        Password: "SecureP@ss123",
    }

    if err := v.Validate(user); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
        return
    }

    fmt.Println("Validation passed!")
}
```

## Available Validators

### Decimal Validators

Validate decimal values with precise comparison:

```go
type Product struct {
    Price    string `validate:"required,dgte=0,dlte=999999.99"`
    Discount string `validate:"dgte=0,dlte=100"`
    Quantity string `validate:"dgt=0"`
}
```

**Tags:**

- `decimal` - Validates decimal format (precision and scale)
- `dgt=value` - Decimal greater than
- `dgte=value` - Decimal greater than or equal
- `dlt=value` - Decimal less than
- `dlte=value` - Decimal less than or equal
- `deq=value` - Decimal equal to
- `dneq=value` - Decimal not equal to

### Conditional Decimal Validators

Validate decimals conditionally based on other field values:

```go
type Payment struct {
    Method string `validate:"required,oneof=cash credit_card"`
    Amount string `validate:"required,decimal_if=Method cash dgte 0 decimal 10:0"`
    // Amount must be integer (0 decimals) if Method is "cash"
}
```

**Tag:**

- `decimal_if=FieldName FieldValue Operator CompareValue DecimalFormat`
  - Only validates if FieldName equals FieldValue
  - Supports: `dgt`, `dgte`, `dlt`, `dlte`, `deq`, `dneq`
  - DecimalFormat: `precision:scale` (e.g., `10:2`)

### Phone Number Validators

Validate international phone numbers in E.164 format:

```go
type Contact struct {
    PhoneGeneral string `validate:"mobile_e164"`        // Any country
    PhoneTH      string `validate:"mobile_e164=TH"`     // Thai mobile
    PhoneUS      string `validate:"mobile_e164=US"`     // US mobile
    PhoneGB      string `validate:"mobile_e164=GB"`     // UK mobile
}
```

**Supported Countries:**

- `TH` - Thailand
- `US` - United States
- `GB` - United Kingdom
- `FR` - France
- And many more...

### URL Validators

Validate URL formats:

```go
type Website struct {
    Homepage string `validate:"url"`              // Any valid URL
    SecureAPI string `validate:"https_url"`       // HTTPS only
}
```

### Password Strength Validator

Validate password complexity:

```go
type Account struct {
    Password string `validate:"password_strength"`
}

// Or use standalone validation
err := xvalidator.ValidatePasswordStrength("MyP@ssw0rd")
```

**Requirements:**

- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one digit
- At least one special character: `!@#$%^&*()_+-=[]{}|;:,.<>?`

## Examples

See the [_examples/](_examples/) directory for 131 comprehensive examples across 10 categories.

### 📚 Example Categories

#### Getting Started
- [**basic/**](_examples/basic/) - Basic validation (10 examples)
  - Required fields, email, min/max validation
  
#### Feature-Specific
- [**decimal/**](_examples/decimal/) - Decimal validation (22 examples)
  - Standard comparisons (dgt, dgte, dlt, dlte, deq, dneq)
  - Conditional validation (decimal_if)
  - Precision and scale validation
  
- [**phone/**](_examples/phone/) - E.164 phone validation (12 examples)
  - International formats (TH, US, UK, FR, DE, JP, etc.)
  
- [**url/**](_examples/url/) - URL validation (13 examples)
  - HTTP/HTTPS validation with paths and parameters
  
- [**password/**](_examples/password/) - Password strength (14 examples)
  - Complexity requirements and edge cases

#### Advanced Usage
- [**nested-struct/**](_examples/nested-struct/) - Nested structures (10 examples)
  - Complex object validation with dive
  
- [**custom-messages/**](_examples/custom-messages/) - Error translation (11 examples)
  - Localized messages and custom formats
  
- [**advanced/**](_examples/advanced/) - Custom validators (5 examples)
  - Thai ID Card, Business Hours, Custom validators
  
- [**real-world/**](_examples/real-world/) - Real scenarios (27 examples)
  - User registration, Payment processing, E-commerce orders
  
- [**best-practices/**](_examples/best-practices/) - Recommended patterns (7 examples)
  - String validation + decimal.Decimal workflow

### 🚀 Running Examples

```bash
# Run any example
cd _examples/basic
go run main.go

# Decimal has 2 files (run separately)
cd _examples/decimal
go run main.go
go run conditional.go

# Real-world has 3 files (run individually)
cd _examples/real-world
go run user-registration.go
go run payment.go
go run ecommerce.go
```

For complete documentation, see [_examples/README.md](_examples/README.md)

## Best Practices

### Decimal Handling for Monetary Values

**✅ Recommended Pattern:**

```go
import "github.com/shopspring/decimal"

// 1. Use string in API/struct layer
type PriceRequest struct {
    Amount string `json:"amount" validate:"required,decimal=10:2,dgt=0"`
}

// 2. Validate using xvalidator
if err := validator.StructTranslated(req); err != nil {
    return err
}

// 3. Convert to decimal.Decimal for calculations
amount, _ := decimal.NewFromString(req.Amount)
tax := amount.Mul(decimal.NewFromFloat(0.07))
total := amount.Add(tax)

// 4. Format back to string for output
return total.StringFixed(2)
```

**Why this pattern?**
- ✅ Prevents floating-point precision errors
- ✅ Validates before processing
- ✅ Maintains accuracy in calculations
- ✅ Consistent between validation and business logic

**❌ Avoid:**
```go
// Don't use float64 for money
type BadRequest struct {
    Amount float64 `json:"amount"` // ❌ Precision errors!
}

// Don't use decimal.Decimal directly in API structs
type BadRequest struct {
    Amount decimal.Decimal `json:"amount"` // ❌ JSON marshaling issues!
}
```

See [_examples/best-practices/decimal-workflow.go](_examples/best-practices/decimal-workflow.go) for complete working examples.

### Custom Validators

Create reusable custom validators:

```go
func validateBusinessHours(fl validator.FieldLevel) bool {
    timeStr := fl.Field().String()
    t, err := time.Parse("15:04", timeStr)
    if err != nil {
        return false
    }
    
    start, _ := time.Parse("15:04", "09:00")
    end, _ := time.Parse("15:04", "17:00")
    
    return t.After(start) && t.Before(end)
}

// Register
v.GetValidator().RegisterValidation("business_hours", validateBusinessHours)
```

See [_examples/advanced/main.go](_examples/advanced/main.go) for more custom validator examples.

## Testing

Run the test suite:

```bash
go test -v ./...
```

Run tests with coverage:

```bash
go test -v -cover ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Built on top of the excellent [go-playground/validator](https://github.com/go-playground/validator) library.

## Support

For questions and issues, please open an issue on [GitHub](https://github.com/hotfixfirst/go-xvalidator/issues).
