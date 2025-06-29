package database

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"manifold.mleku.dev/chk"
	"manifold.mleku.dev/event"
	"manifold.mleku.dev/filter"
	"manifold.mleku.dev/p256k"
)

// TestQueryEvents tests the QueryEvents function with various filter combinations.
func TestQueryEvents(t *testing.T) {
	// Create a temporary directory for the database
	tempDir, err := os.MkdirTemp("", "manifold-test-db")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new database
	db := New()
	if err = db.Init(tempDir); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Generate test events
	var events []*event.E
	events, err = generateTestEvents(10)
	if err != nil {
		t.Fatalf("Failed to generate test events: %v", err)
	}

	// Store events in the database
	for _, ev := range events {
		if err = db.StoreEvent(ev); err != nil {
			t.Fatalf("Failed to store event: %v", err)
		}
	}

	// Test cases
	t.Run("FilterByIds", func(t *testing.T) {
		testFilterByIds(t, db, events)
	})

	t.Run("FilterByAuthors", func(t *testing.T) {
		testFilterByAuthors(t, db, events)
	})

	t.Run("FilterByTags", func(t *testing.T) {
		testFilterByTags(t, db, events)
	})

	t.Run("FilterByAuthorsAndTags", func(t *testing.T) {
		testFilterByAuthorsAndTags(t, db, events)
	})

	t.Run("FilterByTimestamp", func(t *testing.T) {
		testFilterByTimestamp(t, db, events)
	})

	t.Run("FilterByAuthorAndTimestamp", func(t *testing.T) {
		testFilterByAuthorAndTimestamp(t, db, events)
	})

	t.Run("FilterByTagsAndTimestamp", func(t *testing.T) {
		testFilterByTagsAndTimestamp(t, db, events)
	})

	t.Run("FilterByAuthorTagsAndTimestamp", func(t *testing.T) {
		testFilterByAuthorTagsAndTimestamp(t, db, events)
	})

	t.Run("FilterWithNoSpecificCriteria", func(t *testing.T) {
		testFilterWithNoSpecificCriteria(t, db, events)
	})

	t.Run("SortingAscending", func(t *testing.T) {
		testSortingAscending(t, db, events)
	})

	t.Run("SortingDescending", func(t *testing.T) {
		testSortingDescending(t, db, events)
	})

	// Tests for negation features
	t.Run("FilterByNotIds", func(t *testing.T) {
		testFilterByNotIds(t, db, events)
	})

	t.Run("FilterByNotAuthors", func(t *testing.T) {
		testFilterByNotAuthors(t, db, events)
	})

	t.Run("FilterByNotTags", func(t *testing.T) {
		testFilterByNotTags(t, db, events)
	})

	t.Run("FilterByCombinedNegations", func(t *testing.T) {
		testFilterByCombinedNegations(t, db, events)
	})
}

// generateTestEvents generates a set of test events with various properties.
func generateTestEvents(count int) ([]*event.E, error) {
	events := make([]*event.E, count)

	// Create signers for different authors
	signers := make([]*p256k.Signer, 3)
	for i := range signers {
		signers[i] = new(p256k.Signer)
		if err := signers[i].Generate(); chk.E(err) {
			return nil, err
		}
	}

	// Create events with different properties
	for i := 0; i < count; i++ {
		// Use different authors
		authorIndex := i % len(signers)
		signer := signers[authorIndex]

		// Create event
		ev := &event.E{
			Pubkey:    signer.Pub(),
			Timestamp: time.Now().Add(time.Duration(i) * time.Minute).Unix(), // Different timestamps
			Content:   []byte(fmt.Sprintf("Test content %d", i)),
			Tags:      &event.Tags{},
		}

		// Add tags based on index
		if i%2 == 0 {
			*ev.Tags = append(*ev.Tags, event.Tag{
				Key:   []byte("type"),
				Value: []byte("text"),
			})
		}
		if i%3 == 0 {
			*ev.Tags = append(*ev.Tags, event.Tag{
				Key:   []byte("category"),
				Value: []byte("test"),
			})
		}
		if i%4 == 0 {
			*ev.Tags = append(*ev.Tags, event.Tag{
				Key:   []byte("importance"),
				Value: []byte("high"),
			})
		}

		// Sign the event
		if err := ev.Sign(signer); chk.E(err) {
			return nil, err
		}

		events[i] = ev
	}

	return events, nil
}

