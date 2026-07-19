package postgres

import (
	"context"
	"fmt"
	"net"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
	"github.com/Murilovisque/gocore/gcpag"
	"github.com/Murilovisque/gocore/gcrepo"
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
			t.Fatalf("erro no insert: %v", err)
		}
		row := exec.QueryRow(ctx, "SELECT name FROM test_users WHERE name = $1", "Murilo")
		var name string
		if err := row.Scan(&name); err != nil {
			t.Fatalf("erro no scan: %v", err)
		}

		if name != "Murilo" {
			t.Fatalf("expected Murilo, but '%s'", name)
		}
	})

	t.Run("InsertAndSelectPaginated", func(t *testing.T) {
		expectedNames := []any{"Bola", "Bolinha", "Baiana", "Bolinha", "Bunito"}
		_, err := exec.Exec(ctx, "INSERT INTO test_users (name) VALUES ($1), ($2), ($3), ($4), ($5)", expectedNames...)
		if err != nil {
			t.Fatalf("erro no insert: %v", err)
		}
		const baseSelectQuery = "SELECT id, name FROM test_users WHERE name like $1"
		baseSelectArgs := []any{"B%"}
		const pageSize = 2
		testPage := func(pageReq gcpag.PaginatedRequest[testUserIdt], skipNames, exptItemsSize int) gcpag.PaginatedResponse[testUserIdt, testUserModel] {
			{ //log table state
				rows, err := exec.Query(t.Context(), baseSelectQuery, baseSelectArgs...)
				if err != nil {
					t.Fatal(err)
				}
				targetRows, err := gcrepo.CollectRows(rows, func(row gcrepo.SqlRow) (map[string]any, error) {
					var id int
					var name string
					err := row.Scan(&id, &name)
					target := map[string]any{
						"id":   id,
						"name": name,
					}
					return target, err
				})
				if err != nil {
					t.Fatal(err)
				}
				t.Log("table state base query:", targetRows)
			}
			reqArg := gcrepo.QueryPaginatedParams[testUserIdt, testUserModel]{
				ScanRow: func(row *testUserModel) []any {
					return []any{&row.idt, &row.name}
				},
				IdtColumn:  "id",
				QueryItems: baseSelectQuery,
				QueryArgs:  baseSelectArgs,
				FieldColumn: gcopt.Of(func(fld gcfield.FieldNameOrdered, placeHolder int) gcopt.Optional[gcrepo.SubQueryPaginatedFieldOrdered] {
					if fld != "name" {
						return gcopt.Empty[gcrepo.SubQueryPaginatedFieldOrdered]()
					}
					return gcopt.Of(gcrepo.SubQueryPaginatedFieldOrdered{
						SubQuery:   "select name from test_users where id = $" + strconv.Itoa(placeHolder),
						ColumnName: "name",
					})
				}),
				LastQueryPlaceHolder: 1,
			}
			response, err := gcrepo.QueryPaginated(t.Context(), exec, pageReq, reqArg)
			t.Logf("response %v", response)
			if err != nil {
				t.Fatal(err)
			}
			if len(response.Items) != exptItemsSize {
				t.Fatalf("expected total items %d, but %d", exptItemsSize, len(response.Items))
			}
			for i := 0; i < len(response.Items); i++ {
				if response.Items[i].name != expectedNames[i+skipNames] {
					t.Fatalf("expected '%s', but '%s'", expectedNames[i], response.Items[i].name)
				}
			}
			if page, ok := response.SelfPage.Take(); !ok {
				t.Fatal("expected self page")
			} else {
				expec := response.Items[0].idt
				vl := page.Idt
				if expec != page.Idt {
					t.Fatalf("expected %v, but %v", expec, vl)
				}
			}
			return response
		}
		testNavegation := func(pageReq gcpag.PaginatedRequest[testUserIdt]) {
			t.Log("test first page")
			res := testPage(pageReq, 0, 2)
			if res.Items[0].name != expectedNames[0] || res.Items[1].name != expectedNames[1] {
				t.Fatalf("expected same items")
			} else if res.FirstPage.MustTake().Idt != res.Items[0].Idt() {
				t.Fatalf("expected page is be the first: %d not match %d", res.FirstPage.MustTake().Idt, res.Items[0].idt)
			} else if res.LastPage.IsPresent() {
				t.Fatalf("expected last page is empty")
			} else if res.PreviousPage.IsPresent() {
				t.Fatal("expected first page does not have previous")
			}

			t.Log("test second page")
			res = testPage(res.NextPageAsRequest().MustTake(), 2, 2)
			if res.Items[0].name != expectedNames[2] || res.Items[1].name != expectedNames[3] {
				t.Fatalf("expected same items")
			} else if res.FirstPage.IsPresent() {
				t.Fatalf("expected first page is empty")
			} else if res.LastPage.IsPresent() {
				t.Fatalf("expected last page is empty")
			} else if !res.PreviousPage.IsPresent() {
				t.Fatal("expected previous page in second page")
			}

			t.Log("test last page")
			res = testPage(res.NextPageAsRequest().MustTake(), 4, 1)
			if res.Items[0].name != expectedNames[4] {
				t.Fatalf("expected same items")
			} else if res.FirstPage.IsPresent() {
				t.Fatalf("expected first page is empty")
			} else if res.LastPage.MustTake().Idt != res.Items[0].idt {
				t.Fatalf("expected page is be the last: %d not match %d", res.LastPage.MustTake().Idt, res.Items[0].idt)
			} else if res.NextPage.IsPresent() {
				t.Fatal("expected last page does not have next")
			}

			t.Log("test returning to second page")
			res = testPage(res.PreviousPageAsRequest().MustTake(), 2, 2)
			if res.Items[0].name != expectedNames[2] || res.Items[1].name != expectedNames[3] {
				t.Fatalf("expected same items")
			} else if res.FirstPage.IsPresent() {
				t.Fatalf("expected first page is empty")
			} else if res.LastPage.IsPresent() {
				t.Fatalf("expected last page is empty")
			} else if !res.NextPage.IsPresent() {
				t.Fatal("expected second page to have next")
			}

			t.Log("test returning to first page")
			res = testPage(res.PreviousPageAsRequest().MustTake(), 0, 2)
			if res.Items[0].name != expectedNames[0] || res.Items[1].name != expectedNames[1] {
				t.Fatalf("expected same items")
			} else if res.FirstPage.MustTake().Idt != res.Items[0].Idt() {
				t.Fatalf("expected page is be the first: %d not match %d", res.FirstPage.MustTake().Idt, res.Items[0].idt)
			} else if res.LastPage.IsPresent() {
				t.Fatalf("expected last page is empty")
			} else if res.PreviousPage.IsPresent() {
				t.Fatal("expected first page does not have previous")
			} else if !res.NextPage.IsPresent() {
				t.Fatal("expected second page to have next")
			}
		}
		t.Log("ascending test")
		expectedNames = []any{"Bola", "Bolinha", "Baiana", "Bolinha", "Bunito"}
		pageReq := gcpag.PaginatedRequest[testUserIdt]{
			Size: pageSize,
		}
		testNavegation(pageReq)

		t.Log("descending test")
		slices.Reverse(expectedNames)
		pageReq = gcpag.PaginatedRequest[testUserIdt]{
			Size:  pageSize,
			Order: gcpag.Desc,
		}
		testNavegation(pageReq)

		t.Log("ascending by name test")
		expectedNames = []any{"Baiana", "Bola", "Bolinha", "Bolinha", "Bunito"}
		pageReq = gcpag.PaginatedRequest[testUserIdt]{
			Size:      pageSize,
			SortField: gcopt.Of(gcfield.FieldNameOrdered("name")),
		}
		testNavegation(pageReq)
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
		sql := fmt.Sprintf("INSERT INTO test_users (name) VALUES (%s), (%s), (%s)", gcrepo.PlaceHolderRange(exec.Syntax(), 1, 4)...)
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

type testUserIdt int

func (t testUserIdt) String() string {
	return strconv.Itoa(int(t))
}

type testUserModel struct {
	idt  testUserIdt
	name string
}

func (t testUserModel) Idt() testUserIdt {
	return t.idt
}
