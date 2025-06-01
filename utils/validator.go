package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var messages []string
	for _, err := range v {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, ", ")
}

// ValidateStruct validates a struct based on validation tags
func ValidateStruct(s interface{}) ValidationErrors {
	var errors ValidationErrors
	
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	
	// Handle pointer to struct
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	
	if v.Kind() != reflect.Struct {
		return errors
	}
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		// Get validation tag
		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}
		
		// Get JSON tag for field name
		jsonTag := fieldType.Tag.Get("json")
		fieldName := fieldType.Name
		if jsonTag != "" && jsonTag != "-" {
			fieldName = strings.Split(jsonTag, ",")[0]
		}
		
		// Parse validation rules
		rules := strings.Split(validateTag, ",")
		for _, rule := range rules {
			if err := validateField(field, fieldName, rule); err != nil {
				errors = append(errors, *err)
			}
		}
	}
	
	return errors
}

// validateField validates a single field based on a rule
func validateField(field reflect.Value, fieldName, rule string) *ValidationError {
	parts := strings.Split(rule, "=")
	ruleName := parts[0]
	var ruleValue string
	if len(parts) > 1 {
		ruleValue = parts[1]
	}
	
	switch ruleName {
	case "required":
		if isEmptyValue(field) {
			return &ValidationError{
				Field:   fieldName,
				Message: "is required",
			}
		}
	case "min":
		if minVal, err := strconv.Atoi(ruleValue); err == nil {
			if field.Kind() == reflect.String {
				if len(field.String()) < minVal {
					return &ValidationError{
						Field:   fieldName,
						Message: fmt.Sprintf("must be at least %d characters", minVal),
					}
				}
			}
		}
	case "max":
		if maxVal, err := strconv.Atoi(ruleValue); err == nil {
			if field.Kind() == reflect.String {
				if len(field.String()) > maxVal {
					return &ValidationError{
						Field:   fieldName,
						Message: fmt.Sprintf("must be at most %d characters", maxVal),
					}
				}
			}
		}
	case "email":
		if field.Kind() == reflect.String {
			email := field.String()
			if email != "" && !isValidEmail(email) {
				return &ValidationError{
					Field:   fieldName,
					Message: "must be a valid email address",
				}
			}
		}
	}
	
	return nil
}

// isEmptyValue checks if a value is empty
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}