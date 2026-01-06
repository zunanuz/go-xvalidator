# Nested Struct Validation Examples

This example demonstrates validation of complex nested structures, arrays, and slices.

## What You'll Learn

- Nested struct validation
- Pointer field validation
- Array and slice validation
- `dive` validator for array elements
- Multi-level nesting
- Validation of structs within arrays

## Running the Example

```bash
cd nested-struct
go run main.go
```

## Nested Validation Tags

### Core Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `dive` | Validate each element in array/slice | `validate:"dive,email"` |
| `required` | Field must not be nil/empty | `validate:"required"` |
| `omitempty` | Skip if empty | `validate:"omitempty,dive"` |
| `min=n` | Minimum array length | `validate:"min=1,dive"` |
| `max=n` | Maximum array length | `validate:"max=10,dive"` |

## Code Examples

### 1. Simple Nested Struct

```go
type Address struct {
    Street  string `validate:"required,min=5"`
    City    string `validate:"required"`
    Country string `validate:"required,len=2"`
}

type User struct {
    Name    string   `validate:"required"`
    Address *Address `validate:"required"`
}

user := User{
    Name: "John Doe",
    Address: &Address{
        Street:  "123 Main St",
        City:    "Bangkok",
        Country: "TH",
    },
}
```

### 2. Array Validation with `dive`

```go
type Contact struct {
    Phones []string `validate:"required,min=1,dive,mobile_e164"`
    Emails []string `validate:"required,min=1,dive,email"`
}

contact := Contact{
    Phones: []string{
        "+66812345678",  // Each validated
        "+66823456789",  // Each validated
    },
    Emails: []string{
        "user@example.com",  // Each validated
    },
}
```

**How `dive` works:**
- `required,min=1` - Array itself must have at least 1 element
- `dive` - Apply validation to each element
- `mobile_e164` or `email` - Validation rule for each element

### 3. Array of Structs

```go
type OrderItem struct {
    ProductID string `validate:"required"`
    Quantity  int    `validate:"required,min=1"`
    Price     string `validate:"required,dgte=0"`
}

type Order struct {
    Items []OrderItem `validate:"required,min=1,dive"`
}

order := Order{
    Items: []OrderItem{
        {
            ProductID: "PROD-001",
            Quantity:  2,
            Price:     "29.99",
        },
        {
            ProductID: "PROD-002",
            Quantity:  1,
            Price:     "49.99",
        },
    },
}
```

### 4. Multi-Level Nesting

```go
type Company struct {
    Name      string   `validate:"required"`
    Address   Address  `validate:"required"`
    Employees []User   `validate:"required,min=1,dive"`
}

company := Company{
    Name: "Tech Corp",
    Address: Address{
        Street:  "456 Business Blvd",
        City:    "Bangkok",
        Country: "TH",
    },
    Employees: []User{
        {
            Name: "Employee 1",
            Address: &Address{
                Street:  "123 Home St",
                City:    "Bangkok",
                Country: "TH",
            },
        },
    },
}
```

### 5. Optional Nested Fields

```go
type UserProfile struct {
    Name    string   `validate:"required"`
    Address *Address `validate:"omitempty"`  // Optional but validate if present
    Phones  []string `validate:"omitempty,dive,mobile_e164"`  // Optional array
}

// Valid - no address
profile := UserProfile{
    Name:    "John",
    Address: nil,  // ✅ OK (omitempty)
}

// Valid - with address
profile := UserProfile{
    Name: "John",
    Address: &Address{
        Street:  "123 Main St",
        City:    "Bangkok",
        Country: "TH",
    },  // ✅ Validated
}
```

## Common Patterns

### ✅ Valid Examples

```go
// Nested struct
user := User{
    Name: "John",
    Address: &Address{
        Street:  "123 Main St",
        City:    "Bangkok",
        Country: "TH",
    },
}

// Array validation
contact := Contact{
    Phones: []string{"+66812345678"},  // ✅ All valid
    Emails: []string{"user@example.com"},  // ✅ All valid
}

// Empty optional array
product := Product{
    Name:       "Item",
    Categories: []string{},  // ✅ OK if omitempty
}

// Nested arrays
company := Company{
    Employees: []User{
        {Name: "Emp1", Address: &Address{...}},
        {Name: "Emp2", Address: &Address{...}},
    },  // ✅ Each user validated
}
```

