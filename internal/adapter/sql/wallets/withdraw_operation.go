package wallets

import (
	"context"
	"errors"
	"fmt"

	"github.com/SergeiGD/golang-template/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
)

func (r *walletsRepository) WithdrawFromWallet(ctx context.Context, walletId uuid.UUID, amount decimal.Decimal) (*models.WalletModel, error) {
	updateQ := `
		UPDATE wallets
		SET balance = balance-$1
		WHERE id = $2
		RETURNING id, balance, created_at
	`
	insertQ := `
		INSERT into operations
		(wallet_id, type, amount)
		VALUES ($1, $2, $3)
	`

	tx, err := r.client.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("client.Begin: %w", err)
	}
	defer func() {
		if err != nil && tx != nil {
			if transErr := tx.Rollback(ctx); transErr != nil {
				r.logger.WithContext(ctx).
					Warn("failed to rollback transaction: %w", transErr)
			}
		}
	}()

	item := &models.WalletModel{}
	err = tx.QueryRow(ctx, updateQ, amount, walletId).Scan(
		&item.Id,
		&item.Balance,
		&item.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, &DatabaseError{Query: updateQ, Err: err}
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}

	_, err = tx.Exec(ctx, insertQ, walletId, models.WithdrawOperation, amount)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, &DatabaseError{Query: updateQ, Err: err}
		}
		return nil, err
	}

	err = tx.Commit(ctx)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, &DatabaseError{Query: updateQ, Err: err}
		}
		return nil, err
	}

	return item, nil
}
