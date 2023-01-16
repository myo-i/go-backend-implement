package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// トランザクションを開始
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	// 処理を実行
	err = fn(q)
	if err != nil {
		// エラーが発生すればロールバックし、ロールバックでもエラーが発生すれば処理終了
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}

	// 保存
	return tx.Commit()
}

// TransferTxParams 一連の処理の譲渡に必要な入力のパラメータ
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult 一連の処理の譲渡の結果
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`     // 作成された転送レコード
	FromAccount Account  `json:"from_account"` // 残高が更新されたアカウント
	ToAccount   Account  `json:"to_account"`   // 上記アカウントが更新された後の説明となるアカウント
	FromEntry   Entry    `json:"from_entry"`   // お金が出て行ったアカウントのエントリ
	ToEntry     Entry    `json:"to_entry"`     // お金が入ったアカウントのエントリ
}

// store_testでキーが必要になるのでここで定義
var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add acount entries, and update accounts balnce within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	//
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		// Transfer作成
		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Entry作成
		fmt.Println(txName, "create fromEntry")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			// お金が出ていくのでマイナス
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create toEntry")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			// お金が入ってくるのでマイナス
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		// Account 作成
		// Update account's balance
		fmt.Println(txName, "get fromAccount")
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "get toAccount")
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "update toAccount")
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
