package structs

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"social-network/server/logs"

	"golang.org/x/crypto/bcrypt"
)

func (b *Birthdate) UnmarshalJSON(data []byte) error {
	// Unmarshal as string
	var raw string
	logs.Println("Unmarshalling Birthdate:", string(data))
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid birthdate format: %w", err)
	}
	if raw == "" {
		return fmt.Errorf("empty feild (birthdate) expected format (YYYY-MM-DD)")
	}

	// Parse using expected format (customize if needed)
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return fmt.Errorf("invalid birthdate format (expected YYYY-MM-DD): %w", err)
	}

	*b = Birthdate(parsed)
	return nil
}

func (g *Gender) UnmarshalJSON(data []byte) error {
	// Unmarshal as string
	var raw string
	logs.Println("Unmarshalling Gender:", string(data))
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid gender format: %w", err)
	}

	// Parse using expected format (customize if needed)
	for i, t := range []string{"male", "female", "DFK"} {
		if raw == t {
			*g = Gender(i)
			return nil
		}
	}

	return fmt.Errorf("invalid gender %q valid ones are ['male','female','DFK']", raw)
}

func (g Gender) MarshalJSON() ([]byte, error) {
	names := []string{"male", "female", "DFK"}

	if int(g) < 0 || int(g) >= len(names) {
		return nil, fmt.Errorf("invalid gender value: %d", g)
	}

	// Marshal the string as JSON string literal
	return json.Marshal(names[g])
}

func (n *NameOrEmail) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// simple email detection
	if Email(s).IsValid() == nil {
		email := Email(s)
		n.Validator = email
	} else if Name(s).IsValid() == nil {
		name := Name(s)
		n.Validator = name
	}
	return errors.New("not a valid email or name")
}

func (n NameOrEmail) String() string {
	switch v := n.Validator.(type) {
	case Name:
		return string(v)
	case Email:
		return string(v)
	default:
		return ""
	}
}

func JsonRestrictedDecoder(data []byte, destination interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(destination)
}

// Implement the driver.Valuer interface for writing to DB
func (b Birthdate) Value() (driver.Value, error) {
	t := time.Time(b)
	if t.IsZero() {
		logs.Fatal("Birthdate is zero value, cannot store in DB")
		return nil, errors.New("birthdate is zero value, cannot store in DB")
	}
	return t.Format("2006-01-02"), nil // Store date only, no time
}

// Implement the sql.Scanner interface for reading from DB
func (b *Birthdate) Scan(value interface{}) error {
	if value == nil {
		*b = Birthdate(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*b = Birthdate(v)
		return nil
	case []byte:
		// Parse date string from DB, usually in "YYYY-MM-DD"
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		*b = Birthdate(t)
		return nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		*b = Birthdate(t)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Birthdate", value)
	}
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
	if _, err := bcrypt.Cost([]byte(p)); err != nil {
		err := bcrypt.CompareHashAndPassword([]byte(p), password)
		if err != nil {
			logs.Println("Password comparison failed:", err)
			return false
		}
		logs.Println("Password comparison succeeded")
	} else if _, err := bcrypt.Cost([]byte(password)); err != nil {
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
