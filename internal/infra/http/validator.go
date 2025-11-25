package infrahttp

import "context"

type Validator interface {
	Validate(ctx context.Context, s interface{}) error
}