### ❌ Invalid Examples

```go
// Missing nested field
user := User{
    Name: "John",
    Address: &Address{
        Street:  "123 Main St",
        City:    "",  // ❌ Required
        Country: "TH",
    },
}

// Invalid array element
contact := Contact{
    Phones: []string{
        "+66812345678",  // ✅ Valid
        "0812345678",    // ❌ Invalid format
    },
}

// Empty required array
contact := Contact{
    Phones: []string{},  // ❌ min=1 required
}

// Invalid nested struct in array
company := Company{
    Employees: []User{
        {
            Name:  "Employee",
            Email: "not-an-email",  // ❌ Invalid email
        },
    },
}
```

## Validation Tags Explained

### Array Constraints

```go
// Minimum length
Tags []string `validate:"min=1,dive,required"`
// Must have at least 1 element, each must not be empty

// Maximum length
Tags []string `validate:"max=10,dive,min=2"`
// Max 10 elements, each must be at least 2 chars

// Exact length
Codes []string `validate:"len=3,dive,len=2"`
// Exactly 3 elements, each exactly 2 chars

// Range
Items []string `validate:"min=1,max=100,dive,required"`
// Between 1 and 100 elements, each required
```

### Nested Struct Constraints

```go
// Required pointer
Address *Address `validate:"required"`
// Must not be nil, and Address fields are validated

// Optional pointer
Address *Address `validate:"omitempty"`
// Can be nil, but validated if present

// Required embedded
Address Address `validate:"required"`
// Fields are always validated
```

## Complex Nesting Example

```go
type Organization struct {
    Name       string       `validate:"required"`
    Departments []Department `validate:"required,min=1,dive"`
}

type Department struct {
    Name      string    `validate:"required"`
    Manager   User      `validate:"required"`
    Employees []User    `validate:"required,min=1,dive"`
    Budget    string    `validate:"required,dgte=0"`
}

type User struct {
    Name    string   `validate:"required"`
    Email   string   `validate:"required,email"`
    Phone   string   `validate:"required,mobile_e164"`
    Address *Address `validate:"required"`
}

type Address struct {
    Street  string `validate:"required,min=5"`
    City    string `validate:"required"`
    Country string `validate:"required,len=2"`
}
```

**Validation Flow:**
1. Organization validated
2. Each Department in array validated
3. Each Department's Manager (User) validated
4. Each Employee (User) in Department validated
5. Each User's Address validated

## Error Messages

```
address.city is a required field
phones[1] must be a valid E.164 mobile number
employees[0].email must be a valid email address
items[2].quantity must be 1 or greater
```

Error messages include the path to the invalid field.

## Performance Tips

1. **Use pointers** for large nested structs to avoid copying
2. **Use omitempty** for optional fields to skip validation
3. **Limit array sizes** with max constraints
4. **Cache validator** - don't create new instance per validation

## Common Pitfalls

### ❌ Forgetting `dive` for arrays

```go
// Wrong - validates array, not elements
Emails []string `validate:"required,email"`

// Correct - validates each element
Emails []string `validate:"required,dive,email"`
```

### ❌ Nil pointer without omitempty

```go
// Will fail if Address is nil
Address *Address `validate:"required"`

// Use omitempty if optional
Address *Address `validate:"omitempty"`
```

### ❌ Empty required arrays

```go
// Will fail
Items []OrderItem `validate:"required,min=1,dive"`
Items: []OrderItem{}  // Empty array

// Must have at least one
Items: []OrderItem{{...}}
```

## Next Steps

- [**custom-messages/**](../custom-messages/) - Customize error messages
- [**real-world/**](../real-world/) - See complex validation in action
- [**advanced/**](../advanced/) - Advanced validation techniques

## Related Documentation

- [Main README](../../README.md)
- [Basic Examples](../basic/)
- [go-playground/validator dive](https://github.com/go-playground/validator)
