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
		"insert into transfers (account_id, amount, iscredit, status) values ($1, $2, $3, $4) RETURNING id, status",
		transfer.AccountId,
		transfer.Amount,
		transfer.IsCredit,
		transfer.Status,
	)

	err := row.Scan(&transfer.Id, &transfer.Status)
	if err != nil {
		fmt.Println("err storing transfer", err)
		return err
	}

	return nil
}

func UpdateTransfer(tx pgx.Tx, transfer *Transfer, status string) error {
	row := tx.QueryRow(
		context.Background(),
		"UPDATE transfers SET status = + $1, updated_at=clock_timestamp()  WHERE id = $2 RETURNING status",
		status,
		transfer.Id,
	)

	err := row.Scan(&transfer.Status)
	if err != nil {
		fmt.Println("err updating transfer", err)
		return err
	}

	return nil
}
