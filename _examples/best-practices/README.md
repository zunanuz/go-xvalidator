# Best Practices

This folder demonstrates best practices and recommended patterns for using go-xvalidator in production applications.

## Examples

### 1. decimal-workflow.go
Complete workflow showing how to properly use string validation with decimal calculations:
- Accept string input from API/user
- Validate using `decimal` validators (dgt, dgte, dlt, dlte, deq, dneq)
- Convert to `decimal.Decimal` for accurate calculations
- Return string output

This pattern prevents floating-point precision errors and maintains consistency between validation and business logic.

## Key Principles

1. **Use strings for monetary values** - Accept as string, validate, then convert to decimal.Decimal internally
2. **Validate early** - Validate input before processing business logic
3. **Use custom validators** - Create reusable validators for domain-specific rules
4. **Translate messages** - Provide localized error messages for better UX
5. **Leverage struct tags** - Use comprehensive validation tags for complex objects
6. **Handle edge cases** - Test zero values, negative values, boundary conditions

## Running Examples

```bash
cd best-practices
go run decimal-workflow.go
```
