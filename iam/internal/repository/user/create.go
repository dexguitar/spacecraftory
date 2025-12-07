package user

import (
	"context"
	"errors"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
)

type userID struct {
	ID string `db:"id"`
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil && errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	// Create user in users table
	userInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("login", "email", "notification_methods", "password", "created_at", "updated_at").
		Values(user.Info.Login, user.Info.Email, user.Info.NotificationMethods, user.Password, time.Now(), time.Now()).
		Suffix("RETURNING id")

	query, args, err := userInsert.ToSql()
	if err != nil {
		return "", err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[userID])
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return result.ID, nil
}
