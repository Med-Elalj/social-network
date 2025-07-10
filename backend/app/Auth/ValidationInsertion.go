package auth

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"social-network/app/modules"
)

func GenerateNickname(firstName, lastName string) string {
	firstName = strings.ToLower(firstName)
	lastName = strings.ToLower(lastName)

	if len(firstName) == 0 || len(lastName) == 0 {
		return ""
	}

	// Try with increasing portions of firstName
	for i := 1; i <= len(firstName); i++ {
		nickname := firstName[:i] + lastName
		if !NicknameExists(nickname) {
			return nickname
		}
	}

	// Fallback: add number suffix
	for suffix := 1; suffix <= 9999; suffix++ {
		nickname := firstName + lastName + strconv.Itoa(suffix)
		if !NicknameExists(nickname) {
			return nickname
		}
	}

	return ""
}

func NicknameExists(nickname string) bool {
	_, exists := EntryExists("display_name", nickname, "profile", true)
	return exists
}

func IsValidNickname(nickname string) bool {
	if len(nickname) < 3 || len(nickname) > 13 {
		return false
	}
	if !regexp.MustCompile(`^[a-zA-Z_]{3,30}$`).MatchString(nickname) {
		return false
	}
	if strings.Contains(nickname, " ") {
		return false
	}
	return true
}

// Validate register form
func (r Register) ValidateRegister() []string {
	var errors []string
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]{3,}@[a-zA-Z0-9.\-]{3,}\.[a-zA-Z]{2,}$`)
	layout := "2006-01-02"

	// Username
	if !IsValidNickname(r.UserName) {
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

// Insert new user
func InsertUser(user Register) (int64, error) {
	tx, err := modules.DB.Begin()
	if err != nil {
		return -1, err
	}

	var avatar sql.NullString
	if avatar.String == "" {
		avatar = sql.NullString{String: "", Valid: false}
	} else {
		avatar = sql.NullString{String: user.Avatar, Valid: true}
	}

	var about sql.NullString
	if user.About == "" {
		about = sql.NullString{String: "", Valid: false}
	} else {
		about = sql.NullString{String: string(user.About), Valid: true}
	}

	// Insert into profile table
	res, err := tx.Exec(`
    INSERT INTO profile (
        email, first_name, last_name, display_name, date_of_birth, gender, avatar, description, is_user
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1)
`,
		user.Email,
		user.Fname,
		user.Lname,
		user.UserName,
		user.Birthdate,
		user.Gender,
		avatar,
		about,
	)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	profileID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// Insert only the password into user table, linking by profile ID
	_, err = tx.Exec(`INSERT INTO user (id, password_hash) VALUES (?, ?)`,
		profileID,
		HashPassword(user.Password),
	)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	if err := tx.Commit(); err != nil {
		return -1, err
	}

	return profileID, nil
}

func hasSpecial(s string) bool {
	for _, ch := range s {
		if unicode.IsPunct(ch) || unicode.IsSymbol(ch) {
			return true
		}
	}
	return false
}

func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 30 {
		return false
	}

	hasDigit := regexp.MustCompile(`[0-9]`)
	hasLower := regexp.MustCompile(`[a-z]`)
	hasUpper := regexp.MustCompile(`[A-Z]`)

	return hasDigit.MatchString(password) &&
		hasLower.MatchString(password) &&
		hasUpper.MatchString(password) &&
		hasSpecial(password)
}
