# Decimal Validation Examples

This example demonstrates decimal validation features for precise numeric comparisons.

## What You'll Learn

- Decimal comparison operators (greater than, less than, equal)
- Precision and scale validation
- Conditional decimal validation based on other fields
- Working with monetary values
- Integer-only decimal validation

## Running the Examples

```bash
# Run main decimal examples
cd decimal
go run main.go

# Run conditional decimal examples
go run conditional.go
```

## Decimal Validators

### Comparison Operators

| Tag | Description | Example |
|-----|-------------|---------|
| `dgt=value` | Decimal greater than | `validate:"dgt=0"` |
| `dgte=value` | Decimal greater than or equal | `validate:"dgte=0"` |
| `dlt=value` | Decimal less than | `validate:"dlt=1000000"` |
| `dlte=value` | Decimal less than or equal | `validate:"dlte=100"` |
| `deq=value` | Decimal equal to | `validate:"deq=100"` |
| `dneq=value` | Decimal not equal to | `validate:"dneq=0"` |

### Precision and Scale

| Tag | Description | Example |
|-----|-------------|---------|
| `decimal` | Default precision (38:18) | `validate:"decimal"` |
| `decimal=scale` | Specify scale only | `validate:"decimal=2"` |
| `decimal=precision:scale` | Specify both | `validate:"decimal=10:2"` |
| `decimal=0` | Integer only (no decimals) | `validate:"decimal=0"` |

### Conditional Decimal

| Tag | Description | Example |
|-----|-------------|---------|
| `decimal_if=rule@Field=value` | Apply rule if condition met | `validate:"decimal_if=2@Type=credit"` |

## Code Examples

### 1. Basic Decimal Comparison

```go
type PriceData struct {
    RegularPrice string `validate:"required,dgte=0"`
    SalePrice    string `validate:"required,dgte=0,dlt=1000000"`
    Discount     string `validate:"required,dgte=0,dlte=100"`
    MinOrder     string `validate:"required,dgt=0"`
}
```

### 2. Decimal with Precision/Scale

```go
type Invoice struct {
    Subtotal string `validate:"required,decimal=10:2,dgte=0"`
    Tax      string `validate:"required,decimal=10:2,dgte=0"`
    Total    string `validate:"required,decimal=10:2,dgte=0"`
}
```

**Precision:Scale Explanation:**
- `10:2` means 10 total digits, 2 after decimal point
- Max value: 99999999.99 (8 digits before, 2 after)
- Valid: `1234.56`, `99999999.99`
- Invalid: `123.456` (too many decimal places), `123456789.99` (too many digits)

### 3. Conditional Decimal Validation

```go
type PaymentMethod struct {
    Type   string `validate:"required,oneof=credit debit cash"`
    Amount string `validate:"required,decimal_if=2@Type=credit,decimal_if=2@Type=debit,decimal_if=0@Type=cash"`
}
```

**Explanation:**
- If `Type=credit`: Amount must have exactly 2 decimal places
- If `Type=debit`: Amount must have exactly 2 decimal places
- If `Type=cash`: Amount must be integer (0 decimal places)

### 4. Equality Validators

```go
type BankAccount struct {
    MinBalance string `validate:"required,deq=100"`    // Must equal exactly 100
    FeeAmount  string `validate:"required,dneq=0"`     // Cannot be 0
}
```

## Common Use Cases

### ✅ Valid Examples

```go
// Valid price comparisons
price := PriceData{
    RegularPrice: "99.99",
    SalePrice:    "79.99",
    Discount:     "20.00",
    MinOrder:     "1.00",  // Greater than 0
}

// Valid precision/scale
invoice := Invoice{
    Subtotal: "1000.50",  // 10:2 format
    Tax:      "70.04",
    Total:    "1070.54",
}

// Valid conditional - credit card with decimals
payment := PaymentMethod{
    Type:   "credit",
    Amount: "100.50",  // 2 decimal places OK for credit
}

// Valid conditional - cash without decimals
payment := PaymentMethod{
    Type:   "cash",
    Amount: "100",  // Integer OK for cash
}
```

### ❌ Invalid Examples

```go
// Negative price
price := PriceData{
    RegularPrice: "-10.00",  // ❌ must be >= 0
}

// Discount over 100%
price := PriceData{
    Discount: "120.00",  // ❌ must be <= 100
}

// MinOrder not greater than 0
price := PriceData{
    MinOrder: "0.00",  // ❌ must be > 0
}

// Too many decimal places
invoice := Invoice{
    Subtotal: "1000.505",  // ❌ max 2 decimal places
}

// Cash payment with decimals
payment := PaymentMethod{
    Type:   "cash",
    Amount: "100.50",  // ❌ cash must be integer
}

// Fee amount is zero
account := BankAccount{
    FeeAmount: "0.00",  // ❌ must not equal 0
}
```

## Precision/Scale Guide

| Format | Total Digits | Decimal Places | Example Valid | Example Invalid |
|--------|--------------|----------------|---------------|-----------------|
| `10:2` | 10 | 2 | `12345678.99` | `123456789.99` |
| `10:0` | 10 | 0 (integer) | `1234567890` | `123.45` |
| `5:2` | 5 | 2 | `123.45` | `1234.56` |
| `8:3` | 8 | 3 | `12345.678` | `12345.6789` |
| Default (38:18) | 38 | 18 | Very large | - |

## Error Messages

```
regular_price must be a valid decimal greater than or equal to 0
discount must be a valid decimal less than or equal to 100
min_order must be a valid decimal greater than 0
subtotal must be a valid decimal(10:2) format
amount must be a valid decimal(0) format when Type equals cash
```

## Tips

1. **Use string type** for decimal values to avoid floating-point precision issues
2. **Specify precision/scale** for financial calculations
3. **Use conditional validation** for different payment methods or product types
4. **Combine with other validators** for complete validation

## Next Steps

- [**phone/**](../phone/) - Phone number validation
- [**url/**](../url/) - URL validation
- [**real-world/**](../real-world/) - See decimal validation in action

## Related Documentation

- [Main README](../../README.md)
- [Basic Examples](../basic/)
