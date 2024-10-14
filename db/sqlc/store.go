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

var txKey = struct{}{}

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
		if arg.FromAccountID > arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		// Check if there was an error during addMoney.
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	if err != nil {
		return
	}

	return account1, account2, nil
}