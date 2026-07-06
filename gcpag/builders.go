package gcpag

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Murilovisque/gocore/gcfield"
	"github.com/Murilovisque/gocore/gcopt"
)

func ParseRequestFromHttp[I gcfield.IdtOrdered](req *http.Request, params ParseRequestFromHttpParams[I]) (PaginatedRequest[I], error) {
	q := req.URL.Query()
	size, err := strconv.Atoi(q.Get(httpParamPageSize))
	if size < 1 || err != nil {
		size = params.DefaultSize
	}
	order, err := ParseOrder(q.Get(httpParamPageOrder))
	if err != nil {
		order = params.DefaultOrder
	}
	p := PaginatedRequest[I]{Size: size, Order: order}
	checkPoints := []struct {
		field string
		pos   StartPosition
		ori   Orientation
	}{
		{httpParamPageStartIdt, StartAt, NextPage},
		{httpParamPageAfterIdt, AfterAt, NextPage},
		{httpParamReversePageStartIdt, StartAt, PreviousPage},
		{httpParamReversePageAfterIdt, AfterAt, PreviousPage},
	}
	for _, cp := range checkPoints {
		if val := q.Get(cp.field); val != "" {
			if idt, ok, err := params.IdtParser(val); err != nil {
				return p, err
			} else if ok {
				p.Idt = gcopt.Of(idt)
				p.StartPosition = cp.pos
				p.Orientation = cp.ori
				break
			}
		}
	}
	if fp, ok := params.Field.Take(); ok {
		for _, name := range fp.AllowedNames() {
			if val := q.Get(httpParamPageSortField + name); val != "" {
				if ok, err := fp.Parse(name, val); err != nil {
					return p, err
				} else if ok {
					break
				}
			}
		}
	}
	return p, nil
}


func ParseOrder(vl string) (Order, error) {
	vl = strings.ToLower(vl)
	switch vl {
	case "asc":
		return Asc, nil
	case "desc":
		return Desc, nil
	default:
		return Asc, fmt.Errorf("gcpag: invalid order value '%s'", vl)
	}
}
