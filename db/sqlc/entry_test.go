package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/techschool/simple_bank/util"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: sql.NullInt64{Int64: account.ID, Valid: true},
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.ID)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t, createRandomAccount(t))
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t, createRandomAccount(t))

	entry1, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry1)

	require.Equal(t, entry.ID, entry1.ID)
	require.Equal(t, entry.AccountID, entry1.AccountID)
	require.Equal(t, entry.Amount, entry1.Amount)
	require.WithinDuration(t, entry.CreatedAt, entry1.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	entry := createRandomEntry(t, createRandomAccount(t))
	arg := UpdateEntryParams{
		ID:     entry.ID,
		Amount: util.RandomMoney(),
	}

	entry1, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry1)

	require.Equal(t, arg.ID, entry1.ID)
	require.Equal(t, entry.AccountID, entry1.AccountID)
	require.Equal(t, arg.Amount, entry1.Amount)
	require.WithinDuration(t, entry.CreatedAt, entry1.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t, createRandomAccount(t))

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	entry1, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry1)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t, createRandomAccount(t))
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
