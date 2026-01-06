# Password Strength Validation Examples

This example demonstrates password strength validation with comprehensive security requirements.

## What You'll Learn

- Password strength validation
- Minimum length requirements
- Character type requirements (uppercase, lowercase, digits, special)
- Maximum length constraints
- Optional password fields
- Password confirmation matching
- Standalone password validation

## Running the Example

```bash
cd password
go run main.go
```

## Password Requirements

The `password_strength` validator enforces these requirements:

| Requirement | Rule |
|-------------|------|
| **Minimum Length** | At least 8 characters |
| **Maximum Length** | No more than 100 characters |
| **Uppercase Letter** | At least one (A-Z) |
| **Lowercase Letter** | At least one (a-z) |
| **Digit** | At least one (0-9) |
| **Special Character** | At least one from: `!@#$%^&*()_+-=[]{}|;:,.<>?` |

## Password Validators

### Validation Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `password_strength` | Comprehensive password validation | `validate:"password_strength"` |

### Combining with Other Tags

```go
// Required password
Password string `validate:"required,password_strength"`

// Optional but strong if provided
Password string `validate:"omitempty,password_strength"`

// With confirmation
NewPassword     string `validate:"required,password_strength"`
ConfirmPassword string `validate:"required,eqfield=NewPassword"`
```

## Code Examples

### 1. User Registration

```go
type UserRegistration struct {
    Username string `validate:"required,min=3,max=20"`
    Email    string `validate:"required,email"`
    Password string `validate:"required,password_strength"`
}

user := UserRegistration{
    Username: "johndoe",
    Email:    "john@example.com",
    Password: "MyP@ssw0rd123",  // ✅ Strong password
}
```

### 2. Password Change

```go
type PasswordChange struct {
    CurrentPassword string `validate:"required"`
    NewPassword     string `validate:"required,password_strength"`
    ConfirmPassword string `validate:"required,eqfield=NewPassword"`
}

change := PasswordChange{
    CurrentPassword: "OldP@ss123",
    NewPassword:     "NewP@ss456",
    ConfirmPassword: "NewP@ss456",  // Must match NewPassword
}
```

### 3. Optional Strong Password

```go
type Account struct {
    Username string `validate:"required"`
    Password string `validate:"omitempty,password_strength"`
}

// Valid - empty password
account := Account{
    Username: "admin",
    Password: "",  // ✅ OK because of omitempty
}

// Valid - strong password
account := Account{
    Username: "admin",
    Password: "Str0ng@Pass",  // ✅ Meets requirements
}

// Invalid - weak password
account := Account{
    Username: "admin",
    Password: "weak",  // ❌ Doesn't meet requirements
}
```

### 4. Standalone Validation

```go
// Using struct validation
v, _ := xvalidator.NewValidator()
password := "SecureP@ss123"
if err := v.Var(password, "password_strength"); err != nil {
    // Password is weak
}

// Using dedicated function
password := "MyV@lid123"
if err := xvalidator.ValidatePasswordStrength(password); err != nil {
    // Password is weak
}
```

## Valid Password Examples

### ✅ Strong Passwords

```go
"MyP@ssw0rd123"      // All requirements met
"Secure!Pass1"       // Minimum length, all types
"C0mpl3x@Pwd"        // Good mix of characters
"Str0ng#P@ss"        // Multiple special chars
"Valid_Pass123"      // Underscore is special char
"Test@2024!"         // With year
"My$ecure9Pass"      // Dollar sign
"P@ssw0rd_2024"      // Combination
"Adm1n@Str0ng!"      // Admin password
"Us3r#S3cur3"        // Leetspeak style
```

### ❌ Weak Passwords

```go
// Too short (< 8 characters)
"Short1!"            // ❌ Only 7 characters
"Pass1!"             // ❌ Only 6 characters

// Missing uppercase
"mypassword123!"     // ❌ No uppercase
"test@password1"     // ❌ No uppercase

// Missing lowercase
"MYPASSWORD123!"     // ❌ No lowercase
"TEST@PASSWORD1"     // ❌ No lowercase

// Missing digit
"MyPassword!"        // ❌ No digit
"Test@Password"      // ❌ No digit

// Missing special character
"MyPassword123"      // ❌ No special char
"TestPassword1"      // ❌ No special char

// Too long (> 100 characters)
"MyP@ssw0rd" + string(make([]byte, 91))  // ❌ Over 100 chars

// Common weak passwords
"Password123!"       // ❌ Too common
"Qwerty123!"         // ❌ Too common
"12345678A!"         // ❌ Too common
```

## Accepted Special Characters

The following special characters are accepted:

```
!  @  #  $  %  ^  &  *  (  )  _  +  -  =  [  ]  {  }  |  ;  :  ,  .  <  >  ?
```

### Examples by Special Character

