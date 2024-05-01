package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"
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

func ClearJsonFields(jsonArray []interface{}, fieldsToRemove []string) ([]interface{}, error) {
	var modifiedJsonArray []interface{}

	for _, obj := range jsonArray {
		jsonBytes, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}

		var jsonMap map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
			return nil, err
		}

		for _, field := range fieldsToRemove {
			delete(jsonMap, field)
		}

		modifiedJsonBytes, err := json.Marshal(jsonMap)
		if err != nil {
			return nil, err
		}

		var modifiedObj interface{}
		if err := json.Unmarshal(modifiedJsonBytes, &modifiedObj); err != nil {
			return nil, err
		}

		modifiedJsonArray = append(modifiedJsonArray, modifiedObj)
	}

	return modifiedJsonArray, nil
}

func PrintFormattedJSON[T any](items []T) {
	for _, item := range items {
		// Converti l'elemento in JSON.
		jsonData, err := json.MarshalIndent(item, "", " ")
		if err != nil {
			fmt.Println("Errore durante la conversione in JSON:", err)
		}
		// Stampa il JSON formattato.
		fmt.Println(string(jsonData))
	}
}

func DownloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func ReadFileContent(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
