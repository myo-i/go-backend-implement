package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	// このテストではaccount1からaccount2への送金の記録n回を考慮している
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before", account1.Balance, account2.Balance)

	// 同時性を注意深く意識しないといけないため、いくつかのgoroutineで処理をする
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
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
	existed := make(map[int]bool)
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

		// check from accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		// check to accounts
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// chek account's balance(残高)
		fmt.Println(">> tx", fromAccount.Balance, toAccount.Balance)
		balance1 := account1.Balance - fromAccount.Balance
		balance2 := toAccount.Balance - account2.Balance
		require.Equal(t, balance1, balance2)
		require.True(t, balance1 > 0)
		require.True(t, balance1%amount == 0) // amount, 2 * amount, 3 * amount, ..., n * amount（トランザクションの回数分、amountの量が増えるため）

		k := int(balance1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check updated balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, int64(n)*amount+account2.Balance, updateAccount2.Balance)

}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	// このテストではnが偶数回の時account1からaccount2へ
	// nが奇数回の時account2からaccount1への送金の記録をテストしている
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before", account1.Balance, account2.Balance)

	// account1からaccount2へ送金、その後account2からaccount1へ送金
	// 上記で１セットのため偶数回でセットする
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			// go funcでは条件が満たされていない場合、テスト全体が停止するという保証がないできない
			// そのため、処理の結果を一旦保持して後で確認する
			errs <- err
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check updated balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after", updateAccount1.Balance, updateAccount2.Balance)

	// 偶数回目でaccount1からaccount2へ、奇数回目でaccount2からaccount1へ送金しているので
	// 処理の前後で残高に変化がない
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)

}
