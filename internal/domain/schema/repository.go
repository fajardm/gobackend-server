package schema

import "context"

type Repository interface {
	All(ctx context.Context) (Schemas, error)
	FindByClassName(ctx context.Context, className ClassName) (*Schema, error)
	Exists(ctx context.Context, className ClassName) (bool, error)
	Create(ctx context.Context, schema Schema) error
	Update(ctx context.Context, schema Schema) error
	Delete(ctx context.Context, className ClassName) error
}
