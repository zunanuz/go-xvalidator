# Real-World Examples

This folder contains practical, real-world validation scenarios demonstrating how to use go-xvalidator in production applications.

## Examples

### 1. [ecommerce.go](ecommerce.go) - E-Commerce Order Validation

Complete e-commerce order validation including:
- Shopping cart items with quantities and prices
- Shipping and billing addresses
- Decimal calculations (subtotal, tax, shipping, total)
- Nested struct validation
- Array validation with `dive`

**Run:**
```bash
go run ecommerce.go
```

**Key Features:**
- Multi-item cart validation
- Address validation with phone numbers (E.164)
- Decimal precision for money calculations
- Optional billing address
- Discount and tax validation

### 2. [user-registration.go](user-registration.go) - User Registration Form

Complete user registration form validation including:
- Personal information (name, email, phone)
- Account security (username, password with strength validation)
- Password confirmation matching
- Profile information (bio, website, date of birth)
- Address details
- Terms and conditions acceptance
- Age verification

**Run:**
```bash
go run user-registration.go
```

**Key Features:**
- Password strength validation
- Password confirmation with `eqfield`
- Optional fields with `omitempty`
- Phone number validation (E.164)
- URL validation for website
- Boolean validation for terms acceptance
- Age restriction (18+)

### 3. [payment.go](payment.go) - Payment Processing

Payment gateway integration validation including:
- Transaction and order IDs
- Amount validation with decimal precision
- Multiple payment methods (credit card, debit card, bank transfer, e-wallet)
- Conditional card details validation
- Secure callback URLs (HTTPS only)
- Currency and fee calculations

**Run:**
```bash
go run payment.go
```

**Key Features:**
- Conditional validation (`required_if` for card details)
- HTTPS-only URLs for security
- Decimal validation for financial amounts
- Maximum amount limits
- Card number and CVV validation
- Multiple payment method support

## Common Patterns

### 1. Nested Struct Validation

All examples use nested structs:
```go
type Order struct {
    ShippingAddress Address `validate:"required"`
    Items          []Item   `validate:"required,min=1,dive"`
}
```

### 2. Decimal Money Calculations

Financial amounts use decimal validation:
```go
type Payment struct {
    Amount string `validate:"required,decimal=10:2,dgte=0"`
    Tax    string `validate:"required,decimal=10:2,dgte=0"`
    Total  string `validate:"required,decimal=10:2,dgt=0"`
}
```

### 3. Conditional Validation

Payment method determines required fields:
```go
type Payment struct {
    Method      string       `validate:"required,oneof=card ewallet"`
    CardDetails *CardDetails `validate:"required_if=Method card"`
}
```

### 4. Array Validation

Cart items and phone numbers use `dive`:
```go
type Order struct {
    Items  []CartItem `validate:"required,min=1,dive"`
    Phones []string   `validate:"required,dive,mobile_e164"`
}
```

### 5. Security URLs

Payment callbacks require HTTPS:
```go
type PaymentRequest struct {
    SuccessURL string `validate:"required,https_url"`
    WebhookURL string `validate:"required,https_url"`
}
```

## Best Practices Demonstrated

### 1. **Use Translated Errors**
```go
if err := v.StructTranslated(data); err != nil {
    // User-friendly error messages
}
```

### 2. **Decimal for Money**
```go
// Always use string with decimal validation
Amount string `validate:"required,decimal=10:2,dgte=0"`
```

### 3. **Required vs Optional**
```go
// Required
Email string `validate:"required,email"`

// Optional but validated if present
Website string `validate:"omitempty,url"`
```

### 4. **Nested Validation**
```go
// Validate entire nested struct
Address Address `validate:"required"`

// Validate optional nested struct
BillingAddress *Address `validate:"omitempty"`
```

### 5. **Security First**
```go
// Always HTTPS for sensitive operations
WebhookURL string `validate:"required,https_url"`

// Strong passwords
Password string `validate:"required,password_strength"`
```

## Integration Examples

### HTTP Handler

```go
func CreateOrder(w http.ResponseWriter, r *http.Request) {
    var order EcommerceOrder
    json.NewDecoder(r.Body).Decode(&order)
    
    v, _ := xvalidator.NewValidator()
    if err := v.StructTranslated(order); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Process valid order
}
```

### Service Layer

```go
func (s *UserService) Register(req UserRegistration) error {
    v, _ := xvalidator.NewValidator()
    if err := v.StructTranslated(req); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Hash password and create user
    return s.repository.Create(req)
}
```

### CLI Application

```go
func main() {
    var config AppConfig
    // Load config from file or flags
    
    v, _ := xvalidator.NewValidator()
    if err := v.StructTranslated(config); err != nil {
        fmt.Fprintf(os.Stderr, "Invalid configuration: %v\n", err)
        os.Exit(1)
    }
}
```

## Validation Checklist

When implementing validation, consider:

- ✅ Required vs optional fields
- ✅ String length constraints
- ✅ Numeric ranges
- ✅ Email and URL formats
- ✅ Phone number format (E.164)
- ✅ Password strength
- ✅ Decimal precision for money
- ✅ Array/slice validation with `dive`
- ✅ Nested struct validation
- ✅ Conditional validation
- ✅ Security (HTTPS, strong passwords)
- ✅ User-friendly error messages

## Error Handling

All examples use `StructTranslated()` for user-friendly errors:

```go
// Before: Raw error
Key: 'Order.Items[0].Quantity' Error:Field validation for 'Quantity' failed on the 'min' tag

// After: Translated error
items[0].quantity must be 1 or greater
```

## Testing

Each example includes both valid and invalid test cases:
- ✅ Valid complete data
- ❌ Missing required fields
- ❌ Invalid formats (email, phone, URL)
- ❌ Out of range values
- ❌ Failed conditional validation
- ❌ Invalid nested data

## Next Steps

- [**advanced/**](../advanced/) - Advanced validation techniques
- [**custom-messages/**](../custom-messages/) - Customizing error messages
- [Main README](../../README.md)

## Production Tips

1. **Cache validator instance** - Create once, reuse
2. **Use pointers for large structs** - Avoid copying
3. **Validate early** - Before database operations
4. **Log validation errors** - For debugging
5. **Return specific errors** - Help users fix issues
6. **Test edge cases** - Empty arrays, nil pointers, etc.

## Related Documentation

- [Main README](../../README.md)
- [Basic Examples](../basic/)
- [Nested Struct Examples](../nested-struct/)
