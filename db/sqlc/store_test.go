package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// 同時性を注意深く意識しないといけないため、いくつかのgoroutineで処理をする
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			// go funcでは条件が満たされていない場合、テスト全体が停止するという保証がないできない
			// そのため、処理の結果を一旦保持して後で確認する
			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check trancfer(db)
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// アカウントのIDからトランザクションのデータを取るわけではない
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check from entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// アカウントのIDからエントリのデータを取るわけではない
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check to entries
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// アカウントのIDからエントリのデータを取るわけではない
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO: chek account's balance

	}
}
