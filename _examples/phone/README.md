# Phone Number Validation Examples

This example demonstrates phone number validation using E.164 international format.

## What You'll Learn

- E.164 phone number format validation
- Country-specific phone validation
- Multiple phone number fields
- Optional phone fields
- International phone numbers

## Running the Example

```bash
cd phone
go run main.go
```

## Phone Validators

### E.164 Format

The E.164 format is the international standard for phone numbers:
- Starts with `+` followed by country code
- Contains only digits after the `+`
- Length: 7-15 digits (including country code)
- No spaces, dashes, or parentheses

**Format:** `+[country code][subscriber number]`

### Validation Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `mobile_e164` | Any valid E.164 format | `validate:"mobile_e164"` |
| `mobile_e164=TH` | Thailand mobile only | `validate:"mobile_e164=TH"` |
| `mobile_e164=US` | US mobile only | `validate:"mobile_e164=US"` |
| `mobile_e164=GB` | UK mobile only | `validate:"mobile_e164=GB"` |
| `mobile_e164=FR` | France mobile only | `validate:"mobile_e164=FR"` |

### Supported Countries

| Code | Country | Format Example | Country Code |
|------|---------|----------------|--------------|
| `TH` | Thailand | `+66812345678` | +66 |
| `US` | United States | `+12025551234` | +1 |
| `GB` | United Kingdom | `+447912345678` | +44 |
| `FR` | France | `+33612345678` | +33 |
| `AU` | Australia | `+61412345678` | +61 |
| `DE` | Germany | `+491701234567` | +49 |
| `JP` | Japan | `+819012345678` | +81 |
| `CN` | China | `+8613912345678` | +86 |

## Code Examples

### 1. Basic Phone Validation

```go
type Contact struct {
    Name  string `validate:"required"`
    Phone string `validate:"required,mobile_e164"`
}

contact := Contact{
    Name:  "John Doe",
    Phone: "+66812345678",  // ✅ Valid
}
```

### 2. Country-Specific Validation

```go
type UserContact struct {
    PhoneTH string `validate:"mobile_e164=TH"`
    PhoneUS string `validate:"mobile_e164=US"`
    PhoneGB string `validate:"mobile_e164=GB"`
}

contact := UserContact{
    PhoneTH: "+66812345678",  // ✅ Thai number
    PhoneUS: "+12025551234",  // ✅ US number
    PhoneGB: "+447912345678", // ✅ UK number
}
```

### 3. Optional Phone Fields

```go
type UserProfile struct {
    PrimaryPhone string `validate:"required,mobile_e164"`
    WorkPhone    string `validate:"omitempty,mobile_e164"`
    HomePhone    string `validate:"omitempty,mobile_e164"`
}

// Valid - only primary phone required
profile := UserProfile{
    PrimaryPhone: "+66812345678",
    WorkPhone:    "",  // ✅ OK (omitempty)
    HomePhone:    "",  // ✅ OK (omitempty)
}
```

### 4. Single Phone Validation

```go
v, _ := xvalidator.NewValidator()

phone := "+66812345678"
if err := v.Var(phone, "mobile_e164"); err != nil {
    // Invalid phone number
}

// Country-specific
usPhone := "+12025551234"
if err := v.Var(usPhone, "mobile_e164=US"); err != nil {
    // Invalid US phone number
}
```

## E.164 Format Examples

### ✅ Valid Phone Numbers

```go
// General E.164 format
"+66812345678"    // Thailand
"+12025551234"    // USA
"+447912345678"   // UK
"+33612345678"    // France
"+61412345678"    // Australia
"+491701234567"   // Germany
"+819012345678"   // Japan

// Country-specific
Contact{PhoneTH: "+66812345678"}  // Thai mobile
Contact{PhoneUS: "+12025551234"}  // US mobile
Contact{PhoneGB: "+447912345678"} // UK mobile
```

### ❌ Invalid Phone Numbers

```go
// Missing + prefix
"66812345678"      // ❌ Must start with +

// Contains non-digits
"+66-81-234-5678"  // ❌ No dashes
"+66 81 234 5678"  // ❌ No spaces
"+66(81)2345678"   // ❌ No parentheses
"+66ABC345678"     // ❌ No letters

// Too short or too long
"+661234"          // ❌ Too short (< 7 digits)
"+661234567890123456" // ❌ Too long (> 15 digits)

// Wrong country code for specific validation
Contact{
    PhoneTH: "+12025551234",  // ❌ US number, expecting TH
}
```

## Country Code Reference

### Thailand (+66)
- Mobile: `+668XXXXXXXX` or `+669XXXXXXXX`
- Example: `+66812345678`

### United States (+1)
- Mobile: `+1XXXXXXXXXX` (10 digits after +1)
- Example: `+12025551234`

### United Kingdom (+44)
- Mobile: `+447XXXXXXXXX`
- Example: `+447912345678`

### France (+33)
- Mobile: `+336XXXXXXXX` or `+337XXXXXXXX`
- Example: `+33612345678`

## Error Messages

```
phone_general must be a valid E.164 mobile number
phone_th must be a valid E.164 mobile number for TH
phone_us must be a valid E.164 mobile number for US
```

## Common Mistakes

### ❌ Forgetting the + prefix
```go
phone := "66812345678"  // Wrong
phone := "+66812345678" // Correct
```

### ❌ Including spaces or formatting
```go
phone := "+66 81 234 5678"  // Wrong
phone := "+66812345678"     // Correct
```

### ❌ Wrong country code
```go
// Expecting TH but providing US
Contact{
    PhoneTH: "+12025551234",  // Wrong country
}
```

### ❌ Using landline instead of mobile
```go
// Some validators check for mobile-specific patterns
landline := "+6625551234"  // May be invalid
mobile := "+66812345678"   // Correct
```

## Tips

1. **Always use + prefix** for E.164 format
2. **Store as string** to preserve leading zeros
3. **No formatting** - store raw E.164 format
4. **Country-specific validation** when you know the expected country
5. **Use omitempty** for optional phone fields
6. **Validate before storage** to ensure data quality

## Next Steps

- [**url/**](../url/) - URL validation
- [**password/**](../password/) - Password validation
- [**real-world/**](../real-world/) - See phone validation in context

## Related Documentation

- [Main README](../../README.md)
- [E.164 Format Specification](https://en.wikipedia.org/wiki/E.164)
