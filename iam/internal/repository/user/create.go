package user

import (
	"context"
	"errors"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

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
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	// create user
	userInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("login", "email", "password", "created_at", "updated_at").
		Values(user.Info.Login, user.Info.Email, user.Password, time.Now(), time.Now()).
		Suffix("RETURNING id")

	query, args, err := userInsert.ToSql()
	if err != nil {
		return "", err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return "", model.ErrUserAlreadyExists
		}
		return "", err
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[userID])
	if err != nil {
		return "", err
	}

	// insert notification methods
	if len(user.Info.NotificationMethods) > 0 {
		notificationInsert := sq.Insert("notification_methods").
			PlaceholderFormat(sq.Dollar).
			Columns("user_uuid", "provider_name", "target")

		for _, method := range user.Info.NotificationMethods {
			notificationInsert = notificationInsert.Values(result.ID, method.ProviderName, method.Target)
		}

		query, args, err = notificationInsert.ToSql()
		if err != nil {
			return "", err
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return "", err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return "", err
	}

	return result.ID, nil
}
