package postgres

import (
	"context"
	"fmt"
	"net"
	"slices"
	"testing"
	"time"

	"github.com/Murilovisque/gocore/gcr"
)

func TestPostgresIntegration(t *testing.T) {
	ctx := context.Background()
	connStr := "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"

	checkPostgresAvailable(t)

	res, err := BuildPool(ctx, connStr)
	if err != nil {
		t.Fatalf("falha ao conectar no banco: %v", err)
	}

	exec := res.Executor()
	_, err = exec.Exec(ctx, `CREATE TABLE IF NOT EXISTS test_users (id SERIAL PRIMARY KEY, name TEXT)`)
	if err != nil {
		t.Fatalf("falha ao criar tabela: %v", err)
	}

	t.Cleanup(func() {
		_, err := exec.Exec(ctx, `DROP TABLE test_users`)
		if err != nil {
			t.Fatal(err)
		}
		err = res.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("InsertAndSelect", func(t *testing.T) {
		_, err := exec.Exec(ctx, "INSERT INTO test_users (name) VALUES ($1)", "Murilo")
		if err != nil {
			t.Errorf("erro no insert: %v", err)
		}

		row := exec.QueryRow(ctx, "SELECT name FROM test_users WHERE name = $1", "Murilo")
		var name string
		if err := row.Scan(&name); err != nil {
			t.Errorf("erro no scan: %v", err)
		}

		if name != "Murilo" {
			t.Errorf("expected Murilo, but '%s'", name)
		}
	})

	t.Run("TransactionCommit", func(t *testing.T) {
		txRes, err := res.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		_, err = txRes.Executor().Exec(ctx, "INSERT INTO test_users (name) VALUES ($1)", "TransactionUser")
		if err != nil {
			txRes.Rollback(ctx)
			t.Fatal(err)
		}
		if err := txRes.Commit(ctx); err != nil {
			t.Fatal(err)
		}
		var count int
		row := exec.QueryRow(ctx, "SELECT count(*) FROM test_users WHERE name = $1", "TransactionUser")
		row.Scan(&count)
		if count != 1 {
			t.Errorf("transaction not persisted")
		}
	})

	t.Run("SyntaxAndRows", func(t *testing.T) {
		const expectedLimit = 2
		var expectedNames = []any{"zulu", "zé", "zacarias"}
		sql := fmt.Sprintf("INSERT INTO test_users (name) VALUES (%s), (%s), (%s)", gcr.PlaceHolderRange(exec.Syntax(), 1, 4)...)
		_, err := exec.Exec(ctx, sql, expectedNames...)
		if err != nil {
			t.Errorf("erro no insert: %v", err)
		}
		result := []any{}
		sql = fmt.Sprintf("SELECT name FROM test_users where name like %s order by id %s", exec.Syntax().PlaceHolder(1), exec.Syntax().LimitStatement(2))
		rows, err := exec.Query(ctx, sql, "z%", expectedLimit)
		if err != nil {
			t.Fatal(err)
		}
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				t.Errorf("erro no scan: %v", err)
			}
			result = append(result, name)
		}
		if err = rows.Err(); err != nil {
			t.Fatal(err)
		} else if err = rows.Close(); err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(expectedNames[:2], result) {
			t.Fatalf("expected %v, but %v", expectedNames, result)
		}
	})
}

func checkPostgresAvailable(t *testing.T) {
	t.Helper()

	timeout := 1 * time.Second
	_, err := net.DialTimeout("tcp", "localhost:5432", timeout)
	if err != nil {
		t.Skip("postgres is not running. Skipping integration tests")
	}
}
