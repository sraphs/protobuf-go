package timestamppb_test

import (
	"encoding/json"
	"testing"
	time "time"

	"github.com/jinzhu/now"
	. "google.golang.org/protobuf/datatypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDate(t *testing.T) {
	type UserWithDate struct {
		ID   uint
		Name string
		Date *timestamppb.Timestamp
	}

	DB.Migrator().DropTable(&UserWithDate{})
	if err := DB.Migrator().AutoMigrate(&UserWithDate{}); err != nil {
		t.Errorf("failed to migrate, got error: %v", err)
	}

	curTime := time.Now().UTC()
	beginningOfDay := now.New(curTime).BeginningOfDay()

	user := UserWithDate{Name: "jinzhu", Date: timestamppb.New(curTime)}
	DB.Create(&user)

	result := UserWithDate{}
	if err := DB.First(&result, "name = ? AND date = ?", "jinzhu", timestamppb.New(curTime)).Error; err != nil {
		t.Fatalf("Failed to find record with date")
	}

	AssertEqual(t, result.Date, beginningOfDay)
}

func TestJSONEncoding(t *testing.T) {
	date := timestamppb.New(time.Now())
	b, err := json.Marshal(date)
	if err != nil {
		t.Fatalf("failed to encode datatypes.Date: %v", err)
	}

	var got *timestamppb.Timestamp
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("failed to decode to datatypes.Date: %v", err)
	}

	AssertEqual(t, date, got)
}
