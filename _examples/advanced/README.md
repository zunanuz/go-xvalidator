# Advanced Examples

This folder contains advanced validation techniques including custom validators, struct-level validation, and complex validation patterns.

## Overview

Learn how to extend go-xvalidator with:

- Custom validation functions
- Struct-level cross-field validation
- Parameterized validators
- Chained validators
- Conditional logic in validators
- Integration patterns

## Examples

### 1. Custom Validator - Thai ID Card

Validates Thai national ID card numbers using the MOD 11 checksum algorithm.

**Features:**

- 13-digit format validation
- MOD 11 checksum verification
- Custom error messages

```go
v.GetValidator().RegisterValidation("thai_id_card", validateThaiIDCard)

type ThaiCitizen struct {
    IDCard string `validate:"required,thai_id_card"`
}
```

### 2. Custom Validator - Business Hours

Validates time is within business hours (9:00 AM - 5:00 PM).

**Features:**

- Time parsing and validation
- Range checking
- Custom time format support

```go
v.GetValidator().RegisterValidation("business_hours", validateBusinessHours)

type Appointment struct {
    Time string `validate:"required,business_hours"`
}
```

### 3. Custom Validator - Thai Bank Account

Validates Thai bank account format: XXX-X-XXXXX-X.

**Features:**

- Regex pattern matching
- Format validation
- Hyphen-separated format

```go
v.GetValidator().RegisterValidation("thai_bank_account", validateThaiBankAccount)

type BankTransfer struct {
    FromAccount string `validate:"required,thai_bank_account"`
}
```

### 4. Custom Validator - Strong Password

Validates password meets strong criteria (12+ chars, upper, lower, digit, special).

**Features:**

- Minimum length check
- Character type requirements
- Custom strength rules

```go
v.GetValidator().RegisterValidation("strong_password", validateStrongPassword)

type UserAccount struct {
    Password string `validate:"required,strong_password"`
}
```

### 5. Custom Validator - Future Date

Validates date is in the future (for events, bookings, etc.).

**Features:**

- Date parsing
- Future date check
- Time comparison

```go
v.GetValidator().RegisterValidation("future_date", validateFutureDate)

type Event struct {
    StartDate string `validate:"required,future_date"`
}
```

### 6. Custom Validator - Decimal Range

Validates decimal values within specified min/max range.

**Features:**

- Parameterized validator factory
- Decimal precision handling
- Multiple range validators

```go
v.GetValidator().RegisterValidation("product_price", 
    validateDecimalRange(1.00, 1000000.00))

type Product struct {
    Price string `validate:"required,decimal=10:2,product_price"`
}
```

### 7. Custom Validator - Thai Phone Number

Validates Thai phone number format (10 digits, 0XXXXXXXXX).

**Features:**

- Prefix validation (02, 06, 08, 09)
- Length check
- Digit validation

```go
v.GetValidator().RegisterValidation("thai_phone", validateThaiPhone)

type Contact struct {
    ThaiPhone string `validate:"required,thai_phone"`
}
```

### 8. Struct-Level Cross-Field Validation

Validates relationships between multiple fields (e.g., check-in before check-out).

**Features:**

- Multiple field comparison
- Date range validation
- Custom struct-level errors

```go
v.GetValidator().RegisterStructValidation(validateDateRange, Booking{})

type Booking struct {
    CheckIn  string `validate:"required"`
    CheckOut string `validate:"required"`
}
```

### 9. Chaining Multiple Custom Validators

Combines multiple custom validators on a single field.

**Features:**

- Multiple validator composition
- Sequential validation
- Combined validation rules

```go
type APIKey struct {
    Name string `validate:"required,no_spaces,alphanum_underscore"`
}
```

### 10. Custom Validator with Parameters

Creates validators that accept parameters from validation tags.

**Features:**

- Parameter extraction from tags
- Dynamic validation rules
- Reusable validators

```go
v.GetValidator().RegisterValidation("min_age", validateMinAge)

type DriverLicense struct {
    Birthdate string `validate:"required,min_age=18"`
}
```

## Running the Examples

```bash
cd _examples/advanced
go run main.go
```

## Custom Validator Patterns

### 1. Simple Field Validator

```go
func validateCustomField(fl validator.FieldLevel) bool {
    value := fl.Field().String()
    // Your validation logic
    return isValid
}
```

### 2. Validator Factory (with configuration)

```go
func validateRange(min, max float64) validator.Func {
    return func(fl validator.FieldLevel) bool {
        value := parseValue(fl.Field())
        return value >= min && value <= max
    }
}
```

### 3. Struct-Level Validator

```go
func validateStruct(sl validator.StructLevel) {
    data := sl.Current().Interface().(YourType)
    
    if !isValid(data) {
        sl.ReportError(data.Field, "Field", "Field", "custom_tag", "")
    }
}
```

### 4. Parameterized Validator

```go
func validateWithParam(fl validator.FieldLevel) bool {
    value := fl.Field().String()
    param := fl.Param() // Get parameter from tag
    
    // Use param in validation logic
    return validate(value, param)
}
```

