package commons

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
)

type ValidatorDirective struct {
	Validate *validator.Validate
}

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

func (v *ValidatorDirective) Binding(ctx context.Context, obj interface{},
	next graphql.Resolver, constraint string) (interface{}, error) {

	val, err := next(ctx)
	if err != nil {
		panic(err)
	}

	err = Validate.Var(val, constraint)
	if err != nil {
		return nil, err
	}

	return val, nil
}
