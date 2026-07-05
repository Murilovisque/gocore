package gcrepo

import "fmt"

type RowTo[T any] func(row SqlRow) (target T, err error)

func CollectRows[T any](rows SqlRows, fnRowTo RowTo[T]) ([]T, error) {
	defer rows.Close()
	var items []T
	for rows.Next() {
		m, err := fnRowTo(rows)
		if err != nil {
			return items, fmt.Errorf("gcrepo: convert item of the rows failed. Cause %w", err)
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return items, fmt.Errorf("gcrepo: convert item after rows iteration failed. Cause %w", err)
	}
	return items, nil
}
