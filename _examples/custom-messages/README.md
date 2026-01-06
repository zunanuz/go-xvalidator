# Custom Error Messages Examples

This example demonstrates how to use translated, user-friendly error messages.

## What You'll Learn

- Using `StructTranslated()` for translated errors
- Using `VarTranslated()` for single field validation
- Comparing raw vs translated error messages
- Custom error messages for decimal validators
- Error messages for custom validators (phone, URL, password)
- Understanding error message formatting

## Running the Example

```bash
cd custom-messages
go run main.go
```

## Validation Methods

### With Translation

| Method | Description | Example |
|--------|-------------|---------|
| `StructTranslated()` | Validate struct with translated errors | `v.StructTranslated(user)` |
| `VarTranslated()` | Validate single field with translation | `v.VarTranslated(email, "email")` |

### Without Translation

| Method | Description | Example |
|--------|-------------|---------|
| `Struct()` | Validate struct with raw errors | `v.Struct(user)` |
| `Var()` | Validate single field | `v.Var(email, "email")` |
| `Validate()` | General validation (auto-translated) | `v.Validate(user)` |

## Code Examples

### 1. Basic Translated Validation

```go
v, _ := xvalidator.NewValidator()

type User struct {
    Name  string `validate:"required,min=3"`
    Email string `validate:"required,email"`
}

user := User{
    Name:  "AB",  // Too short
    Email: "invalid",
}

// With translation
if err := v.StructTranslated(user); err != nil {
    fmt.Println(err)
    // Output: name must be at least 3 characters in length; email must be a valid email address
}
```

### 2. Single Field Validation

```go
email := "not-an-email"

// With translation
if err := v.VarTranslated(email, "required,email"); err != nil {
    fmt.Println(err)
    // Output: must be a valid email address
}

// Without translation (raw)
if err := v.Var(email, "required,email"); err != nil {
    fmt.Println(err)
    // Output: Key: '' Error:Field validation for '' failed on the 'email' tag
}
```

### 3. Comparing Raw vs Translated

```go
type Product struct {
    Name  string `json:"name" validate:"required,min=3"`
    Price string `json:"price" validate:"required,dgte=0"`
}

invalid := Product{
    Name:  "AB",
    Price: "-10",
}

// Raw errors
v.Struct(invalid)
// Output: Key: 'Product.Name' Error:Field validation for 'Name' failed on the 'min' tag
//         Key: 'Product.Price' Error:Field validation for 'Price' failed on the 'dgte' tag

// Translated errors
v.StructTranslated(invalid)
// Output: name must be at least 3 characters in length; price must be a valid decimal greater than or equal to 0
```

## Error Message Formats

### Built-in Validators

```go
// Required field
`validate:"required"`
// Error: field_name is a required field

// Min length
`validate:"min=3"`
// Error: field_name must be at least 3 characters in length

// Max length
`validate:"max=50"`
// Error: field_name must be a maximum of 50 characters in length

// Email
`validate:"email"`
// Error: field_name must be a valid email address

// URL
`validate:"url"`
// Error: field_name must be a valid URL

// Numeric range
`validate:"min=0,max=100"`
// Error: field_name must be 0 or greater
// Error: field_name must be 100 or less
```

### Custom Validators

```go
// Decimal validators
`validate:"dgte=0"`
// Error: field_name must be a valid decimal greater than or equal to 0

`validate:"dlt=1000"`
// Error: field_name must be a valid decimal less than 1000

`validate:"decimal=10:2"`
// Error: field_name must be a valid decimal(10:2) format

// Phone validator
`validate:"mobile_e164"`
// Error: field_name must be a valid E.164 mobile number

`validate:"mobile_e164=TH"`
// Error: field_name must be a valid E.164 mobile number for TH

// HTTPS URL
`validate:"https_url"`
// Error: field_name must be a valid HTTPS URL

// Password strength
`validate:"password_strength"`
// Error: password must be at least 8 characters long
// Error: password must contain at least one: uppercase letter, digit
```

