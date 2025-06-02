package structs

import (
	"errors"
	"regexp"
	"time"

	"social-network/server/logs"

	"golang.org/x/crypto/bcrypt"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z_]{3,30}$`)

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

// var minAge = 13

func (b Birthdate) IsValid() error {
	t := time.Time(b)
	if t.After(time.Now()) {
		return errors.New("cannot be in the future")
	}
	// Uncomment the following lines if you want to enforce a minimum and max age
	// if t.Year() < 1900 {
	// 	return errors.New("year must be greater than 1900")
	// }
	// if t.Year() > time.Now().Year()-minAge {
	// 	return fmt.Errorf("you must be at least %d years old", minAge)
	// }
	if t.IsZero() {
		return errors.New("cannot be empty")
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
	return l.Password.IsValid()
}

func (p *Password) Hash() {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*p), bcrypt.DefaultCost)
	if err != nil {
		logs.Fatalln(err.Error())
		return
	}
	logs.Printf("Hashing password:%q\ngot: %q", *p, Password(bytes))
	*p = Password(bytes)
}

func (p Password) Verify(password []byte) bool {
	logs.Printf("Verifying password: %q against hash: %q", string(password), string(p))
	if _, err := bcrypt.Cost([]byte(p)); err == nil {
		err := bcrypt.CompareHashAndPassword([]byte(p), password)
		if err != nil {
			logs.Println("Password comparison failed:", err)
			return false
		}
		logs.Println("Password comparison succeeded")
	} else if _, err := bcrypt.Cost([]byte(password)); err == nil {
		err := bcrypt.CompareHashAndPassword([]byte(password), []byte(p))
		if err != nil {
			logs.Println("Password comparison failed:", err)
			return false
		}
		logs.Println("Password comparison succeeded")
	} else {
		logs.Println("Password comparison failed: invalid hash or password format")
		return false
	}
	return true
}
