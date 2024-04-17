package utils

import (
	"fmt"
	"reflect"
	"net/mail"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
    return err == nil
}

func Padding(s string, length int) string {
	for len(s) < length {
		s = s + "="
	}
	return s
}

func UnPadding(s string, length int) string {
	for len(s) > length {
		s = s[:len(s)-1]
	}
	return s
}

func CheckAllFieldsHaveValue(s interface{}) bool {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				return false
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() == 0 {
				return false
			}
		case reflect.Ptr:
			if field.IsNil() {
				return false
			}
		default:
			fmt.Printf("Tipo di campo non gestito: %s\n", field.Type())
		}
	}

	return true
}