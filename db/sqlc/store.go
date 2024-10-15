package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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

type txKeyType string

var txKey = txKeyType("txKey")


func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
    var result TransferTxResult

    // Start the transaction
    err := store.execTx(ctx, func(q *Queries) error {
        var err error

        // Type conversion from TransferTxParams to CreateTransferParams
        result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
        if err != nil {
            return fmt.Errorf("TransferTx - failed to create transfer: %w", err)
        }

        // Create entries for the FromAccount and ToAccount
        result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.FromAccountID,
            Amount:    -arg.Amount,
        })
        if err != nil {
            return fmt.Errorf("TransferTx - failed to create from entry: %w", err)
        }

        log.Printf("Created FromEntry: %+v", result.FromEntry)

        result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.ToAccountID,
            Amount:    arg.Amount,
        })
        if err != nil {
            return fmt.Errorf("TransferTx - failed to create to entry: %w", err)
        }

        log.Printf("Created ToEntry: %+v", result.ToEntry)

        // Update the account balances
        result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
        if err != nil {
            return fmt.Errorf("TransferTx - failed to update account balances: %w", err)
        }

        return nil
    })

    if err != nil {
        return result, err
    }

    return result, nil
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