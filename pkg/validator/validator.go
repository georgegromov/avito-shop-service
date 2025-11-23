package validator

import (
	"avito-shop-service/internal/infra/http"
	"context"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {

	opts := []validator.Option{
		validator.WithRequiredStructEnabled(),
	}

	v := validator.New(opts...)
	return &Validator{validate: v}
}

func (v *Validator) Validate(ctx context.Context, s interface{}) error {
	return v.validate.StructCtx(ctx, s)
}

var _ http.Validator = (*Validator)(nil)
