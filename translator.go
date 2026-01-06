package xvalidator

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
)

// setupTranslator creates and configures an English translator for validation messages
func setupTranslator(v *validator.Validate) (ut.Translator, error) {
	// Setup English translator
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	// Register default English translations
	err := en_trans.RegisterDefaultTranslations(v, trans)
	if err != nil {
		return nil, fmt.Errorf("failed to register default translations: %w", err)
	}

	// Register custom translations for our custom validators
	err = registerCustomTranslations(v, trans)
	if err != nil {
		return nil, fmt.Errorf("failed to register custom translations: %w", err)
	}
	return trans, nil
}

// formatTranslatedErrors converts validator errors to user-friendly translated messages
func formatTranslatedErrors(validationErrors validator.ValidationErrors, translator ut.Translator) error {
	var messages []string
	for _, err := range validationErrors {
		translatedMsg := err.Translate(translator)
		messages = append(messages, translatedMsg)
	}
	return fmt.Errorf("%s", strings.Join(messages, "; "))
}

// registerDecimalTranslation registers decimal validation translation with custom formatting
func registerDecimalTranslation(v *validator.Validate, trans ut.Translator) error {
	// Register main decimal translation
	err := v.RegisterTranslation("decimal", trans, func(ut ut.Translator) error {
		return ut.Add("decimal", "{0} must be a decimal with precision ≤ {1} and scale ≤ {2}", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		if param == "" {
			// Use default values when no parameter specified
			translated, _ := ut.T("decimal", fe.Field(),
				fmt.Sprintf("%d", DefaultPrecision),
				fmt.Sprintf("%d", DefaultScale))
			return translated
		}

		// Parse parameters to get precision and scale
		precision, scale := parseDecimalParams(param)

		// Special case for integer format (scale = 0)
		if scale == 0 {
			return fmt.Sprintf("%s must be an integer format (no decimal places)", fe.Field())
		}

		translated, _ := ut.T("decimal", fe.Field(),
			fmt.Sprintf("%d", precision),
			fmt.Sprintf("%d", scale))
		return translated
	})
	if err != nil {
		return fmt.Errorf("failed to register decimal translation: %w", err)
	}

	return nil
}

// registerDecimalIfTranslation registers decimal_if validation translation with custom formatting
func registerDecimalIfTranslation(v *validator.Validate, trans ut.Translator) error {
	// Register main decimal_if translation
	err := v.RegisterTranslation("decimal_if", trans, func(ut ut.Translator) error {
		return ut.Add("decimal_if", "{0} must be a decimal with precision ≤ {1} and scale ≤ {2} when {3} equals '{4}'", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		if param == "" {
			// Use default error message when no parameter specified
			return fmt.Sprintf("%s conditional decimal validation failed", fe.Field())
		}

		// Parse parameters to get rule, field, and expected value
		rule, field, expect, err := parseDecimalIfParam(param)
		if err != nil {
			return fmt.Sprintf("%s conditional decimal validation failed", fe.Field())
		}

		// Parse decimal rule to get precision and scale
		precision, scale := parseDecimalParams(rule)

		// Special case for integer format (scale = 0)
		if scale == 0 {
			return fmt.Sprintf("%s must be an integer format (no decimal places) when %s equals '%s'",
				fe.Field(), field, expect)
		}

		// Check if we have a specific rule or using defaults
		if rule == "" {
			return fmt.Sprintf("%s must be a decimal with default precision and scale when %s equals '%s'",
				fe.Field(), field, expect)
		}

		translated, _ := ut.T("decimal_if", fe.Field(),
			fmt.Sprintf("%d", precision),
			fmt.Sprintf("%d", scale),
			field, expect)
		return translated
	})
	if err != nil {
		return fmt.Errorf("failed to register decimal_if translation: %w", err)
	}

	return nil
}

// registerPasswordStrengthTranslation registers password_strength validation translation with custom formatting
func registerPasswordStrengthTranslation(v *validator.Validate, trans ut.Translator) error {
	// Define special characters as constant to avoid escaping issues
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	// Register password_strength translation without parameter placeholders
	err := v.RegisterTranslation("password_strength", trans, func(ut ut.Translator) error {
		return ut.Add("password_strength", "must contain at least 8 characters with: uppercase letter (A-Z), lowercase letter (a-z), digit (0-9), and special character", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		// Build message with special characters defined separately
		return fmt.Sprintf("%s must contain at least 8 characters with: uppercase letter (A-Z), lowercase letter (a-z), digit (0-9), and special character (%s)", fe.Field(), specialChars)
	})
	if err != nil {
		return fmt.Errorf("failed to register password_strength translation: %w", err)
	}

	return nil
}

// registerCustomTranslations registers English translations for our custom validators
func registerCustomTranslations(v *validator.Validate, trans ut.Translator) error {
	// Register decimal translations first
	err := registerDecimalTranslation(v, trans)
	if err != nil {
		return err
	}

	// Register decimal_if translation
	err = registerDecimalIfTranslation(v, trans)
	if err != nil {
		return err
	}

	// Register password_strength translation
	err = registerPasswordStrengthTranslation(v, trans)
	if err != nil {
		return err
	}

	// Register translations for other validators
	translations := map[string]struct {
		tag         string
		translation string
		override    bool
	}{
		"dgt": {
			tag:         "dgt",
			translation: "{0} must be greater than {1}",
			override:    false,
		},
		"dgte": {
			tag:         "dgte",
			translation: "{0} must be greater than or equal to {1}",
			override:    false,
		},
		"dlt": {
			tag:         "dlt",
			translation: "{0} must be less than {1}",
			override:    false,
		},
		"dlte": {
			tag:         "dlte",
			translation: "{0} must be less than or equal to {1}",
			override:    false,
		},
		"deq": {
			tag:         "deq",
			translation: "{0} must be equal to {1}",
			override:    false,
		},
		"dneq": {
			tag:         "dneq",
			translation: "{0} must not be equal to {1}",
			override:    false,
		},
		"https_url": {
			tag:         "https_url",
			translation: "{0} must be a valid HTTPS URL",
			override:    false,
		},
		"mobile_e164": {
			tag:         "mobile_e164",
			translation: "{0} must be a valid mobile number in E.164 format (e.g., +66812345678)",
			override:    false,
		},
		"iso4217": {
			tag:         "iso4217",
			translation: "{0} must be a valid ISO 4217 currency code (e.g., THB, USD, EUR)",
			override:    false,
		},
	}

	for _, t := range translations {
		err := v.RegisterTranslation(t.tag, trans, func(ut ut.Translator) error {
			return ut.Add(t.tag, t.translation, t.override)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			if fe.Param() != "" {
				translated, _ := ut.T(t.tag, fe.Field(), fe.Param())
				return translated
			}
			translated, _ := ut.T(t.tag, fe.Field())
			return translated
		})
		if err != nil {
			return fmt.Errorf("failed to register translation for %s: %w", t.tag, err)
		}
	}
	return nil
}
