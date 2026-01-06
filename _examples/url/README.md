# URL Validation Examples

This example demonstrates URL validation for web addresses and API endpoints.

## What You'll Learn

- Basic URL validation (HTTP and HTTPS)
- HTTPS-only URL validation
- Optional URL fields
- URL with query parameters
- API endpoint validation

## Running the Example

```bash
cd url
go run main.go
```

## URL Validators

### Validation Tags

| Tag | Description | Accepts | Example |
|-----|-------------|---------|---------|
| `url` | Any valid URL | HTTP, HTTPS, FTP, etc. | `validate:"url"` |
| `https_url` | HTTPS only | HTTPS only | `validate:"https_url"` |

### When to Use Each

| Validator | Use Case | Example |
|-----------|----------|---------|
| `url` | General websites, flexible protocol | User websites, blog URLs |
| `https_url` | Secure endpoints, APIs, webhooks | API endpoints, payment URLs |

## Code Examples

### 1. Basic URL Validation

```go
type Website struct {
    Homepage string `validate:"required,url"`
    BlogURL  string `validate:"omitempty,url"`
}

// Valid examples
website := Website{
    Homepage: "https://example.com",        // ✅ HTTPS
}

website := Website{
    Homepage: "http://legacy.example.com",  // ✅ HTTP also valid
}
```

### 2. HTTPS-Only Validation

```go
type APIConfig struct {
    BaseURL    string `validate:"required,https_url"`
    WebhookURL string `validate:"required,https_url"`
}

// Valid - HTTPS only
api := APIConfig{
    BaseURL:    "https://api.example.com",     // ✅
    WebhookURL: "https://webhook.example.com", // ✅
}

// Invalid - HTTP not allowed
api := APIConfig{
    BaseURL:    "http://api.example.com",  // ❌ Must be HTTPS
}
```

### 3. Mixed Requirements

```go
type SocialMedia struct {
    Website     string `validate:"omitempty,url"`        // Any protocol
    LinkedInURL string `validate:"omitempty,https_url"`  // HTTPS only
    GitHubURL   string `validate:"omitempty,https_url"`  // HTTPS only
}

// Valid
social := SocialMedia{
    Website:     "http://mysite.com",              // ✅ HTTP OK for website
    LinkedInURL: "https://linkedin.com/in/user",   // ✅ HTTPS required
    GitHubURL:   "https://github.com/user",        // ✅ HTTPS required
}
```

### 4. Optional URL Fields

```go
type UserProfile struct {
    Name    string `validate:"required"`
    Website string `validate:"omitempty,url"`  // Optional but must be valid if provided
}

// Valid - empty optional field
profile := UserProfile{
    Name:    "John Doe",
    Website: "",  // ✅ OK (omitempty)
}

// Valid - with URL
profile := UserProfile{
    Name:    "John Doe",
    Website: "https://johndoe.com",  // ✅ Valid URL
}

// Invalid - malformed URL
profile := UserProfile{
    Name:    "John Doe",
    Website: "not-a-url",  // ❌ Invalid format
}
```

### 5. Single URL Validation

```go
v, _ := xvalidator.NewValidator()

// Validate any URL
url := "https://example.com"
if err := v.Var(url, "url"); err != nil {
    // Invalid URL
}

// Validate HTTPS only
secureURL := "https://api.example.com"
if err := v.Var(secureURL, "https_url"); err != nil {
    // Invalid or not HTTPS
}
```

## Valid URL Examples

### ✅ Valid for `url` Tag

```go
"https://example.com"
"http://example.com"
"https://www.example.com"
"https://subdomain.example.com"
"https://example.com:8080"
"https://example.com/path/to/page"
"https://example.com?query=value"
"https://example.com#section"
"ftp://ftp.example.com"
"https://example.com/api/v1/users?page=1&limit=10"
```

### ✅ Valid for `https_url` Tag

```go
"https://example.com"
"https://www.example.com"
"https://api.example.com"
"https://example.com:443"
"https://example.com/webhook"
"https://example.com/api/v1/callback?token=abc"
"https://192.168.1.1"
"https://[::1]"  // IPv6
```

### ❌ Invalid URLs

```go
// Missing protocol
"example.com"           // ❌ Must have http:// or https://
"www.example.com"       // ❌ Must have protocol

// Malformed
"ht!tp://example.com"   // ❌ Invalid protocol
"https:/example.com"    // ❌ Missing slash
"https://exam ple.com"  // ❌ Space in domain

// Invalid for https_url tag
"http://example.com"    // ❌ Must be HTTPS
"ftp://example.com"     // ❌ Must be HTTPS

// Empty or whitespace
""                      // ❌ Empty string
"   "                   // ❌ Only whitespace
```

## Use Cases

### 1. Website Configuration

```go
type Website struct {
    Homepage string `validate:"required,url"`        // Any protocol OK
    BlogURL  string `validate:"omitempty,url"`       // Optional
    DocsURL  string `validate:"omitempty,url"`       // Optional
}
```

**Why `url`?** User-facing websites might use HTTP for backwards compatibility.

### 2. API Configuration

```go
type APIConfig struct {
    BaseURL     string `validate:"required,https_url"`  // Must be secure
    WebhookURL  string `validate:"required,https_url"`  // Must be secure
    CallbackURL string `validate:"omitempty,https_url"` // Optional but secure
}
```

**Why `https_url`?** APIs should always use HTTPS for security.

### 3. OAuth Configuration

```go
type OAuthConfig struct {
    AuthURL     string `validate:"required,https_url"`
    TokenURL    string `validate:"required,https_url"`
    RedirectURL string `validate:"required,https_url"`
}
```

**Why `https_url`?** OAuth requires HTTPS for security.

### 4. Payment Integration

```go
type PaymentConfig struct {
    PaymentURL string `validate:"required,https_url"`
    WebhookURL string `validate:"required,https_url"`
    ReturnURL  string `validate:"required,https_url"`
}
```

**Why `https_url`?** Payment endpoints must be secure.

## URL Components

```
https://api.example.com:443/v1/users?page=1&limit=10#results
  |      |      |        |    |     |                 |
  |      |      |        |    |     |                 └─ Fragment
  |      |      |        |    |     └─ Query parameters
  |      |      |        |    └─ Path
  |      |      |        └─ Port
  |      |      └─ Domain
  |      └─ Subdomain
  └─ Scheme/Protocol
```

All these components are validated when using `url` or `https_url` tags.

## Error Messages

```
homepage must be a valid URL
base_url must be a valid HTTPS URL
```

## Tips

1. **Use `https_url` for APIs** to enforce security
2. **Use `url` for user-facing websites** for flexibility
3. **Use `omitempty`** for optional URL fields
4. **Store with protocol** - always include `http://` or `https://`
5. **Validate early** - check URLs before making HTTP requests

## Security Considerations

### ✅ Good Practices

- Use `https_url` for sensitive operations
- Validate URLs before making HTTP requests
- Store full URLs with protocol
- Check URL scheme before processing

### ❌ Bad Practices

- Accepting HTTP for sensitive operations
- Auto-prepending protocols without validation
- Assuming all URLs are safe
- Not validating user-provided URLs

## Next Steps

- [**password/**](../password/) - Password strength validation
- [**nested-struct/**](../nested-struct/) - Complex validation
- [**real-world/**](../real-world/) - See URL validation in context

## Related Documentation

- [Main README](../../README.md)
- [URL Standard](https://url.spec.whatwg.org/)
