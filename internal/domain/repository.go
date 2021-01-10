package domain

import "context"

type Repository interface {
	ClassExist(ctx context.Context, className ClassName) (bool, error)
	GetAllClasses(ctx context.Context) (Classes, error)
	GetClass(ctx context.Context, className ClassName) (*Class, error)
	CreateClass(ctx context.Context, class Class) error
	UpdateClass(ctx context.Context, class Class) error
}
