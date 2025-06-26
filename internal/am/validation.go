package am

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ValidationField represents a field-specific validation error.
type ValidationField struct {
	Name    string
	Value   string
	Message string
}

// Validation holds the result of one or many validation checks.
// It can collect multiple error messages unless strict mode is enabled.
// Use IsValid() to check result, and Error() or JSON() to report.
type Validation struct {
	Errors []string
	Fields map[string]ValidationField
}

// Add adds an error message to the validation result.
func (v *Validation) Add(err string) {
	v.Errors = append(v.Errors, err)
}

// AddFieldError adds a field-specific error to the validation result.
func (v *Validation) AddFieldError(field, value, message string) {
	v.Errors = append(v.Errors, message)
	if v.Fields == nil {
		v.Fields = make(map[string]ValidationField)
	}
	v.Fields[field] = ValidationField{
		Name:    field,
		Value:   value,
		Message: message,
	}
}

// IsValid returns true if there are no validation errors.
func (v Validation) IsValid() bool {
	return len(v.Errors) == 0
}

// HasErrors returns true if there are validation errors.
// This is the inverse of IsValid and can make code more readable by avoiding negations.
func (v Validation) HasErrors() bool {
	return !v.IsValid()
}

// Error returns all validation errors as a comma-separated string.
func (v Validation) Error() string {
	return strings.Join(v.Errors, ", ")
}

// JSON returns all validation errors as a JSON string.
func (v Validation) JSON() string {
	data, _ := json.Marshal(v.Errors)
	return string(data)
}

// FieldMsg returns the message for a given field tag name, or an empty string if none exists.
func (v Validation) FieldMsg(field string) string {
	if v.Fields == nil {
		return ""
	}

	if f, ok := v.Fields[field]; ok {
		return f.Message
	}

	return ""
}

// Validator is a function type that performs validation on any input.
type Validator func(v any) (Validation, error)

// ComposeValidators allows combining multiple Validator functions into one.
// By default, it will collect all validation errors.
// If strict is true, it will stop after the first error.
func ComposeValidators(fns ...Validator) Validator {
	return func(v any) (Validation, error) {
		out := Validation{}
		for _, fn := range fns {
			res, err := fn(v)
			if err != nil {
				return out, err
			}
			out.Errors = append(out.Errors, res.Errors...)
			if res.Fields != nil {
				if out.Fields == nil {
					out.Fields = make(map[string]ValidationField)
				}
				for k, field := range res.Fields {
					out.Fields[k] = field
				}
			}
		}
		return out, nil
	}
}

// ComposeValidatorsStrict allows combining multiple Validator functions into one
// and stops after the first error is found.
func ComposeValidatorsStrict(fns ...Validator) Validator {
	return func(v any) (Validation, error) {
		out := Validation{}
		for _, fn := range fns {
			res, err := fn(v)
			if err != nil {
				return out, err
			}
			out.Errors = append(out.Errors, res.Errors...)
			if res.Fields != nil {
				if out.Fields == nil {
					out.Fields = make(map[string]ValidationField)
				}
				for k, field := range res.Fields {
					out.Fields[k] = field
				}
			}
			if len(out.Errors) > 0 {
				break
			}
		}
		return out, nil
	}
}

// Common validation primitives:

// MinLength validates that a string field has a minimum length.
func MinLength(field, val string, min int) Validator {
	return func(_ any) (Validation, error) {
		v := Validation{}
		if len(val) < min {
			v.AddFieldError(field, val, fmt.Sprintf("%s: must be at least %d characters", field, min))
		}
		return v, nil
	}
}

// MaxLength validates that a string field has a maximum length.
func MaxLength(field, val string, max int) Validator {
	return func(_ any) (Validation, error) {
		v := Validation{}
		if len(val) > max {
			v.AddFieldError(field, val, fmt.Sprintf("%s: must be at most %d characters", field, max))
		}
		return v, nil
	}
}

// Equals validates that two string fields are equal.
func Equals(field string, a, b string) Validator {
	return func(_ any) (Validation, error) {
		v := Validation{}
		if a != b {
			v.Add(fmt.Sprintf("%s: fields do not match", field))
		}
		return v, nil
	}
}

// GreaterThan validates that an integer field is greater than a value.
func GreaterThan(field string, a, b int) Validator {
	return func(_ any) (Validation, error) {
		v := Validation{}
		if a <= b {
			v.Add(fmt.Sprintf("%s must be greater than %d", field, b))
		}
		return v, nil
	}
}