// testFilterByIds tests filtering events by IDs.
func testFilterByIds(t *testing.T, db *D, events []*event.E) {
	// Get IDs of the first 3 events
	ids := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		var err error
		ids[i], err = events[i].Id()
		if err != nil {
			t.Fatalf("Failed to get event ID: %v", err)
		}
	}

	// Create filter with IDs
	f := filter.F{
		Ids: ids,
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != len(ids) {
		t.Fatalf("Expected %d events, got %d", len(ids), len(result))
	}

	// Check that all IDs are in the result
	for _, id := range ids {
		found := false
		for _, resultId := range result {
			if bytes.Equal(id, resultId) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("ID %x not found in result", id)
		}
	}
}

// testFilterByAuthors tests filtering events by authors.
func testFilterByAuthors(t *testing.T, db *D, events []*event.E) {
	// Get the pubkey of the first author
	pubkey := events[0].Pubkey

	// Count events with this pubkey
	expectedCount := 0
	for _, ev := range events {
		if bytes.Equal(ev.Pubkey, pubkey) {
			expectedCount++
		}
	}

	// Create filter with author
	f := filter.F{
		Authors: [][]byte{pubkey},
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify each result has the correct author
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}
		if !bytes.Equal(ev.Pubkey, pubkey) {
			t.Fatalf("Event has incorrect author")
		}
	}
}

// testFilterByTags tests filtering events by tags.
func testFilterByTags(t *testing.T, db *D, events []*event.E) {
	// Create filter with tag
	f := filter.F{
		Tags: filter.TagMap{
			"type": {[]byte("text")},
		},
	}

	// Count events with this tag
	expectedCount := 0
	for _, ev := range events {
		for _, tag := range *ev.Tags {
			if bytes.Equal(tag.Key, []byte("type")) && bytes.Equal(tag.Value, []byte("text")) {
				expectedCount++
				break
			}
		}
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify each result has the correct tag
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		found := false
		for _, tag := range *ev.Tags {
			if bytes.Equal(tag.Key, []byte("type")) && bytes.Equal(tag.Value, []byte("text")) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Event does not have the expected tag")
		}
	}
}

// testFilterByAuthorsAndTags tests filtering events by both authors and tags.
func testFilterByAuthorsAndTags(t *testing.T, db *D, events []*event.E) {
	// Get the pubkey of the first author
	pubkey := events[0].Pubkey

	// Create filter with author and tag
	f := filter.F{
		Authors: [][]byte{pubkey},
		Tags: filter.TagMap{
			"type": {[]byte("text")},
		},
	}

	// Count events with this author and tag
	expectedCount := 0
	for _, ev := range events {
		if !bytes.Equal(ev.Pubkey, pubkey) {
			continue
		}

		for _, tag := range *ev.Tags {
			if bytes.Equal(tag.Key, []byte("type")) && bytes.Equal(tag.Value, []byte("text")) {
				expectedCount++
				break
			}
		}
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify each result has the correct author and tag
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		if !bytes.Equal(ev.Pubkey, pubkey) {
			t.Fatalf("Event has incorrect author")
		}

		found := false
		for _, tag := range *ev.Tags {
			if bytes.Equal(tag.Key, []byte("type")) && bytes.Equal(tag.Value, []byte("text")) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Event does not have the expected tag")
		}
	}
}

