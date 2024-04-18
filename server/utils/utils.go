package utils

import (
	"fmt"
	"net/mail"
	"reflect"
	"strconv"
	"time"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
    return err == nil
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

func CalculateExpireTimeInt64(minutes string) int64 {
	i, err := strconv.Atoi(minutes)
	if err != nil {
		fmt.Println("[!] Error converting string to int")
		panic(err)
	}

	return time.Now().Add(time.Minute * time.Duration(i)).Unix()
}

func CalculateExpireTime(minutes string) time.Time {
	i, err := strconv.Atoi(minutes)
	if err != nil {
		fmt.Println("[!] Error converting string to int")
		panic(err)
	}

	return time.Now().Add(time.Minute * time.Duration(i))
}