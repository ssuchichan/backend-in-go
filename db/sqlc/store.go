package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions.
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store.
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction.
func (store *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction.
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to another.
// It creates a transfer record, add account entries, and update account's balance within a single database transaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// 旧写法
		//account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		//if err != nil {
		//	return err
		//}
		//
		//result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		//	ID:      arg.FromAccountID,
		//	Balance: account1.Balance - arg.Amount,
		//})
		//if err != nil {
		//	return err
		//}
		//
		//account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		//if err != nil {
		//	return err
		//}
		//
		//result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		//	ID:      arg.ToAccountID,
		//	Balance: account2.Balance + arg.Amount,
		//})
		//if err != nil {
		//	return err
		//}

		// 现处理ID小的账户，保持处理顺序一致，防止死锁
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			//result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			//	Amount: -arg.Amount,
			//	ID:     arg.FromAccountID,
			//})
			//if err != nil {
			//	return err
			//}
			//result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			//	Amount: arg.Amount,
			//	ID:     arg.ToAccountID,
			//})
			//if err != nil {
			//	return err
			//}
		} else {
			//result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				Amount: arg.Amount,
				ID:     arg.ToAccountID,
			})
			if err != nil {
				return err
			}
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				Amount: -arg.Amount,
				ID:     arg.FromAccountID,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})

	return
}
