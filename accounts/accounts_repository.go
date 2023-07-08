package accounts

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func UpdateBalance(tx pgx.Tx, accountId string, amount float32) error {
	_, err := tx.Exec(
		context.Background(),
		"UPDATE accounts SET balance = balance + $1 WHERE id = $2",
		amount,
		accountId,
	)

	return err
}
