package filter

import (
	"bytes"
	"testing"
)

func TestMarshalUnmarshal(t *testing.T) {
	// Test case 1: Filter with Ids only
	f1 := &F{
		Ids: [][]byte{
			[]byte("id1"),
			[]byte("id2"),
		},
		// These should be ignored when marshaling
		Authors: [][]byte{[]byte("author1")},
		Sort:    "asc",
	}

	data1, err := f1.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal filter with Ids: %v", err)
	}

	// Verify that only Ids are included
	if !bytes.Contains(data1, []byte("IDS:")) {
		t.Errorf("Marshaled data should contain IDS sentinel")
	}
	if bytes.Contains(data1, []byte("AUTHORS:")) {
		t.Errorf("Marshaled data should not contain AUTHORS sentinel when Ids are present")
	}
	if bytes.Contains(data1, []byte("SORT:")) {
		t.Errorf("Marshaled data should not contain SORT sentinel when Ids are present")
	}

	// Unmarshal back
	f1Unmarshaled := &F{}
	if err := f1Unmarshaled.Unmarshal(data1); err != nil {
		t.Fatalf("Failed to unmarshal filter with Ids: %v", err)
	}

	// Verify Ids are preserved
	if len(f1Unmarshaled.Ids) != 2 {
		t.Errorf("Expected 2 Ids, got %d", len(f1Unmarshaled.Ids))
	}
	// Verify Sort defaults to descending
	if f1Unmarshaled.Sort != "desc" {
		t.Errorf("Expected Sort to be 'desc', got '%s'", f1Unmarshaled.Sort)
	}

	// Test case 2: Filter with various fields
	f2 := &F{
		Authors:    [][]byte{[]byte("author1"), []byte("author2")},
		NotAuthors: [][]byte{[]byte("notauthor1")},
		Tags: TagMap{
			"tag1": [][]byte{[]byte("value1"), []byte("value2")},
			"tag2": [][]byte{[]byte("value3")},
		},
		NotTags: TagMap{
			"nottag1": [][]byte{[]byte("notvalue1")},
		},
		Since: 1000,
		Until: 2000,
		Sort:  "asc",
	}

	data2, err := f2.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal filter with various fields: %v", err)
	}

	// Verify all fields are included
	if !bytes.Contains(data2, []byte("AUTHORS:")) {
		t.Errorf("Marshaled data should contain AUTHORS sentinel")
	}
	if !bytes.Contains(data2, []byte("NOTAUTHORS:")) {
		t.Errorf("Marshaled data should contain NOTAUTHORS sentinel")
	}
	if !bytes.Contains(data2, []byte("TAGS:")) {
		t.Errorf("Marshaled data should contain TAGS sentinel")
	}
	if !bytes.Contains(data2, []byte("NOTTAGS:")) {
		t.Errorf("Marshaled data should contain NOTTAGS sentinel")
	}
	if !bytes.Contains(data2, []byte("SINCE:")) {
		t.Errorf("Marshaled data should contain SINCE sentinel")
	}
	if !bytes.Contains(data2, []byte("UNTIL:")) {
		t.Errorf("Marshaled data should contain UNTIL sentinel")
	}
	if !bytes.Contains(data2, []byte("SORT:")) {
		t.Errorf("Marshaled data should contain SORT sentinel when not 'desc'")
	}

	// Unmarshal back
	f2Unmarshaled := &F{}
	if err := f2Unmarshaled.Unmarshal(data2); err != nil {
		t.Fatalf("Failed to unmarshal filter with various fields: %v", err)
	}

	// Verify fields are preserved
	if len(f2Unmarshaled.Authors) != 2 {
		t.Errorf("Expected 2 Authors, got %d", len(f2Unmarshaled.Authors))
	}
	if len(f2Unmarshaled.NotAuthors) != 1 {
		t.Errorf("Expected 1 NotAuthor, got %d", len(f2Unmarshaled.NotAuthors))
	}
	if len(f2Unmarshaled.Tags) != 2 {
		t.Errorf("Expected 2 Tags, got %d", len(f2Unmarshaled.Tags))
	}
	if len(f2Unmarshaled.NotTags) != 1 {
		t.Errorf("Expected 1 NotTag, got %d", len(f2Unmarshaled.NotTags))
	}
	if f2Unmarshaled.Since != 1000 {
		t.Errorf("Expected Since to be 1000, got %d", f2Unmarshaled.Since)
	}
	if f2Unmarshaled.Until != 2000 {
		t.Errorf("Expected Until to be 2000, got %d", f2Unmarshaled.Until)
	}
	if f2Unmarshaled.Sort != "asc" {
		t.Errorf("Expected Sort to be 'asc', got '%s'", f2Unmarshaled.Sort)
	}

	// Test case 3: Filter with default Sort
	f3 := &F{
		Authors: [][]byte{[]byte("author1")},
		// Sort not specified, should default to "desc"
	}

	data3, err := f3.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal filter with default Sort: %v", err)
	}

	// Verify Sort is not included (defaults to descending)
	if bytes.Contains(data3, []byte("SORT:")) {
		t.Errorf("Marshaled data should not contain SORT sentinel when Sort is 'desc'")
	}

	// Unmarshal back
	f3Unmarshaled := &F{}
	if err := f3Unmarshaled.Unmarshal(data3); err != nil {
		t.Fatalf("Failed to unmarshal filter with default Sort: %v", err)
	}

	// Verify Sort defaults to descending
	if f3Unmarshaled.Sort != "desc" {
		t.Errorf("Expected Sort to be 'desc', got '%s'", f3Unmarshaled.Sort)
	}

	// Test case 4: Invalid combination - Ids with other fields
	invalidData := []byte("IDS:aWQx\nAUTHORS:YXV0aG9yMQ==")
	f4 := &F{}
	err = f4.Unmarshal(invalidData)
	if err == nil {
		t.Errorf("Expected error when unmarshaling Ids with other fields, got nil")
	}
}
