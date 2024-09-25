package endpoint

import "errors"

var (
	ErrPaginationLimit  = errors.New("pagniation limit error, must be between 10 and 100")
	ErrUnsupportedOrder = errors.New("unsupported order, must be ASC or DESC")
	ErrMissingOrder     = errors.New("order is missing, must be ASC or DESC")
	ErrMissingLimit     = errors.New("limit is missing, must be between 10 and 100")
)

const (
	MaxPaginationLimit = 100
	MinPaginationLimit = 10
)

type PagniationParams struct {
	Limit  int    `query:"limit"`
	Cursor string `query:"cursor"`
	Order  string `query:"order"`
}
