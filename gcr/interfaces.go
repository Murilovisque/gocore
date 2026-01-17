package gcr

import (
	"context"
)

type SqlResource interface {
	Executor() SqlExecutor
	Begin(ctx context.Context) (SqlTxResource, error)
	Syntax() SqlSyntax
	Close() error
}

type SqlTxResource interface {
	SqlResource
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type SqlExecutor interface {
	Exec(ctx context.Context, sql string, args ...any) (SqlResult, error)
	Query(ctx context.Context, sql string, args ...any) (SqlRows, error)
	QueryRow(ctx context.Context, sql string, args ...any) SqlRow
	Syntax() SqlSyntax
}

type SqlSyntax interface {
	PlaceHolder(position int) string
	LimitStatement(position int) string
}

type SqlResult interface {
	RowsAffected() (int64, error)
}

type SqlRows interface {
	Close() error
	Err() error
	Next() bool
	Scan(...any) error
}

type SqlRow interface {
	Scan(...any) error
}
