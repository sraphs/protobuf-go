package durationpb_test

import (
	"testing"

	. "google.golang.org/protobuf/datatypes"
	. "gorm.io/gorm/utils/tests"

	"google.golang.org/protobuf/types/known/durationpb"
)

func TestDuration(t *testing.T) {
	if SupportedDriver("mysql", "postgres", "sqlite", "sqlserver") {
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
}
