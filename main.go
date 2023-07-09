package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sync"

	"github.com/arantesxyz/crdbledgertest/accounts"
	"github.com/arantesxyz/crdbledgertest/database"
	"github.com/arantesxyz/crdbledgertest/transfers"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getAccountIds() []string {
	return []string{"50fa342e-5fb7-48e2-8c9b-a34157117b45", "8639b936-86b0-4109-bd42-f2848eb497d0"}
}

func generateRandomTransfer() transfers.Transfer {
	randAmount := float32(rand.Intn(200) - 80)
	accounts := getAccountIds()

	accId := accounts[(int(math.Abs(float64(randAmount))) % len(accounts))]

	return transfers.Transfer{
		AccountId: accId,
		Amount:    randAmount,
		IsCredit:  randAmount >= 0,
		Status:    "PROCESSING",
	}
}

func processTransfer(pool *pgxpool.Pool, transfer *transfers.Transfer) error {
	crdbpgx.ExecuteTx(context.Background(), pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return transfers.CreateTransfer(tx, transfer)
	})
	fmt.Println("Transfer created!", transfer)

	err := crdbpgx.ExecuteTx(context.Background(), pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		err := accounts.UpdateBalance(tx, transfer.AccountId, transfer.Amount)
		if err != nil {
			fmt.Println("err update balance", err)
			return errors.New("balanceerror")
		}

		err = transfers.UpdateTransfer(tx, transfer, "APPROVED")
		if err != nil {
			fmt.Println("err updating transfer", err)
			return err
		}

		return nil
	})

	if err != nil {
		if err.Error() == "balanceerror" {
			fmt.Println("Error trying to update balance, cancelling the transaction")
			return crdbpgx.ExecuteTx(context.Background(), pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
				return transfers.UpdateTransfer(tx, transfer, "CANCELLED")
			})
		}
		fmt.Println("processing error", err)
	}

	fmt.Println("Balance updated")

	return nil
}

func exec(pool *pgxpool.Pool, amount int) {
	for i := 0; i < amount; i++ {
		transfer := generateRandomTransfer()

		processTransfer(pool, &transfer)
		fmt.Println("Processed!", transfer)
	}
}

func main() {
	pool := database.CreatePool(os.Getenv("DATABASE_URI"))
	defer pool.Close()

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exec(pool, 10)
		}()
	}

	wg.Wait()
}
