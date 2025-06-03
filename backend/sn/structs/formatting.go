package structs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"social-network/server/logs"
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
	for i, t := range []string{"male", "female", "Attack Helicopter"} {
		if raw == t {
			*g = Gender(i)
			return nil
		}
	}

	return fmt.Errorf("invalid gender %q valid ones are ['male','female','Attack Helicopter']", raw)
}

func (g Gender) MarshalJSON() ([]byte, error) {
	names := []string{"male", "female", "Attack Helicopter"}

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
		n.Input = email
	} else if Name(s).IsValid() == nil {
		name := Name(s)
		n.Input = name
	} else {
		logs.Println("Unmarshalling NameOrEmail:", s, "||||", Email(s).IsValid(), Name(s).IsValid(), "||||", n.Input)
		return errors.New("not a valid email or name")
	}
	return nil
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

func (pc Pcategories) Value() (driver.Value, error) {
	if pc == nil {
		logs.Fatal("Birthdate is zero value, cannot store in DB")
		return nil, errors.New("birthdate is zero value, cannot store in DB")
	}
	return strings.Join(pc, "|"), nil // Store date only, no time
}

// Implement the sql.Scanner interface for reading from DB
func (pc *Pcategories) Scan(value interface{}) error {
	if value == nil {
		*pc = Pcategories{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*pc = strings.Split(string(v), "|")
		return nil
	case string:
		*pc = strings.Split(v, "|")
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Categories", value)
	}
}

func (n NameOrEmail) Value() (driver.Value, error) {
	switch v := n.Input.(type) {
	case Name:
		return string(v), nil
	case Email:
		return string(v), nil
	default:
		return nil, fmt.Errorf("unsupported type %T for NameOrEmail", v)
	}
}

func (p *PostPrivacy) UnmarshalJSON(data []byte) error {
	// Unmarshal as string
	var raw string
	logs.Println("Unmarshalling PostPrivacy:", string(data))
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid post Privacy format: %w", err)
	}

	// Parse using expected format (customize if needed)
	for i, t := range []string{"public", "followers", "private"} {
		if raw == t {
			*p = PostPrivacy(i)
			return nil
		}
	}

	return fmt.Errorf("invalid gender %q valid ones are ['public', 'followers', 'private']", raw)
}

func (p PostPrivacy) MarshalJSON() ([]byte, error) {
	names := []string{"public", "followers", "private"}

	if int(p) < 0 || int(p) >= len(names) {
		return nil, fmt.Errorf("invalid privacy value: %d", p)
	}

	// Marshal the string as JSON string literal
	return json.Marshal(names[int(p)])
}