// testFilterByTimestamp tests filtering events by timestamp.
func testFilterByTimestamp(t *testing.T, db *D, events []*event.E) {
	// Find the middle timestamp
	timestamps := make([]int64, len(events))
	for i, ev := range events {
		timestamps[i] = ev.Timestamp
	}

	// Sort timestamps
	for i := 0; i < len(timestamps); i++ {
		for j := i + 1; j < len(timestamps); j++ {
			if timestamps[i] > timestamps[j] {
				timestamps[i], timestamps[j] = timestamps[j], timestamps[i]
			}
		}
	}

	// Use the middle timestamp as the since value
	middleIndex := len(timestamps) / 2
	since := timestamps[middleIndex]

	// Create filter with since timestamp
	f := filter.F{
		Since: since,
	}

	// Count events with timestamp >= since
	expectedCount := 0
	for _, ev := range events {
		if ev.Timestamp >= since {
			expectedCount++
		}
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify each result has the correct timestamp
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		if ev.Timestamp < since {
			t.Fatalf("Event has timestamp %d, which is less than since %d", ev.Timestamp, since)
		}
	}
}

// testFilterByAuthorAndTimestamp tests filtering events by both author and timestamp.
func testFilterByAuthorAndTimestamp(t *testing.T, db *D, events []*event.E) {
	// Get the pubkey of the first author
	pubkey := events[0].Pubkey

	// Find the middle timestamp
	timestamps := make([]int64, len(events))
	for i, ev := range events {
		timestamps[i] = ev.Timestamp
	}

	// Sort timestamps
	for i := 0; i < len(timestamps); i++ {
		for j := i + 1; j < len(timestamps); j++ {
			if timestamps[i] > timestamps[j] {
				timestamps[i], timestamps[j] = timestamps[j], timestamps[i]
			}
		}
	}

	// Use the middle timestamp as the since value
	middleIndex := len(timestamps) / 2
	since := timestamps[middleIndex]

	// Create filter with author and since timestamp
	f := filter.F{
		Authors: [][]byte{pubkey},
		Since:   since,
	}

	// Count events with this author and timestamp >= since
	expectedCount := 0
	for _, ev := range events {
		if bytes.Equal(ev.Pubkey, pubkey) && ev.Timestamp >= since {
			expectedCount++
		}
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify each result has the correct author and timestamp
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		if !bytes.Equal(ev.Pubkey, pubkey) {
			t.Fatalf("Event has incorrect author")
		}

		if ev.Timestamp < since {
			t.Fatalf("Event has timestamp %d, which is less than since %d", ev.Timestamp, since)
		}
	}
}

// testFilterByTagsAndTimestamp tests filtering events by both tags and timestamp.
func testFilterByTagsAndTimestamp(t *testing.T, db *D, events []*event.E) {
	// Create filter with tag
	tagKey := "type"
	tagValue := "text"

	// Find the middle timestamp
	timestamps := make([]int64, len(events))
	for i, ev := range events {
		timestamps[i] = ev.Timestamp
	}

	// Sort timestamps
	for i := 0; i < len(timestamps); i++ {
		for j := i + 1; j < len(timestamps); j++ {
			if timestamps[i] > timestamps[j] {
				timestamps[i], timestamps[j] = timestamps[j], timestamps[i]
			}
		}
	}

	// Use the middle timestamp as the since value
	middleIndex := len(timestamps) / 2
	since := timestamps[middleIndex]

	// Create filter with tag and since timestamp
	f := filter.F{
		Tags: filter.TagMap{
			tagKey: {[]byte(tagValue)},
		},
		Since: since,
	}

	// Count events with this tag and timestamp >= since
	expectedCount := 0
	for _, ev := range events {
		if ev.Timestamp < since {
			continue
		}

		for _, tag := range *ev.Tags {
			if bytes.Equal(tag.Key, []byte(tagKey)) && bytes.Equal(tag.Value, []byte(tagValue)) {
				expectedCount++
				break
			}
		}
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify each result has the correct tag and timestamp
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		if ev.Timestamp < since {
			t.Fatalf("Event has timestamp %d, which is less than since %d", ev.Timestamp, since)
		}

		found := false
		for _, tag := range *ev.Tags {
			if bytes.Equal(tag.Key, []byte(tagKey)) && bytes.Equal(tag.Value, []byte(tagValue)) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Event does not have the expected tag")
		}
	}
}

