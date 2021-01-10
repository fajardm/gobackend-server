package domain

import (
	"context"
)

type Service interface {
	FindByClassName(ctx context.Context, className ClassName) (*Class, error)
	Create(ctx context.Context, data Class) (*Class, error)
	Update(ctx context.Context, data Class) (*Class, error)
	Delete(ctx context.Context, className ClassName) error
}

// type UpdatedResult struct {
// 	UpdatedClass   Class
// 	AddedFields    Fields
// 	DeletedFields  Fields
// 	AddedIndexes   Indexes
// 	DeletedIndexes Indexes
// }

// func Update(existing Class, newClass Class) *UpdatedResult {
// 	var (
// 		wg  = sync.WaitGroup{}
// 		res = new(UpdatedResult)
// 	)

// 	go func() {
// 		wg.Add(1)
// 		defer wg.Done()
// 		for key, field := range newClass.Fields {
// 			if _, ok := existing.Fields[key]; !ok {
// 				existing.Fields[key] = field
// 				res.AddedFields[key] = field
// 			}
// 		}
// 	}()

// 	go func() {
// 		wg.Add(1)
// 		defer wg.Done()
// 		for key, field := range existing.Fields {
// 			if _, ok := newClass.Fields[key]; !ok {
// 				existing.Fields.Delete(key)
// 				res.DeletedFields[key] = field
// 			}
// 		}
// 	}()

// 	go func() {
// 		wg.Add(1)
// 		defer wg.Done()
// 		for key, index := range newClass.Indexes {
// 			if _, ok := existing.Indexes[key]; !ok {
// 				existing.Indexes[key] = index
// 				res.AddedIndexes[key] = index
// 			}
// 		}
// 	}()

// 	go func() {
// 		wg.Add(1)
// 		defer wg.Done()
// 		for key, index := range existing.Indexes {
// 			if _, ok := newClass.Indexes[key]; !ok {
// 				existing.Indexes.Delete(key)
// 				res.DeletedIndexes[key] = index
// 			}
// 		}
// 	}()

// 	wg.Wait()

// 	return res
// }
