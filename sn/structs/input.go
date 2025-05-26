package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"
)

type Validator interface {
	IsValid() error
}

// Custom field types implementing Validator

type Name string

type Email string

type Password string

type Birthdate time.Time

// User struct with custom types
type Register struct {
	UserName  Name      `json:"username"`
	Email     Email     `json:"email"`
	Birthdate Birthdate `json:"birthdate"`
	Fname     Name      `json:"fname"`
	Lname     Name      `json:"lname"`
	Password  Password  `json:"password"`
	// Username  string    `json:"username"`
	// Email     string    `json:"email"`
	// Age       string    `json:"age"`
	// Gender    string    `json:"gender"`
	// Avatar    string    `json:"avatar"`
	// Aboutme   string    `json:"aboutme"`
	// Status    string    `json:"status"`
}

type Login struct {
	Username Name
	Password Name
}

// Input interface with IsValid method
type Input interface {
	IsValid() error
}

// Name type implementing Input
// type Name string
var nameRegex = regexp.MustCompile(`[a-zA-Z_]{3,}$`)

func (n Name) IsValid() error {
	if len(n) == 0 {
		return errors.New("cannot be empty")
	}
	return nil
}

// Email type implementing Input
// type Email string

func (e Email) IsValid() error {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if emailRegex.MatchString(string(e)) {
		return nil
	} else {
		return errors.New("malformed or missing data")
	}
}

// Birthdate type implementing Input
// type Birthdate time.Time

func (b Birthdate) IsValid() error {
	t := time.Time(b)
	if t.After(time.Now()) {
		return errors.New("cannot be in the future")
	}
	return nil
}

var passwordRegex = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{8,}$`)

func (p Password) IsValid() error {
	pass := string(p)

	if !passwordRegex.MatchString(pass) {
		return errors.New("password must be at least 8 characters long, and include at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}

func (r Register) Validate() error {
	return validateStruct(r)
}

func (l Login) Validate() error {
	return validateStruct(l)
}

// validateStruct loops through struct fields and calls IsValid() if implemented
func validateStruct(v any) error {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return errors.New("nil pointer passed to ValidateStruct")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.New("ValidateStruct expects a struct or pointer to struct")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Only check exported fields
		if !field.CanInterface() {
			continue
		}

		// Get the JSON tag (fall back to field name if not set)
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = fieldType.Name
		}

		// Check if the field implements the Input interface
		if input, ok := field.Interface().(Input); ok {
			if err := input.IsValid(); err != nil {
				return fmt.Errorf("field '%s' is invalid: %w", jsonTag, err)
			}
		}
	}

	return nil
}

func (b *Birthdate) UnmarshalJSON(data []byte) error {
	// Unmarshal as string
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid birthdate format: %w", err)
	}

	// Parse using expected format (customize if needed)
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return fmt.Errorf("invalid birthdate format (expected YYYY-MM-DD): %w", err)
	}

	*b = Birthdate(parsed)
	return nil
}
