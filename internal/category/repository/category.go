package repository

import (
	"context"

	"ecom/internal/category/entity"
	"ecom/pkg/dbs"
)

type CategoryRepo struct {
	db *dbs.Database
}

func NewCategoryRepository(db *dbs.Database) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, category *entity.Category) error {
	return r.db.Create(ctx, category)
}

func (r *CategoryRepo) Update(ctx context.Context, category *entity.Category) error {
	return r.db.Update(ctx, category)
}

func (r *CategoryRepo) Delete(ctx context.Context, id string) error {
	return r.db.Delete(ctx, &entity.Category{}, dbs.WithQuery(dbs.NewQuery("id = ?", id)))
}

func (r *CategoryRepo) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	var category entity.Category
	if err := r.db.FindById(ctx, id, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepo) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	var category entity.Category
	query := dbs.NewQuery("slug = ?", slug)
	if err := r.db.FindOne(ctx, &category, dbs.WithQuery(query)); err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepo) List(ctx context.Context) ([]entity.Category, error) {
	var categories []entity.Category
	if err := r.db.Find(ctx, &categories, dbs.WithOrder("name")); err != nil {
		return nil, err
	}
	return categories, nil
}
