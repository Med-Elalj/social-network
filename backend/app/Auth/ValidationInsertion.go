package auth

import (
	"database/sql"
	"regexp"
	"time"

	"social-network/app/modules"
)

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
		avatar = sql.NullString{String: string(avatar.String), Valid: true}
	}

	var about sql.NullString
	if user.About == "" {
		about = sql.NullString{String: "", Valid: false}
	} else {
		about = sql.NullString{String: string(user.About), Valid: true}
	}

	res, err := tx.Exec(`INSERT INTO profile (display_name,avatar,description, is_user) VALUES (?,?,?, 1)`, user.UserName, avatar, about)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	profileID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	_, err = tx.Exec(`INSERT INTO user (id, email, first_name, last_name, password_hash, date_of_birth, gender)
	VALUES (?, ?, ?, ?, ?, ?, ?)`,
		profileID, user.Email, user.Fname, user.Lname, HashPassword(user.Password), user.Birthdate, user.Gender)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	if err := tx.Commit(); err != nil {
		return -1, err
	}

	return profileID, nil
}
