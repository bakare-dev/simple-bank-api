package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (repo *Repository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	if err := repo.db.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (repo *Repository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var entity T
	if err := repo.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (repo *Repository[T]) Get(ctx context.Context, opts map[string]any) (*T, error) {
	var entity T
	query := repo.db.WithContext(ctx).Model(new(T))
	for key, value := range opts {
		query = query.Where(key+" = ?", value)
	}
	if err := query.First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (repo *Repository[T]) Update(ctx context.Context, id string, entity *T) error {
	if err := repo.db.WithContext(ctx).Where("id = ?", id).Save(entity).Error; err != nil {
		return err
	}
	return nil
}

func (repo *Repository[T]) Delete(ctx context.Context, id string) error {
	if err := repo.db.WithContext(ctx).Where("id = ?", id).Delete(new(T)).Error; err != nil {
		return err
	}
	return nil
}

func (repo *Repository[T]) GetAll(ctx context.Context, page, size int, opts ...map[string]any) ([]T, error) {
	var entities []T
	query := repo.db.WithContext(ctx).Model(new(T))

	if len(opts) > 0 {
		for key, value := range opts[0] {
			query = query.Where(key+" = ?", value)
		}
	}

	offset := (page - 1) * size
	query = query.Offset(offset).Limit(size)

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (repo *Repository[T]) Count(ctx context.Context, opts ...map[string]any) (int64, error) {
	var total int64
	query := repo.db.WithContext(ctx).Model(new(T))

	if len(opts) > 0 {
		for key, value := range opts[0] {
			query = query.Where(key+" = ?", value)
		}
	}

	err := query.Count(&total).Error
	if err != nil {
		return 0, err
	}

	return total, nil
}
