package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreAndRetrieveContent(t *testing.T) {
	// Prepare the test database
	testDBFile := "test.db"
	db, err := initDB(testDBFile)
	assert.NoError(t, err, "Error initializing database")
	defer os.Remove(testDBFile)

	// Test storeContent
	err = storeHistoryItem(db, "tag1,tag2", "Sample text 1")
	assert.NoError(t, err, "Error storing content")

	err = storeHistoryItem(db, "tag3,tag4", "Sample text 2")
	assert.NoError(t, err, "Error storing content")

	err = storeHistoryItem(db, "tag1,tag3", "Sample text 3")
	assert.NoError(t, err, "Error storing content")

	// Test retrieveContent
	last10, tagged, err := retrieveHistoryItems(db, "tag1,tag2")
	assert.NoError(t, err, "Error retrieving content")

	assert.Equal(t, 3, len(last10), "Incorrect number of last 10 entries")
	assert.Equal(t, 2, len(tagged), "Incorrect number of tagged entries")
}
