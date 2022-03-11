package repository

import (
	"context"
	"github.com/diez37/go-packages/clients/db"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"io"
	"time"
)

const (
	sqlTableName = "logins"
)

type sql struct {
	db     goqu.SQLDatabase
	tracer trace.Tracer
}

func NewSql(db goqu.SQLDatabase, tracer trace.Tracer) Repository {
	return &sql{db: db, tracer: tracer}
}

func (repository *sql) FindByUuid(ctx context.Context, uuid uuid.UUID) (*Login, error) {
	ctx, span := repository.tracer.Start(ctx, "repository.sql:FindByUuid")
	defer span.End()

	span.SetAttributes(attribute.String("uuid", uuid.String()))

	sql, args, err := goqu.From(sqlTableName).Where(goqu.Ex{"uuid": uuid}).ToSQL()
	if err != nil {
		return nil, err
	}

	return repository.find(ctx, sql, args...)
}

func (repository *sql) FindByLogin(ctx context.Context, login string) (*Login, error) {
	ctx, span := repository.tracer.Start(ctx, "repository.sql:FindByLogin")
	defer span.End()

	span.SetAttributes(attribute.String("login", login))

	sql, args, err := goqu.From(sqlTableName).Where(goqu.Ex{"login": login}).ToSQL()
	if err != nil {
		return nil, err
	}

	return repository.find(ctx, sql, args...)
}

func (repository *sql) find(ctx context.Context, sql string, args ...interface{}) (*Login, error) {
	ctx, span := repository.tracer.Start(ctx, "repository.sql:find")
	defer span.End()

	rows, err := repository.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		login := &Login{}

		if err := rows.Scan(&login.Uuid, &login.Login, &login.Ban, &login.CreatedAt, &login.UpdateAt); err != nil {
			return nil, err
		}

		return login, nil
	}

	return nil, db.RecordNotFoundError
}

func (repository *sql) Save(ctx context.Context, login *Login) (*Login, error) {
	ctx, span := repository.tracer.Start(ctx, "repository.sql:Save")
	defer span.End()

	var sql string
	var args []interface{}
	var err error

	now := time.Now()

	if login.Uuid == uuid.Nil {
		span.SetAttributes(attribute.String("action", "insert"))

		login.Uuid = uuid.New()

		login.CreatedAt = &now

		sql, args, err = goqu.Insert(sqlTableName).Rows(login).ToSQL()
	} else {
		span.SetAttributes(attribute.String("action", "update"))

		login.UpdateAt = &now

		sql, args, err = goqu.Update(sqlTableName).Set(login).Where(goqu.Ex{"uuid": login.Uuid}).ToSQL()
	}

	if err != nil {
		return nil, err
	}

	_, err = repository.db.ExecContext(ctx, sql, args...)

	return login, err
}

func (repository *sql) BanByUuid(ctx context.Context, uuid uuid.UUID) (bool, error) {
	ctx, span := repository.tracer.Start(ctx, "repository.sql:BanByUuid")
	defer span.End()

	span.SetAttributes(attribute.String("uuid", uuid.String()))

	sql, args, err := goqu.Update(sqlTableName).
		Set(goqu.Record{"ban": true, "update_at": time.Now()}).
		Where(goqu.Ex{"uuid": uuid}).ToSQL()
	if err != nil {
		return false, err
	}

	return repository.ban(ctx, sql, args...)
}

func (repository *sql) BanByLogin(ctx context.Context, login string) (bool, error) {
	ctx, span := repository.tracer.Start(ctx, "repository.sql:BanByLogin")
	defer span.End()

	span.SetAttributes(attribute.String("login", login))

	sql, args, err := goqu.Update(sqlTableName).
		Set(goqu.Record{"ban": true, "update_at": time.Now()}).
		Where(goqu.Ex{"login": login}).ToSQL()
	if err != nil {
		return false, err
	}

	return repository.ban(ctx, sql, args...)
}

func (repository *sql) ban(ctx context.Context, sql string, args ...interface{}) (bool, error) {
	ctx, span := repository.tracer.Start(ctx, "repository.sql:ban")
	defer span.End()

	result, err := repository.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return false, err
	}

	countUpdatedRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if countUpdatedRows == 0 {
		return false, db.RecordNotFoundError
	}

	return true, nil
}

func (repository *sql) Count(ctx context.Context) (int64, error) {
	ctx, span := repository.tracer.Start(ctx, "repository.sql:Count")
	defer span.End()

	sql, _, err := goqu.From(sqlTableName).Select(goqu.COUNT("uuid")).ToSQL()
	if err != nil {
		return 0, err
	}

	rows, err := repository.db.QueryContext(ctx, sql)
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		count := int64(0)

		if err := rows.Scan(&count); err != nil {
			return 0, err
		}

		return count, nil
	}

	return 0, nil
}

func (repository *sql) Page(ctx context.Context, page uint, limit uint) ([]*Login, error) {
	ctx, span := repository.tracer.Start(ctx, "repository.sql:Page")
	defer span.End()

	span.SetAttributes(attribute.Int("page", int(page)), attribute.Int("limit", int(limit)))

	sql, args, err := goqu.From(sqlTableName).Offset(page * limit).Limit(limit).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := repository.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var logins []*Login

	for rows.Next() {
		login := &Login{}

		if err := rows.Scan(&login.Uuid, &login.Login, &login.Ban, &login.CreatedAt, &login.UpdateAt); err != nil {
			return nil, err
		}

		logins = append(logins, login)
	}

	if len(logins) == 0 {
		return nil, io.EOF
	}

	return logins, nil
}
