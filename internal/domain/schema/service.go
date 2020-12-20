package schema

import "context"

type Service interface {
	FindByClassName(ctx context.Context, className ClassName) (*Schema, error)
	Create(ctx context.Context, data Schema) (*Schema, error)
	Update(ctx context.Context, data Schema) (*Schema, error)
	Delete(ctx context.Context, className ClassName) error
}