## Registration Methods

### Field-Level Validator

```go
v, _ := xvalidator.NewValidator()
v.GetValidator().RegisterValidation("tag_name", validatorFunc)
```

### Struct-Level Validator

```go
v, _ := xvalidator.NewValidator()
v.GetValidator().RegisterStructValidation(structValidatorFunc, StructType{})
```

### Multiple Validators

```go
validators := map[string]validator.Func{
    "thai_id":    validateThaiID,
    "thai_phone": validateThaiPhone,
    "thai_bank":  validateThaiBank,
}

for tag, fn := range validators {
    v.GetValidator().RegisterValidation(tag, fn)
}
```

## Best Practices

### 1. **Error Handling**

```go
func validateSafe(fl validator.FieldLevel) bool {
    defer func() {
        if r := recover(); r != nil {
            // Log error
            return false
        }
    }()
    
    // Validation logic
    return true
}
```

### 2. **Performance**

```go
// Compile regex once
var phoneRegex = regexp.MustCompile(`^0[0-9]{9}$`)

func validatePhone(fl validator.FieldLevel) bool {
    return phoneRegex.MatchString(fl.Field().String())
}
```

### 3. **Reusability**

```go
// Create validator factories
func createRangeValidator(min, max int) validator.Func {
    return func(fl validator.FieldLevel) bool {
        value := fl.Field().Int()
        return value >= int64(min) && value <= int64(max)
    }
}
```

### 4. **Type Safety**

```go
func validateTypeSafe(fl validator.FieldLevel) bool {
    // Check field kind first
    if fl.Field().Kind() != reflect.String {
        return false
    }
    
    value := fl.Field().String()
    // Validate value
    return isValid(value)
}
```

### 5. **Documentation**

```go
// validateThaiIDCard validates Thai national ID card number.
// Format: 13 digits with MOD 11 checksum validation.
// Example: 1234567890123
// Returns: true if valid, false otherwise
func validateThaiIDCard(fl validator.FieldLevel) bool {
    // Implementation
}
```

## Common Custom Validators

### Geographic

- Thai ID card
- Thai phone number
- Thai bank account
- Thai postal code
- Province validation

### Business Rules

- Business hours
- Future dates
- Date ranges
- Age requirements
- Price ranges

### Security

- Strong passwords
- API key formats
- Token validation
- Secure URLs (HTTPS only)
- Encryption requirements

### Financial

- Decimal ranges
- Currency validation
- Account numbers
- Transaction limits
- Tax calculations

### Temporal

- Future dates
- Past dates
- Date ranges
- Age calculation
- Expiration dates

## Integration Examples

### HTTP Middleware

```go
func ValidateMiddleware(v *xvalidator.Validator) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            var data RequestData
            json.NewDecoder(r.Body).Decode(&data)
            
            if err := v.StructTranslated(data); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### gRPC Interceptor

```go
func ValidationInterceptor(v *xvalidator.Validator) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, 
                info *grpc.UnaryServerInfo, 
                handler grpc.UnaryHandler) (interface{}, error) {
        
        if err := v.StructTranslated(req); err != nil {
            return nil, status.Error(codes.InvalidArgument, err.Error())
        }
        
        return handler(ctx, req)
    }
}
```

### Gin Handler

```go
func RegisterUser(c *gin.Context) {
    var user UserRegistration
    
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    v, _ := xvalidator.NewValidator()
    if err := v.StructTranslated(user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Process registration
}
```

### Echo Validator

```go
type CustomValidator struct {
    validator *xvalidator.Validator
}

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.StructTranslated(i)
}

func main() {
    e := echo.New()
    v, _ := xvalidator.NewValidator()
    e.Validator = &CustomValidator{validator: v}
}
```

## Testing Custom Validators

```go
func TestCustomValidator(t *testing.T) {
    v, _ := xvalidator.NewValidator()
    v.GetValidator().RegisterValidation("custom", validateCustom)
    
    tests := []struct {
        name    string
        input   TestStruct
        wantErr bool
    }{
        {"valid", TestStruct{Field: "valid"}, false},
        {"invalid", TestStruct{Field: "invalid"}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := v.StructTranslated(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Error Messages

### Default Error

```go
// Built-in error message
Field validation for 'IDCard' failed on the 'thai_id_card' tag
```

### Custom Error Message

```go
// Register custom translation
v.GetTranslator().AddTranslation("thai_id_card", 
    "รหัสบัตรประชาชนไม่ถูกต้อง")
```

### Parameterized Error

```go
// Error with parameter
v.GetTranslator().AddTranslation("min_age", 
    "อายุต้องไม่น้อยกว่า {0} ปี")
```

## Next Steps

- [Main README](../../README.md)
- [Basic Examples](../basic/)
- [Real-World Examples](../real-world/)
- [Custom Messages](../custom-messages/)

## Resources

- [go-playground/validator Documentation](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [Custom Validators Guide](https://github.com/go-playground/validator#custom-validation-functions)
- [Struct Level Validation](https://github.com/go-playground/validator#struct-level-validation)
