package db

import (
	"context"
	"database/sql"
	"goprojects/simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	// Create a random account for the entry
	account1 := createRandomAccount(t)

	// Prepare the entry creation parameters
	arg := CreateEntryParams{
		AccountID: account1.ID,
		Amount:    util.RandomMoney(),
	}

	// Create the entry
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	// Check the entry values match expectations
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	// Now actually call createRandomEntry to test the CreateEntry functionality
	entry := createRandomEntry(t)
	require.NotEmpty(t, entry)
}

func TestGetEntry(t *testing.T) {
	// Create a random entry
	entry1 := createRandomEntry(t)

	// Fetch the entry by ID
	entry2, err := testQueries.GetAEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	// Check that the fetched entry matches the original entry
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	// Create a random entry
	entry1 := createRandomEntry(t)

	// Prepare update parameters
	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomMoney(),
	}

	// Update the entry
	entry2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	// Check that the updated entry matches expectations
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, arg.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	// Create a random entry
	entry1 := createRandomEntry(t)

	// Delete the entry
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err) // Check for errors during deletion

	// Try to fetch the entry after deletion, should return sql.ErrNoRows
	entry2, err := testQueries.GetAEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestListEntries(t *testing.T) {
	// Create 10 random entries
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	// Prepare list parameters
	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	// List entries
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	// Check that each entry is not empty
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