// testFilterByAuthorTagsAndTimestamp tests filtering events by author, tags, and timestamp.
func testFilterByAuthorTagsAndTimestamp(t *testing.T, db *D, events []*event.E) {
	// Get the pubkey of the first author
	pubkey := events[0].Pubkey

	// Create filter with tag
	tagKey := "type"
	tagValue := "text"

	// Find the middle timestamp
	timestamps := make([]int64, len(events))
	for i, ev := range events {
		timestamps[i] = ev.Timestamp
	}

	// Sort timestamps
	for i := 0; i < len(timestamps); i++ {
		for j := i + 1; j < len(timestamps); j++ {
			if timestamps[i] > timestamps[j] {
				timestamps[i], timestamps[j] = timestamps[j], timestamps[i]
			}
		}
	}

	// Use the middle timestamp as the since value
	middleIndex := len(timestamps) / 2
	since := timestamps[middleIndex]

	// Create filter with author, tag, and since timestamp
	f := filter.F{
		Authors: [][]byte{pubkey},
		Tags: filter.TagMap{
			tagKey: {[]byte(tagValue)},
		},
		Since: since,
	}

	// Count events with this author, tag, and timestamp >= since
	expectedCount := 0
	for _, ev := range events {
		if !bytes.Equal(ev.Pubkey, pubkey) || ev.Timestamp < since {
			continue
		}

		for _, tag := range *ev.Tags {
			if bytes.Equal(tag.Key, []byte(tagKey)) && bytes.Equal(tag.Value, []byte(tagValue)) {
				expectedCount++
				break
			}
		}
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify each result has the correct author, tag, and timestamp
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		if !bytes.Equal(ev.Pubkey, pubkey) {
			t.Fatalf("Event has incorrect author")
		}

		if ev.Timestamp < since {
			t.Fatalf("Event has timestamp %d, which is less than since %d", ev.Timestamp, since)
		}

		found := false
		for _, tag := range *ev.Tags {
			if bytes.Equal(tag.Key, []byte(tagKey)) && bytes.Equal(tag.Value, []byte(tagValue)) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Event does not have the expected tag")
		}
	}
}

// testFilterWithNoSpecificCriteria tests filtering events with no specific criteria.
func testFilterWithNoSpecificCriteria(t *testing.T, db *D, events []*event.E) {
	// Create empty filter
	f := filter.F{}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != len(events) {
		t.Fatalf("Expected %d events, got %d", len(events), len(result))
	}
}

// testSortingAscending tests sorting events in ascending order.
func testSortingAscending(t *testing.T, db *D, events []*event.E) {
	// Create filter with ascending sort
	f := filter.F{
		Sort: "asc",
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != len(events) {
		t.Fatalf("Expected %d events, got %d", len(events), len(result))
	}

	// Verify sorting
	for i := 1; i < len(result); i++ {
		ev1, err := db.GetEventById(result[i-1])
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		ev2, err := db.GetEventById(result[i])
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}
		if ev1.Timestamp > ev2.Timestamp {
			t.Fatalf("Events not sorted in ascending order")
		}
	}
}

