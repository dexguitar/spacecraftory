package user

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/dexguitar/spacecraftory/iam/internal/model"
	"github.com/dexguitar/spacecraftory/iam/internal/repository/converter"
	repoModel "github.com/dexguitar/spacecraftory/iam/internal/repository/model"
)

func (r *userRepository) GetUserByUUID(ctx context.Context, userUUID string) (*model.User, error) {
	userQuery := sq.Select("id", "login", "email", "password").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": userUUID})

	query, args, err := userQuery.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[repoModel.UserRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	methods, err := r.getNotificationMethods(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	return converter.ToModelUserFromRow(&row, methods), nil
}

func (r *userRepository) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	userQuery := sq.Select("id", "login", "email", "password").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"login": login})

	query, args, err := userQuery.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[repoModel.UserRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	methods, err := r.getNotificationMethods(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	return converter.ToModelUserFromRow(&row, methods), nil
}

func (r *userRepository) getNotificationMethods(ctx context.Context, userUUID string) ([]repoModel.NotificationMethodRow, error) {
	query := sq.Select("id", "user_uuid", "provider_name", "target").
		From("notification_methods").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"user_uuid": userUUID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[repoModel.NotificationMethodRow])
}
