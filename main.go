package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/arantesxyz/crdbledgertest/accounts"
	"github.com/arantesxyz/crdbledgertest/database"
	"github.com/arantesxyz/crdbledgertest/transfers"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
)

const testAccountId = "50fa342e-5fb7-48e2-8c9b-a34157117b45"

func generateRandomTransfer() transfers.Transfer {
	randAmount := float32(rand.Intn(200) - 100)
	return transfers.Transfer{
		AccountId: testAccountId,
		Amount:    randAmount,
		IsCredit:  randAmount >= 0,
		Status:    "PROCESSING",
	}
}

func processTransfer(conn *pgx.Conn, transfer *transfers.Transfer) error {
	err := crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return transfers.CreateTransfer(tx, transfer)
	})
	fmt.Println("Transfer created!")

	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		err := accounts.UpdateBalance(tx, testAccountId, transfer.Amount)
		if err != nil {
			fmt.Println("err update balance", err)
			return errors.New("balanceerror")
		}

		err = transfers.UpdateTransfer(tx, transfer.Id, "APPROVED")
		if err != nil {
			fmt.Println("err updating transfer", err)
			return err
		}

		return nil
	})

	if err != nil {
		if err.Error() == "balanceerror" {
			fmt.Println("Error trying to update balance, cancelling the transaction")
			return crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
				return transfers.UpdateTransfer(tx, transfer.Id, "CANCELLED")
			})
		}
		fmt.Println("processing error", err)
	}

	fmt.Println("Balance updated")

	return nil
}

func main() {
	conn := database.GetConnection(os.Getenv("DATABASE_URI"))
	defer conn.Close(context.Background())

	for i := 0; i < 100; i++ {
		transfer := generateRandomTransfer()
		fmt.Println("Transfer: ", transfer)

		processTransfer(conn, &transfer)
		fmt.Println("Processed!", transfer)
	}
}
