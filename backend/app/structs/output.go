package structs

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mattn/go-sqlite3"
)

func SqlConstraint(err *error) bool {
	if *err != nil {
		if sqliteErr, ok := (*err).(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			if len(sqliteErr.Error()) > 26 && sqliteErr.Error()[:26] == "UNIQUE constraint failed: " {
				e := "sorry"
				for _, v := range strings.Fields(sqliteErr.Error()[26:]) {
					v1 := strings.Split(v, ".")
					if len(v1) > 1 {
						e += " " + v1[1]
					} else {
						e += " " + v
					}
				}
				*err = errors.New(e + " already exists")
				return true
			}
		}
	}
	return false
}

// func JsonRestrictedDecoder(data []byte, destination interface{}) error {
// 	dec := json.NewDecoder(bytes.NewReader(data))
// 	dec.DisallowUnknownFields()
// 	fmt.Println("here : ",dec.Decode(destination))
// 	return dec.Decode(destination)
// }

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func JsRespond(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  code,
	})
}

// validateStruct loops through struct fields and calls IsValid() if implemented
// func validateStruct(v any) error {
// 	val := reflect.ValueOf(v)

// 	fmt.Println(val)

// 	if val.Kind() == reflect.Ptr {
// 		if val.IsNil() {
// 			return errors.New("nil pointer passed to ValidateStruct")
// 		}
// 		val = val.Elem()
// 	}

// 	if val.Kind() != reflect.Struct {
// 		return errors.New("ValidateStruct expects a struct or pointer to struct")
// 	}

// 	typ := val.Type()
// 	for i := 0; i < val.NumField(); i++ {
// 		field := val.Field(i)
// 		fieldType := typ.Field(i)

// 		// Only check exported fields
// 		if !field.CanInterface() {
// 			continue
// 		}

// 		// Get the JSON tag (fall back to field name if not set)
// 		jsonTag := fieldType.Tag.Get("json")
// 		if jsonTag == "" {
// 			jsonTag = fieldType.Name
// 		}

// 		// Check if the field implements the Input interface
// 		if input, ok := field.Interface().(Input); ok {
// 			if err := input.IsValid(); err != nil {
// 				return fmt.Errorf("field '%s' is invalid: %w", jsonTag, err)
// 			}
// 		}
// 	}

// 	return nil
// }
