package postgres

import (
	"context"
	"fmt"

	"github.com/Murilovisque/gocore/gcrepo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildPool(ctx context.Context, databaseUrl string) (gcrepo.SqlResource, error) {
	dbpool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to create postgres connection pool. Cause: %w", err)
	}
	return &postgresResource{dbPoll: dbpool}, nil
}

type postgresResource struct {
	dbPoll *pgxpool.Pool
	tx     pgx.Tx
}

func (p *postgresResource) Executor() gcrepo.SqlExecutor {
	return p
}

func (p *postgresResource) Begin(ctx context.Context) (gcrepo.SqlTxResource, error) {
	tx, err := p.dbPoll.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &postgresResource{dbPoll: p.dbPoll, tx: tx}, nil
}

func (p *postgresResource) Commit(ctx context.Context) error {
	if p.tx == nil {
		return fmt.Errorf("no postgres transaction to commit")
	}
	return p.tx.Commit(ctx)
}

func (p *postgresResource) Rollback(ctx context.Context) error {
	if p.tx == nil {
		return fmt.Errorf("no postgres transaction to rollback")
	}
	return p.tx.Rollback(ctx)
}

func (p *postgresResource) Syntax() gcrepo.SqlSyntax {
	return postgresSyntax{}
}

func (p *postgresResource) Close() error {
	if p.tx == nil {
		p.dbPoll.Close()
	}
	return nil
}

func (p *postgresResource) Exec(ctx context.Context, sql string, args ...any) (gcrepo.SqlResult, error) {
	if p.tx != nil {
		result, err := p.tx.Exec(ctx, sql, args...)
		if err != nil {
			return nil, err
		}
		return postgresExecResult{rowsAffected: result.RowsAffected()}, nil
	} else {
		result, err := p.dbPoll.Exec(ctx, sql, args...)
		if err != nil {
			return nil, err
		}
		return postgresExecResult{rowsAffected: result.RowsAffected()}, nil
	}
}

func (p *postgresResource) Query(ctx context.Context, sql string, args ...any) (gcrepo.SqlRows, error) {
	if p.tx != nil {
		rows, err := p.tx.Query(ctx, sql, args...)
		if err != nil {
			return nil, err
		}
		return &postgresRows{rows: rows}, nil
	} else {
		rows, err := p.dbPoll.Query(ctx, sql, args...)
		if err != nil {
			return nil, err
		}
		return &postgresRows{rows: rows}, nil
	}
}

func (p *postgresResource) QueryRow(ctx context.Context, sql string, args ...any) gcrepo.SqlRow {
	if p.tx != nil {
		return p.tx.QueryRow(ctx, sql, args...)
	} else {
		return p.dbPoll.QueryRow(ctx, sql, args...)
	}
}

// rows
type postgresRows struct {
	rows pgx.Rows
}

func (p *postgresRows) Close() error {
	p.rows.Close()
	return nil
}

func (p *postgresRows) Err() error {
	return p.rows.Err()
}

func (p *postgresRows) Next() bool {
	return p.rows.Next()
}

func (p *postgresRows) Scan(args ...any) error {
	return p.rows.Scan(args...)
}

// exec
type postgresExecResult struct {
	rowsAffected int64
}

func (p postgresExecResult) RowsAffected() (int64, error) {
	return p.rowsAffected, nil
}

// syntax
type postgresSyntax struct{}

func (s postgresSyntax) PlaceHolder(i int) string {
	return fmt.Sprintf("$%d", i)
}
func (s postgresSyntax) LimitStatement(i int) string {
	return fmt.Sprintf("LIMIT $%d", i)
}
