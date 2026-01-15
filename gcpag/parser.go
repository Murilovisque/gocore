package gcpag

import (
	"fmt"
	"strings"
)

func StringToOrder(vl string) (Order, error) {
	o := Order(strings.ToLower(vl))
	if o.IsValid() {
		return o, nil
	}
	return o, fmt.Errorf("gcpag: invalid order value '%s'", vl)
}
