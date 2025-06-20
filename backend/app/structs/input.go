package structs

import (
	"regexp"
	"time"

	"social-network/server/logs"

	"golang.org/x/crypto/bcrypt"
)

// verify password
func (p Password) Verify(password []byte) bool {
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

// Validate register form
func (r Register) ValidateRegister() []string {
	var errors []string
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]{3,}@[a-zA-Z0-9.\-]{3,}\.[a-zA-Z]{2,}$`)
	nameRegex := regexp.MustCompile(`^[a-zA-Z_]{3,30}$`)
	layout := "2006-01-02"

	// Username
	if len(r.UserName) < 3 || len(r.UserName) > 13 || !nameRegex.MatchString(string(r.UserName)) {
		errors = append(errors, "Username must be 3-13 characters and use letters or underscores.")
	}

	// Email
	if !emailRegex.MatchString(string(r.Email)) {
		errors = append(errors, "Invalid email format.")
	}

	// Birthdate
	birthdate, err := time.Parse(layout, r.Birthdate)
	if err != nil {
		errors = append(errors, "Birthdate cannot be empty.")
	}
	if birthdate.IsZero() {
		errors = append(errors, "Birthdate cannot be empty.")
	} else if birthdate.After(time.Now()) {
		errors = append(errors, "Birthdate cannot be in the future.")
	} else if birthdate.Year() < 1900 {
		errors = append(errors, "Birth year must be greater than 1900.")
	} else if birthdate.Year() > time.Now().Year()-13 {
		errors = append(errors, "You must be at least 13 years old.")
	}

	// First name
	if len(r.Fname) < 3 || len(r.Fname) > 13 {
		errors = append(errors, "First name must be 3-13 characters.")
	}

	// Last name
	if len(r.Lname) < 3 || len(r.Lname) > 13 {
		errors = append(errors, "Last name must be 3-13 characters.")
	}

	// Password
	if len(r.Password) < 8 {
		errors = append(errors, "Password must be at least 8 characters long.")
	}

	if r.Gender != "male" && r.Gender != "female" {
		errors = append(errors, "Invalid gender. Must be male or female.")
	}

	return errors
}

//Generate hash code
func (p *Password) Hash() {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*p), bcrypt.DefaultCost)
	if err != nil {
		logs.Fatalln(err.Error())
		return
	}
	logs.Printf("Hashing password:%q\ngot: %q", *p, Password(bytes))
	*p = Password(bytes)
}