package gcpag

import (
	"fmt"
	"strings"
)

func StringToOrder(vl string) (Order, error) {
	vl = strings.ToLower(vl)
	switch vl {
	case "asc":
		return Asc, nil
	case "desc":
		return Desc, nil
	default:
		return Order(-1), fmt.Errorf("gcpag: invalid order value '%s'", vl)
	}
}
