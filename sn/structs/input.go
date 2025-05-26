package structs

import (
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

type Gender int

// User struct with custom types
type Register struct {
	UserName  Name      `json:"username"`
	Email     Email     `json:"email"`
	Birthdate Birthdate `json:"birthdate"`
	Fname     Name      `json:"fname"`
	Lname     Name      `json:"lname"`
	Password  Password  `json:"password"`
	Gender    Gender    `json:"gender"`
}

type Login struct {
	NoE      Validator `json:"login"`
	Password Password  `json:"pwd"`
}

// Input interface with IsValid method
type Input interface {
	IsValid() error
}

type NameOrEmail struct {
	Validator
}

var nameRegex = regexp.MustCompile(`^[a-zA-Z_]{3,}$`)

func (n Name) IsValid() error {
	if len(n) == 0 {
		return errors.New("cannot be empty")
	} else if !nameRegex.MatchString(string(n)) {
		return errors.New("invalid charachters used")
	}
	return nil
}

func (e Email) IsValid() error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if emailRegex.MatchString(string(e)) {
		return nil
	} else {
		return errors.New("malformed or missing data")
	}
}

func (b Birthdate) IsValid() error {
	// TODO: max age
	t := time.Time(b)
	if t.After(time.Now()) {
		return errors.New("cannot be in the future")
	}
	return nil
}

func (p Password) IsValid() error {
	if len(string(p)) < 8 {
		return errors.New("password must be at least 8 characters long")
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