## Complete Error Message Examples

### ✅ Decimal Validation

```go
type Payment struct {
    Amount   string `validate:"required,decimal=10:2,dgte=0"`
    Discount string `validate:"required,dgte=0,dlte=100"`
}

// Invalid
payment := Payment{
    Amount:   "1000.505",  // Too many decimals
    Discount: "150",       // Exceeds 100
}

// Error messages:
// - amount must be a valid decimal(10:2) format
// - discount must be a valid decimal less than or equal to 100
```

### ✅ Phone Validation

```go
type Contact struct {
    Phone string `validate:"required,mobile_e164"`
}

contact := Contact{
    Phone: "0812345678",  // Missing +
}

// Error message:
// - phone must be a valid E.164 mobile number
```

### ✅ Password Validation

```go
type User struct {
    Password string `validate:"required,password_strength"`
}

user := User{
    Password: "weak",
}

// Error message:
// - password must be at least 8 characters long
// - password must contain at least one: uppercase letter, digit, special character
```

### ✅ Multiple Errors

```go
type Product struct {
    Name     string `json:"name" validate:"required,min=3"`
    Price    string `json:"price" validate:"required,dgte=0"`
    Category string `json:"category" validate:"required"`
}

product := Product{
    Name:     "",      // Empty
    Price:    "-10",   // Negative
    Category: "",      // Empty
}

// Error message (semicolon-separated):
// name is a required field; price must be a valid decimal greater than or equal to 0; category is a required field
```

## Translation Features

### Field Names

- Uses JSON tag names when available
- Falls back to struct field names
- Converts to human-readable format

```go
type User struct {
    EmailAddress string `json:"email" validate:"required,email"`
}

// Error uses JSON tag name:
// "email must be a valid email address"
// NOT "EmailAddress must be..."
```

### Error Formatting

```go
// Single error
"email must be a valid email address"

// Multiple errors (semicolon-separated)
"name is a required field; email must be a valid email address; age must be 18 or greater"
```

## Use Cases

### 1. API Response Errors

```go
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)
    
    v, _ := xvalidator.NewValidator()
    if err := v.StructTranslated(user); err != nil {
        // Return user-friendly error
        json.NewEncoder(w).Encode(map[string]string{
            "error": err.Error(),
        })
        return
    }
}
```

### 2. Form Validation

```go
func ValidateForm(form RegistrationForm) []string {
    v, _ := xvalidator.NewValidator()
    
    if err := v.StructTranslated(form); err != nil {
        // Split semicolon-separated errors
        return strings.Split(err.Error(), "; ")
    }
    return nil
}
```

### 3. CLI Validation

```go
func ValidateConfig(config Config) {
    v, _ := xvalidator.NewValidator()
    
    if err := v.StructTranslated(config); err != nil {
        fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
        os.Exit(1)
    }
}
```

## Tips

1. **Always use `StructTranslated()`** for user-facing errors
2. **Use `Struct()`** for debugging or logging
3. **Split errors** by "; " for individual messages
4. **Use JSON tags** to control field names in errors
5. **Test error messages** to ensure they're user-friendly

## Method Comparison

| Method | Translation | Use Case |
|--------|-------------|----------|
| `Validate()` | ✅ Auto | General validation |
| `StructTranslated()` | ✅ Yes | User-facing errors |
| `VarTranslated()` | ✅ Yes | Single field user errors |
| `Struct()` | ❌ No | Debugging, logging |
| `Var()` | ❌ No | Debugging, testing |

## Next Steps

- [**real-world/**](../real-world/) - See translation in real applications
- [**advanced/**](../advanced/) - Advanced error handling
- [Main README](../../README.md)

## Related Documentation

- [Validator Documentation](../../README.md)
- [Translation System](../../translator.go)
