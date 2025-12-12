package masker

import (
	"reflect"
	"strings"

	"github.com/ggwhite/go-masker"
)

var sensitiveKeys = map[string]struct{}{
	"password":      {},
	"token":         {},
	"secret":        {},
	"authorization": {},
	"api_key":       {},
	"access_token":  {},
	"refresh_token": {},
	"pin":           {},
	"cvv":           {},
	"pan":           {},
	"credit_card":   {},
	"email":         {},
	"phone":         {},
	"mobile":        {},
}

// Mask recursively masks sensitive fields in the input value.
// It returns a copy of the input with sensitive fields masked.
// Structs are converted to map[string]any to allow replacing non-string fields with masked strings.
func Mask(input any) any {
	if input == nil {
		return nil
	}
	return maskValue(reflect.ValueOf(input))
}

func maskValue(val reflect.Value) any {
	if !val.IsValid() {
		return nil
	}

	switch val.Kind() {
	case reflect.Pointer:
		if val.IsNil() {
			return nil
		}
		// Dereference and continue
		return maskValue(val.Elem())
	case reflect.Interface:
		if val.IsNil() {
			return nil
		}
		return maskValue(val.Elem())
	case reflect.Struct:
		return maskStruct(val)
	case reflect.Map:
		return maskMap(val)
	case reflect.Slice, reflect.Array:
		return maskSlice(val)
	default:
		// Primitive values or others are returned as is
		return val.Interface()
	}
}

func maskStruct(val reflect.Value) map[string]any {
	out := make(map[string]any)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		if !structField.IsExported() {
			continue
		}

		fieldName := structField.Name
		jsonTag := structField.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			parts := strings.Split(jsonTag, ",")
			if len(parts) > 0 {
				fieldName = parts[0]
			}
		}

		// Check if sensitive
		lowerName := strings.ToLower(fieldName)
		if _, isSensitive := sensitiveKeys[lowerName]; isSensitive {
			out[fieldName] = applyMask(field, lowerName)
		} else {
			out[fieldName] = maskValue(field)
		}
	}
	return out
}

func maskMap(val reflect.Value) map[string]any {
	out := make(map[string]any)
	iter := val.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()

		keyStr := k.String()
		lowerKey := strings.ToLower(keyStr)

		if _, isSensitive := sensitiveKeys[lowerKey]; isSensitive {
			out[keyStr] = applyMask(v, lowerKey)
		} else {
			out[keyStr] = maskValue(v)
		}
	}
	return out
}

func maskSlice(val reflect.Value) []any {
	l := val.Len()
	out := make([]any, l)
	for i := range l {
		out[i] = maskValue(val.Index(i))
	}
	return out
}

func applyMask(val reflect.Value, key string) any {
	// Convert value to string for masking
	var strVal string

	switch val.Kind() {
	case reflect.String:
		strVal = val.String()
	default:
		// If it's not a string (e.g. int PIN), we can't easily use go-masker specific functions that expect strings
		// or if we do, we assume the value *should* be treated as a string.
		// For simplicity, let's just use a generic mask for non-string types or convert to string.
		// But go-masker works on strings.
		// Let's rely on any being string?
		if val.CanInterface() {
			if s, ok := val.Interface().(string); ok {
				strVal = s
			} else {
				// Fallback for non-string sensitive data (like int PIN)
				return "******"
			}
		}
	}

	// Use go-masker based on key type
	// Note: go-masker API is usually masker.IDCard(s), masker.Mobile(s), etc.
	// We need to map keys to masker functions.

	switch key {
	case "email":
		return masker.Email(strVal)
	case "mobile", "phone":
		return masker.Mobile(strVal)
	case "password", "token", "secret", "access_token", "refresh_token", "api_key":
		return masker.Password(strVal) // Or generic mask if Password not available, but Password usually exists or we use Overlay
	case "cvv", "pin":
		// short secrets
		return "******"
	default:
		// Generic masking with overlay
		// masker.Overlay(str, start, end, mask)
		// Default to password-like mask
		return masker.Password(strVal)
	}
}
