package code

import "github.com/go-kratos/kratos/v2/errors"

var (
	InvalidId = errors.New(404, "INVALID_ID", "Invalid Id")
)
