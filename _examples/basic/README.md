# Basic Validation Examples

This example demonstrates the fundamental validation features of go-xvalidator.

## What You'll Learn

- Required field validation
- String length constraints (min, max, len)
- Number range validation
- Email validation
- URL validation
- Alphanumeric validation
- Optional field handling
- Single field validation with `Var()`

## Running the Example

```bash
cd basic
go run main.go
```

## Validation Tags Used

### String Validators

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field must not be empty | `validate:"required"` |
| `min=n` | Minimum length for strings | `validate:"min=3"` |
| `max=n` | Maximum length for strings | `validate:"max=50"` |
| `len=n` | Exact length | `validate:"len=8"` |
| `email` | Valid email format | `validate:"email"` |
| `url` | Valid URL format | `validate:"url"` |
| `alphanum` | Alphanumeric characters only | `validate:"alphanum"` |
| `omitempty` | Skip validation if empty | `validate:"omitempty,url"` |

### Number Validators

| Tag | Description | Example |
|-----|-------------|---------|
| `min=n` | Minimum value for numbers | `validate:"min=0"` |
| `max=n` | Maximum value for numbers | `validate:"max=120"` |

## Code Examples

### 1. Basic Struct with Validation

```go
type User struct {
    Name     string `json:"name" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"required,min=18,max=120"`
    Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
    Website  string `json:"website" validate:"omitempty,url"`
}
```

### 2. Creating Validator

```go
v, err := xvalidator.NewValidator()
if err != nil {
    panic(err)
}
```

### 3. Validating Struct

```go
user := User{
    Name:     "John Doe",
    Email:    "john.doe@example.com",
    Age:      25,
    Username: "johndoe123",
}

if err := v.Validate(user); err != nil {
    fmt.Printf("Validation failed: %v\n", err)
}
```

### 4. Validating Single Field

```go
email := "test@example.com"
if err := v.Var(email, "required,email"); err != nil {
    fmt.Printf("Email validation failed: %v\n", err)
}
```

## Common Validation Scenarios

### ✅ Valid Examples

```go
// Valid user
user := User{
    Name:     "John Doe",
    Email:    "john@example.com",
    Age:      25,
    Username: "johndoe",
    Website:  "https://example.com",
}

// Valid with empty optional field
user := User{
    Name:     "Jane Doe",
    Email:    "jane@example.com",
    Age:      30,
    Username: "janedoe",
    Website:  "", // OK because of omitempty
}
```

### ❌ Invalid Examples

```go
// Missing required field
user := User{
    Name:     "",  // ❌ required
    Email:    "john@example.com",
    Age:      25,
    Username: "johndoe",
}

// Invalid email format
user := User{
    Name:     "John Doe",
    Email:    "invalid-email",  // ❌ must be valid email
    Age:      25,
    Username: "johndoe",
}

// Age out of range
user := User{
    Name:     "Young User",
    Email:    "young@example.com",
    Age:      16,  // ❌ must be >= 18
    Username: "younguser",
}

// Username too short
user := User{
    Name:     "John Doe",
    Email:    "john@example.com",
    Age:      25,
    Username: "ab",  // ❌ must be >= 3 characters
}
```

## Error Messages

When validation fails, you'll receive detailed error messages:

```
name is a required field; age must be 18 or greater
```

## Next Steps

- [**decimal/**](../decimal/) - Learn about decimal validation
- [**phone/**](../phone/) - Phone number validation
- [**password/**](../password/) - Password strength validation
- [**nested-struct/**](../nested-struct/) - Complex struct validation

## Related Documentation

- [Main README](../../README.md)
- [go-playground/validator docs](https://github.com/go-playground/validator)