```go
"Test!123Pass"    // ✅ Exclamation mark
"Test@123Pass"    // ✅ At sign
"Test#123Pass"    // ✅ Hash
"Test$123Pass"    // ✅ Dollar
"Test%123Pass"    // ✅ Percent
"Test^123Pass"    // ✅ Caret
"Test&123Pass"    // ✅ Ampersand
"Test*123Pass"    // ✅ Asterisk
"Test_123Pass"    // ✅ Underscore
"Test+123Pass"    // ✅ Plus
"Test-123Pass"    // ✅ Dash
"Test=123Pass"    // ✅ Equals
"Test[123]Pass"   // ✅ Brackets
"Test{123}Pass"   // ✅ Braces
"Test|123Pass"    // ✅ Pipe
"Test:123Pass"    // ✅ Colon
"Test;123Pass"    // ✅ Semicolon
"Test,123Pass"    // ✅ Comma
"Test.123Pass"    // ✅ Period
"Test<123>Pass"   // ✅ Angle brackets
"Test?123Pass"    // ✅ Question mark
```

## Common Use Cases

### 1. User Registration Form

```go
type RegistrationForm struct {
    Email           string `validate:"required,email"`
    Password        string `validate:"required,password_strength"`
    ConfirmPassword string `validate:"required,eqfield=Password"`
    AcceptTerms     bool   `validate:"required,eq=true"`
}
```

### 2. Password Reset

```go
type PasswordReset struct {
    Token           string `validate:"required"`
    NewPassword     string `validate:"required,password_strength"`
    ConfirmPassword string `validate:"required,eqfield=NewPassword"`
}
```

### 3. Account Settings Update

```go
type AccountSettings struct {
    Email          string `validate:"required,email"`
    CurrentPassword string `validate:"required_with=NewPassword"`
    NewPassword     string `validate:"omitempty,password_strength"`
}
```

### 4. Admin Account Creation

```go
type AdminAccount struct {
    Username string `validate:"required,min=4,max=20"`
    Email    string `validate:"required,email"`
    Password string `validate:"required,password_strength,min=12"` // Extra strong
    Role     string `validate:"required,oneof=admin superadmin"`
}
```

## Error Messages

When validation fails, you'll receive specific error messages:

```
password must be at least 8 characters long
password must not exceed 100 characters
password must contain at least one: uppercase letter
password must contain at least one: lowercase letter
password must contain at least one: digit
password must contain at least one: special character (!@#$%^&*()_+-=[]{}|;:,.<>?)
password must contain at least one: uppercase letter, digit, special character (!@#$%^&*()_+-=[]{}|;:,.<>?)
```

## Password Strength Tips

### ✅ Good Practices

1. **Use passphrases**: `MyC@t!sN@med3ob`
2. **Mix character types**: uppercase, lowercase, numbers, symbols
3. **Avoid personal info**: No names, birthdays, common words
4. **Use unique passwords**: Different for each service
5. **Consider length**: Longer is generally stronger

### ❌ Bad Practices

1. **Dictionary words**: `Password123!`
2. **Sequential characters**: `Abc123!@#`
3. **Repeating characters**: `Aaa111!!!`
4. **Keyboard patterns**: `Qwerty123!`
5. **Common substitutions**: `P@ssw0rd`

## Password Validation Levels

### Minimum (8 characters)
```go
"Test@123"  // ✅ Barely meets requirements
```

### Recommended (12+ characters)
```go
"MyS3cure@Pass123"  // ✅ Better security
```

### Strong (16+ characters)
```go
"MyVery$ecure@P@ssw0rd2024!"  // ✅ Excellent security
```

## Integration Examples

### With HTTP Handler

```go
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    var req UserRegistration
    json.NewDecoder(r.Body).Decode(&req)
    
    v, _ := xvalidator.NewValidator()
    if err := v.Validate(req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Password is strong, proceed with registration
}
```

### With Database

```go
func CreateUser(user UserRegistration) error {
    // Validate password strength before hashing
    if err := xvalidator.ValidatePasswordStrength(user.Password); err != nil {
        return fmt.Errorf("weak password: %w", err)
    }
    
    // Hash password and save to database
    hashedPassword := hashPassword(user.Password)
    // ...
}
```

## Security Notes

1. **Always hash passwords** before storing
2. **Never log passwords** in plain text
3. **Use HTTPS** for transmission
4. **Implement rate limiting** for login attempts
5. **Consider multi-factor authentication**
6. **Enforce password expiration** (optional)
7. **Check against breached password lists**

## Next Steps

- [**nested-struct/**](../nested-struct/) - Complex validation scenarios
- [**real-world/**](../real-world/) - See password validation in full applications
- [**custom-messages/**](../custom-messages/) - Customize error messages

## Related Documentation

- [Main README](../../README.md)
- [OWASP Password Guidelines](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
