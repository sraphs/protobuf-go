// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package durationpb

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

var _ sql.Scanner = (*Duration)(nil)
var _ driver.Valuer = (*Duration)(nil)
var _ schema.GormDataTypeInterface = (*Duration)(nil)
var _ migrator.GormDataTypeInterface = (*Duration)(nil)
var _ json.Marshaler = (*Duration)(nil)
var _ json.Unmarshaler = (*Duration)(nil)

// NewDuration is a constructor for Time and returns new Time.
func NewDuration(hour, min, sec, nsec int) *Duration {
	return newDuration(hour, min, sec, nsec)
}

func newDuration(hour, min, sec, nsec int) *Duration {
	return New(time.Duration(
		time.Duration(hour)*time.Hour +
			time.Duration(min)*time.Minute +
			time.Duration(sec)*time.Second +
			time.Duration(nsec)*time.Nanosecond,
	))
}

// Scan implements sql.Scanner interface and scans value into Time,
func (t *Duration) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		t.setFromString(string(v))
	case string:
		t.setFromString(v)
	case time.Time:
		t.setFromTime(v)
	default:
		return errors.New(fmt.Sprintf("failed to scan value: %v", v))
	}

	return nil
}

func (t *Duration) setFromString(str string) {
	var h, m, s, n int
	fmt.Sscanf(str, "%02d:%02d:%02d.%09d", &h, &m, &s, &n)
	*t = *newDuration(h, m, s, n)
}

func (t *Duration) setFromTime(src time.Time) {
	*t = *newDuration(src.Hour(), src.Minute(), src.Second(), src.Nanosecond())
}

// Value implements driver.Valuer interface and returns string format of Time.
func (t *Duration) Value() (driver.Value, error) {
	return t.ToTimeString(), nil
}

// String implements fmt.Stringer interface.
func (t *Duration) ToTimeString() string {
	if nsec := t.nanoseconds(); nsec > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%09d", t.hours(), t.minutes(), t.seconds(), nsec)
	} else {
		// omit nanoseconds unless any value is specified
		return fmt.Sprintf("%02d:%02d:%02d", t.hours(), t.minutes(), t.seconds())
	}
}

func (t *Duration) hours() int {
	return int(time.Duration(t.AsDuration()).Truncate(time.Hour).Hours())
}

func (t *Duration) minutes() int {
	return int((time.Duration(t.AsDuration()) % time.Hour).Truncate(time.Minute).Minutes())
}

func (t *Duration) seconds() int {
	return int((time.Duration(t.AsDuration()) % time.Minute).Truncate(time.Second).Seconds())
}

func (t *Duration) nanoseconds() int {
	return int((time.Duration(t.AsDuration()) % time.Second).Nanoseconds())
}

// GormDataType returns gorm common data type. This type is used for the field's column type.
func (*Duration) GormDataType() string {
	return "time"
}

// GormDBDataType returns gorm DB data type based on the current using database.
func (*Duration) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "TIME"
	case "postgres":
		return "TIME"
	case "sqlserver":
		return "TIME"
	case "sqlite":
		return "TEXT"
	default:
		return ""
	}
}

// MarshalJSON implements json.Marshaler to convert *Duration to json serialization.
func (t *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.ToTimeString())
}

// UnmarshalJSON implements json.Unmarshaler to deserialize json data.
func (t *Duration) UnmarshalJSON(data []byte) error {
	// ignore null
	if string(data) == "null" {
		return nil
	}
	t.setFromString(strings.Trim(string(data), `"`))
	return nil
}
