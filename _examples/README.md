# go-xvalidator Examples

Welcome to the examples directory! This contains comprehensive examples demonstrating all features of go-xvalidator.

## 📚 Table of Contents

### Getting Started

- [**basic/**](basic/) - Basic validation examples
  - Required fields, min/max length, email validation
  - String, number, and basic type validation

### Feature-Specific Examples

- [**decimal/**](decimal/) - Decimal validation
  - Greater than, less than, equal comparisons
  - Conditional decimal validation
  - Precision and scale validation
  - **Note**: Contains 2 files (main.go and conditional.go) - run separately

- [**phone/**](phone/) - Phone number validation
  - E.164 format validation
  - Country-specific validation (TH, US, GB, FR, etc.)

- [**url/**](url/) - URL validation
  - Standard URL validation
  - HTTPS-only URL validation

- [**password/**](password/) - Password strength validation
  - Complex password requirements
  - Strength validation rules

### Advanced Examples

- [**nested-struct/**](nested-struct/) - Complex struct validation
  - Nested struct validation
  - Slice and array validation
  - Deep validation examples

- [**custom-messages/**](custom-messages/) - Custom error messages
  - Customizing error messages
  - Translation setup
  - Multi-language support

- [**advanced/**](advanced/) - Advanced techniques
  - Creating custom validators (Thai ID, Business Hours, Thai Phone, etc.)
  - Custom validation functions with parameters
  - Validator factory patterns

- [**real-world/**](real-world/) - Real-world scenarios
  - User registration forms (user-registration.go)
  - Payment processing validation (payment.go)
  - E-commerce order validation (ecommerce.go)
  - **Note**: Contains 3 separate files - run each individually

- [**best-practices/**](best-practices/) - Best practices & patterns
  - String validation + decimal.Decimal calculation workflow
  - API layer design patterns
  - Accurate monetary calculations

## 🚀 Running Examples

Each example folder contains a `main.go` file that can be run independently:

```bash
# Run basic example
cd basic
go run main.go

# Run decimal examples (2 files)
cd decimal
go run main.go
go run conditional.go

# Run real-world examples (3 files)
cd real-world
go run user-registration.go
go run payment.go
go run ecommerce.go

# Run any specific example
cd <example-folder>
go run main.go
```

## 📖 Example Structure

Each example folder contains:

- `main.go` - Main example code (or multiple .go files)
- `README.md` - Detailed explanation
- Additional `.go` files for specific scenarios (when applicable)

## 🔗 Related Resources

- [Main README](../README.md) - Package documentation
- [API Documentation](https://pkg.go.dev/github.com/hotfixfirst/go-xvalidator)
- [GitHub Repository](https://github.com/hotfixfirst/go-xvalidator)

## 💡 Tips

1. Start with [basic/](basic/) if you're new to go-xvalidator
2. Read [best-practices/](best-practices/) to understand recommended patterns
3. Each example is self-contained and can be run independently
4. Check the README in each folder for detailed explanations
5. Examples are ordered from simple to advanced
6. Some folders have multiple files - see notes above for which ones

## ⚠️ Important Notes

- **decimal/** folder has 2 files: `main.go` and `conditional.go` - run separately
- **real-world/** folder has 3 files - run each individually (not together)
- All examples use `string` type for monetary values (API layer)
- Convert to `decimal.Decimal` internally for calculations (see best-practices/)

## 🤝 Contributing

Have a useful example to share? Contributions are welcome! Please submit a PR with:

- Working code example
- Clear documentation
- Explanation of use case

---

**Happy Validating! 🎉**