// testSortingDescending tests sorting events in descending order.
func testSortingDescending(t *testing.T, db *D, events []*event.E) {
	// Create filter with descending sort
	f := filter.F{
		Sort: "desc",
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != len(events) {
		t.Fatalf("Expected %d events, got %d", len(events), len(result))
	}

	// Verify sorting
	for i := 1; i < len(result); i++ {
		ev1, err := db.GetEventById(result[i-1])
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		ev2, err := db.GetEventById(result[i])
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		if ev1.Timestamp < ev2.Timestamp {
			t.Fatalf("Events not sorted in descending order")
		}
	}
}

// testFilterByNotIds tests filtering events by excluding specific IDs.
func testFilterByNotIds(t *testing.T, db *D, events []*event.E) {
	// Get all event IDs
	allIds := make([][]byte, len(events))
	for i, ev := range events {
		var err error
		allIds[i], err = ev.Id()
		if err != nil {
			t.Fatalf("Failed to get event ID: %v", err)
		}
	}

	// Exclude the first 2 events
	notIds := allIds[:2]

	// Create filter with all IDs and NotIds
	f := filter.F{
		Ids:    allIds,
		NotIds: notIds,
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	expectedCount := len(allIds) - len(notIds)
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify that excluded IDs are not in the result
	for _, notId := range notIds {
		for _, resultId := range result {
			if bytes.Equal(notId, resultId) {
				t.Fatalf("Excluded ID %x found in result", notId)
			}
		}
	}
}

// testFilterByNotAuthors tests filtering events by excluding specific authors.
func testFilterByNotAuthors(t *testing.T, db *D, events []*event.E) {
	// Get the pubkey of the first author
	notAuthor := events[0].Pubkey

	// Count events with this pubkey
	excludedCount := 0
	for _, ev := range events {
		if bytes.Equal(ev.Pubkey, notAuthor) {
			excludedCount++
		}
	}

	// Create filter with NotAuthors
	f := filter.F{
		NotAuthors: [][]byte{notAuthor},
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	expectedCount := len(events) - excludedCount
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify that no events from the excluded author are in the result
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}
		if bytes.Equal(ev.Pubkey, notAuthor) {
			t.Fatalf("Event from excluded author found in result")
		}
	}
}

// testFilterByNotTags tests filtering events by excluding specific tags.
func testFilterByNotTags(t *testing.T, db *D, events []*event.E) {
	// Create filter with NotTags
	notTagKey := "type"
	notTagValue := "text"

	f := filter.F{
		NotTags: filter.TagMap{
			notTagKey: {[]byte(notTagValue)},
		},
	}

	// Count events with this tag
	excludedCount := 0
	for _, ev := range events {
		if ev.Tags != nil {
			for _, tag := range *ev.Tags {
				if bytes.Equal(tag.Key, []byte(notTagKey)) && bytes.Equal(tag.Value, []byte(notTagValue)) {
					excludedCount++
					break
				}
			}
		}
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	expectedCount := len(events) - excludedCount
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify that no events with the excluded tag are in the result
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		if ev.Tags != nil {
			for _, tag := range *ev.Tags {
				if bytes.Equal(tag.Key, []byte(notTagKey)) && bytes.Equal(tag.Value, []byte(notTagValue)) {
					t.Fatalf("Event with excluded tag found in result")
				}
			}
		}
	}
}

// testFilterByCombinedNegations tests filtering events by combining multiple negation criteria.
func testFilterByCombinedNegations(t *testing.T, db *D, events []*event.E) {
	// Get the pubkey of the first author
	notAuthor := events[0].Pubkey

	// Create filter with NotAuthors and NotTags
	notTagKey := "type"
	notTagValue := "text"

	f := filter.F{
		NotAuthors: [][]byte{notAuthor},
		NotTags: filter.TagMap{
			notTagKey: {[]byte(notTagValue)},
		},
	}

	// Count events that should be excluded
	excludedCount := 0
	for _, ev := range events {
		excluded := false

		// Check if author is excluded
		if bytes.Equal(ev.Pubkey, notAuthor) {
			excluded = true
		}

		// Check if tag is excluded
		if !excluded && ev.Tags != nil {
			for _, tag := range *ev.Tags {
				if bytes.Equal(tag.Key, []byte(notTagKey)) && bytes.Equal(tag.Value, []byte(notTagValue)) {
					excluded = true
					break
				}
			}
		}

		if excluded {
			excludedCount++
		}
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	expectedCount := len(events) - excludedCount
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify that no events with the excluded criteria are in the result
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		// Check author
		if bytes.Equal(ev.Pubkey, notAuthor) {
			t.Fatalf("Event from excluded author found in result")
		}

		// Check tags
		if ev.Tags != nil {
			for _, tag := range *ev.Tags {
				if bytes.Equal(tag.Key, []byte(notTagKey)) && bytes.Equal(tag.Value, []byte(notTagValue)) {
					t.Fatalf("Event with excluded tag found in result")
				}
			}
		}
	}
}
