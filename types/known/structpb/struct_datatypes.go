// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package structpb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

var _ sql.Scanner = (*Struct)(nil)
var _ driver.Valuer = (*Struct)(nil)
var _ schema.GormDataTypeInterface = (*Struct)(nil)
var _ migrator.GormDataTypeInterface = (*Struct)(nil)
var _ gorm.Valuer = (*Struct)(nil)
var _ json.Marshaler = (*Struct)(nil)
var _ json.Unmarshaler = (*Struct)(nil)

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *Struct) Scan(val interface{}) error {
	if val == nil {
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", val))
	}
	return m.UnmarshalJSON(ba)
}

// Value return json value, implement driver.Valuer interface
func (st *Struct) Value() (driver.Value, error) {
	if st == nil {
		return nil, nil
	}
	ba, err := st.MarshalJSON()
	return string(ba), err
}

// GormDataType gorm common data type
func (st *Struct) GormDataType() string {
	return "struct"
}

// GormDBDataType gorm db data type
func (*Struct) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

func (st *Struct) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := st.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
