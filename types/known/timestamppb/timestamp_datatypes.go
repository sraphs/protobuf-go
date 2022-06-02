// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package timestamppb

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
)

var _ sql.Scanner = (*Timestamp)(nil)
var _ driver.Valuer = (*Timestamp)(nil)
var _ json.Marshaler = (*Timestamp)(nil)
var _ json.Unmarshaler = (*Timestamp)(nil)

func (date *Timestamp) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	if err != nil || !nullTime.Valid {
		return nil
	}
	*date = *New(nullTime.Time)
	return
}

func (date *Timestamp) Value() (driver.Value, error) {
	if date == nil || !date.IsValid() {
		return nil, nil
	}

	return date.AsTime(), nil
}

func (date *Timestamp) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(date)
}

func (date *Timestamp) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, date)
}
