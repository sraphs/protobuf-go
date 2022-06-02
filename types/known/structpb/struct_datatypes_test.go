// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package structpb_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "google.golang.org/protobuf/datatypes"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/gorm"
)

func TestStruct(t *testing.T) {
	type UserWithStruct struct {
		gorm.Model
		Name       string
		Attributes *structpb.Struct
	}

	DB.Migrator().DropTable(&UserWithStruct{})
	if err := DB.Migrator().AutoMigrate(&UserWithStruct{}); err != nil {
		t.Errorf("failed to migrate, got error: %v", err)
	}

	// Go's json marshaler removes whitespace & orders keys alphabetically
	// use to compare against marshaled []byte of datatypes.JSON
	user1AttrsStr := `{"age":18,"name":"json-1","orgs":{"orga":"orga"},"tags":["tag1","tag2"]}`
	user1Attrs := map[string]interface{}{
		"age":  18,
		"name": "json-1",
		"orgs": map[string]interface{}{
			"orga": "orga",
		},
		"tags": []interface{}{"tag1", "tag2"},
	}
	userStruct, err := structpb.NewStruct(user1Attrs)
	if err != nil {
		t.Fatalf("failed to create struct, got error: %v", err)
	}

	user2Attrs := map[string]interface{}{
		"name": "json-2",
		"age":  28,
		"tags": []interface{}{"tag1", "tag3"},
		"role": "admin",
		"orgs": map[string]interface{}{
			"orgb": "orgb",
		},
	}
	user2Struct, err := structpb.NewStruct(user2Attrs)
	if err != nil {
		t.Fatalf("failed to create struct, got error: %v", err)
	}

	users := []UserWithStruct{{
		Name:       "json-1",
		Attributes: userStruct,
	}, {
		Name:       "json-2",
		Attributes: user2Struct,
	},
		{
			Name:       "json-3",
			Attributes: &structpb.Struct{},
		},
	}

	if err := DB.Create(&users).Error; err != nil {
		t.Errorf("Failed to create users %v", err)
	}

	var result UserWithStruct
	if err := DB.First(&result, JSONQuery("attributes").HasKey("role")).Error; err != nil {
		t.Fatalf("failed to find user with json key, got error %v", err)
	}
	AssertEqual(t, result.Name, users[1].Name)

	var result2 UserWithStruct
	if err := DB.First(&result2, JSONQuery("attributes").HasKey("orgs", "orga")).Error; err != nil {
		t.Fatalf("failed to find user with json key, got error %v", err)
	}
	AssertEqual(t, result2.Name, users[0].Name)

	AssertEqual(t, result2.Attributes.AsMap(), user1Attrs)

	// attributes should not marshal to base64 encoded []byte
	result2Attrs, err := json.Marshal(result2.Attributes)
	if err != nil {
		t.Fatalf("failed to marshal result2.Attributes, got error %v", err)
	}

	AssertEqual(t, string(result2Attrs), user1AttrsStr)

	// []byte should unmarshal into type structpb.Struct
	var j structpb.Struct
	if err := json.Unmarshal([]byte(user1AttrsStr), &j); err != nil {
		t.Fatalf("failed to unmarshal user1Attrs, got error %v", err)
	}

	AssertEqual(t, fmt.Sprint(j.AsMap()), fmt.Sprint(user1Attrs))

	var result3 UserWithStruct
	if err := DB.First(&result3, JSONQuery("attributes").Equals("json-1", "name")).Error; err != nil {
		t.Fatalf("failed to find user with json value, got error %v", err)
	}
	AssertEqual(t, result3.Name, users[0].Name)

	var result4 UserWithStruct
	if err := DB.First(&result4, JSONQuery("attributes").Equals("orgb", "orgs", "orgb")).Error; err != nil {
		t.Fatalf("failed to find user with json value, got error %v", err)
	}
	AssertEqual(t, result4.Name, users[1].Name)

	// FirstOrCreate
	jsonMap := map[string]interface{}{"Attributes": JSON(`{"age":19,"name":"json-1","orgs":{"orga":"orga"},"tags":["tag1","tag2"]}`)}
	if err := DB.Where(&UserWithStruct{Name: "json-1"}).Assign(jsonMap).FirstOrCreate(&UserWithStruct{}).Error; err != nil {
		t.Errorf("failed to run FirstOrCreate")
	}

	var result5 UserWithStruct
	if err := DB.First(&result5, JSONQuery("attributes").Equals(19, "age")).Error; err != nil {
		t.Fatalf("failed to find user with json value, got error %v", err)
	}

	var result6 UserWithStruct
	if err := DB.Where("name = ?", "json-3").First(&result6).Error; err != nil {
		t.Fatalf("failed to find user with json value, got error %v", err)
	}

	AssertEqual(t, result6.Attributes, &structpb.Struct{})

	type UserWithStructPtr struct {
		gorm.Model
		Name       string
		Attributes *structpb.Struct
	}

	DB.Migrator().DropTable(&UserWithStructPtr{})
	if err := DB.Migrator().AutoMigrate(&UserWithStructPtr{}); err != nil {
		t.Errorf("failed to migrate, got error: %v", err)
	}

	jm1, err := structpb.NewStruct(user1Attrs)
	if err != nil {
		t.Fatalf("failed to create struct, got error: %v", err)
	}

	ujmps := []*UserWithStructPtr{
		{
			Name:       "json-4",
			Attributes: jm1,
		},
		{
			Name: "json-5",
		},
	}

	if err := DB.Create(&ujmps).Error; err != nil {
		t.Errorf("Failed to create users %v", err)
	}

	var result7 UserWithStructPtr
	if err := DB.Where("name = ?", "json-4").First(&result7).Error; err != nil {
		t.Fatalf("failed to find user with json value, got error %v", err)
	}

	AssertEqual(t, result7.Attributes, jm1)

	var result8 UserWithStructPtr
	if err := DB.Where("name = ?", "json-5").First(&result8).Error; err != nil {
		t.Fatalf("failed to find user with json value, got error %v", err)
	}

	AssertEqual(t, result8.Attributes, nil)

	var result9 UserWithStructPtr
	if err := DB.Where(result8, "Attributes").First(&result9).Error; err != nil {
		t.Fatalf("failed to find user with json value, got error %v", err)
	}
}
