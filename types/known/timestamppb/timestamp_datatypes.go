// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package timestamppb

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	time "time"

	"google.golang.org/protobuf/encoding/protojson"
)

var _ sql.Scanner = (*Timestamp)(nil)
var _ driver.Valuer = (*Timestamp)(nil)
var _ json.Marshaler = (*Timestamp)(nil)
var _ json.Unmarshaler = (*Timestamp)(nil)

func (date *Timestamp) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = *New(nullTime.Time)
	return
}

func (date *Timestamp) Value() (driver.Value, error) {
	y, m, d := time.Time(date.AsTime()).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date.AsTime()).Location()), nil
}

func (date *Timestamp) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(date)
}

func (date *Timestamp) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, date)
}
