// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package durationpb_test

import (
	"encoding/json"
	"testing"
	"time"

	. "google.golang.org/protobuf/datatypes"

	"google.golang.org/protobuf/types/known/durationpb"
)

func TestDuration(t *testing.T) {
	type UserWithTime struct {
		ID   uint
		Name string
		Time *durationpb.Duration
	}

	DB.Migrator().DropTable(&UserWithTime{})
	if err := DB.Migrator().AutoMigrate(&UserWithTime{}); err != nil {
		t.Fatalf("failed to migrate, got error: %v", err)
	}

	user := UserWithTime{Name: "user1", Time: durationpb.NewDuration(1, 2, 3, 0)}
	DB.Create(&user)

	result := UserWithTime{}
	if err := DB.First(&result, "name = ? AND time = ?", "user1", durationpb.NewDuration(1, 2, 3, 0)).Error; err != nil {
		t.Fatalf("failed to find record with time, got error: %v", err)
	}

	AssertEqual(t, result.Time, durationpb.NewDuration(1, 2, 3, 0))
}

func TestJSONEncoding(t *testing.T) {
	date := durationpb.New(time.Hour)
	b, err := json.Marshal(date)
	if err != nil {
		t.Fatalf("failed to encode datatypes.Date: %v", err)
	}

	var got *durationpb.Duration
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("failed to decode to datatypes.Date: %v", err)
	}

	AssertEqual(t, date, got)
}
