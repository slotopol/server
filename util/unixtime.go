package util

import (
	"encoding/json"
	"encoding/xml"
	"time"

	"gopkg.in/yaml.v3"
)

type Unix int64 // UNIX time in seconds.

const YearOnly = "2006"

func (u Unix) Time() time.Time {
	return time.Unix(int64(u), 0).UTC()
}

func (u Unix) Year() int {
	return u.Time().Year()
}

func (u Unix) IsZero() bool {
	return u == 0
}

// MarshalJSON is standard JSON interface implementation to stream UNIX time.
func (u Unix) MarshalJSON() ([]byte, error) {
	var t = u.Time()
	if t.Hour() != 0 || t.Minute() != 0 || t.Second() != 0 {
		return json.Marshal(t.Format(time.DateTime))
	}
	if t.Month() != 1 || t.Day() != 1 {
		return json.Marshal(t.Format(time.DateOnly))
	}
	return json.Marshal(t.Format(YearOnly))
}

// UnmarshalJSON is standard JSON interface implementation to stream UNIX time.
func (u *Unix) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var layout string
	if len(s) == 4 {
		layout = YearOnly
	} else if len(s) == 10 {
		layout = time.DateOnly
	} else {
		layout = time.DateTime
	}
	var t, err = time.Parse(layout, s)
	if err != nil {
		return err
	}
	*u = Unix(t.Unix())
	return nil
}

// MarshalYAML is YAML marshaler interface implementation to stream UNIX time.
func (u Unix) MarshalYAML() (any, error) {
	var t = u.Time()
	if t.Hour() != 0 || t.Minute() != 0 || t.Second() != 0 {
		return t.Format(time.DateTime), nil
	}
	if t.Month() != 1 || t.Day() != 1 {
		return t.Format(time.DateOnly), nil
	}
	return t.Format(YearOnly), nil
}

// UnmarshalYAML is YAML marshaler interface implementation to stream UNIX time.
func (u *Unix) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	var layout string
	if len(s) == 4 {
		layout = YearOnly
	} else if len(s) == 10 {
		layout = time.DateOnly
	} else {
		layout = time.DateTime
	}
	var t, err = time.Parse(layout, s)
	if err != nil {
		return err
	}
	*u = Unix(t.Unix())
	return nil
}

// MarshalXML is XML marshaler interface implementation to stream UNIX time.
func (u Unix) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var s string
	var t = u.Time()
	if t.Hour() != 0 || t.Minute() != 0 || t.Second() != 0 {
		s = t.Format(time.DateTime)
	} else if t.Month() != 1 || t.Day() != 1 {
		s = t.Format(time.DateOnly)
	} else {
		s = t.Format(YearOnly)
	}
	return e.EncodeElement(s, start)
}

// UnmarshalXML is XML marshaler interface implementation to stream UNIX time.
func (u *Unix) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	var layout string
	if len(s) == 4 {
		layout = YearOnly
	} else if len(s) == 10 {
		layout = time.DateOnly
	} else {
		layout = time.DateTime
	}
	var t, err = time.Parse(layout, s)
	if err != nil {
		return err
	}
	*u = Unix(t.Unix())
	return nil
}

func Year(year int) Unix {
	return Unix(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Unix())
}

func Date(year int, month time.Month, day int) Unix {
	return Unix(time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Unix())
}
