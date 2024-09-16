package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	//a query only performs a single operation in a table thus needs to extend the struct functionality in golang via composition
	*Queries
	//creating a new transaction
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

//function to execite a geeneric database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error{
	//create a new db transaction
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	//call New() function with the created transaction to get back a new queries object
	//The New() is provided by sqlc
	q := New(tx)
	//we have the queries that runs within the transaction so we can call the input function with that queries and get back an error
	err = fn(q)
	//if the error is not nil, rollback the transaction
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
//the struct contains all the input parameters needed to transfer money between two accounts
type TransferTxParams struct {
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        int64     `json:"amount"`
}

//the struct contains the resultof the transfer transaction
type TransferTxResult struct {
	Transfer     Transfer     `json:"transfer"`
	FromAccount  Account      `json:"from_account"`
	ToAccount    Account      `json:"to_account"`
	FromEntry    Entry        `json:"from_entry"`
	ToEntry      Entry        `json:"to_entry"`
}

//function to perform money transfer transacton
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	//create an empty result
	var result TransferTxResult
	// call the store store.ExecTx function to create and run the database transaction
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
			Amount: -arg.Amount,			
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

		return nil
	})

	return result, err
}