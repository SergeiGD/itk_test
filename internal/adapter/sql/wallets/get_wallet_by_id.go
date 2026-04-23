package wallets

import (
	"context"
	"errors"

	"github.com/SergeiGD/golang-template/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *walletsRepository) GetWalletById(ctx context.Context, walletId uuid.UUID) (*models.WalletModel, error) {
	q := `
		SELECT id, balance, created_at FROM wallets
		WHERE id = $1
	`
	item := &models.WalletModel{}

	err := r.client.QueryRow(ctx, q, walletId).Scan(
		&item.Id,
		&item.Balance,
		&item.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, &DatabaseError{Query: q, Err: err}
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWalletNotFound
		}

		return nil, err
	}

	return item, nil
}
