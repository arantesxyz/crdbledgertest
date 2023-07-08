package transfers

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Transfer struct {
	Id        string
	AccountId string
	Amount    float32
	IsCredit  bool
	Status    string
}

func CreateTransfer(tx pgx.Tx, transfer *Transfer) error {
	row := tx.QueryRow(
		context.Background(),
		"insert into transfers (account_id, amount, iscredit, status) values ($1, $2, $3, $4) RETURNING id",
		transfer.AccountId,
		transfer.Amount,
		transfer.IsCredit,
		transfer.Status,
	)

	err := row.Scan(&transfer.Id)
	if err != nil {
		fmt.Println("err create transfer", err)
		return err
	}

	return nil
}

func UpdateTransfer(tx pgx.Tx, txid string, status string) error {
	_, err := tx.Exec(
		context.Background(),
		"UPDATE transfers SET status = + $1 WHERE id = $2",
		status,
		txid,
	)

	return err
}
