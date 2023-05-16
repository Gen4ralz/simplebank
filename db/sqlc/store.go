package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provide all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	// buid a new Store object and return it
	return &Store{
		db: db,	// db is the input sql.DB
		Queries: New(db), // Queries is created by calling the New function with that db object
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID	int64	`json:"form_account_id"`
	ToAccountID		int64	`json:"to_account_id"`
	Amount			int64	`json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer	Transfer	`json:"transfer"`
	FromAccount	Account		`json:"from_account"`	// from_account after balance is updated
	ToAccount	Account		`json:"to_account"`		// to_account after balance is updated
	FromEntry	Entry		`json:"from_entry"`		// The entry of the from account which records that money is moving out
	ToEntry		Entry		`json:"to_entry"`		// The entry of the to account which records that money is moving in
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update account's balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	// create and run a new database transaction (transfer record)
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return nil
		}

		// add account entries
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
			// transaction will be rolled back
		}

		// add account entries
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
			// transaction will be rolled back
		}

		// TODO: update account's balance

		return nil
	})

	return result, err
}