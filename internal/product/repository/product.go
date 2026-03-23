package repository

import (
	"context"

	"ecom/internal/product/entity"
	"ecom/pkg/dbs"
)

type ProductRepo struct {
	db *dbs.Database
}

func NewProductRepository(db *dbs.Database) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, product *entity.Product) error {
	return r.db.Create(ctx, product)
}

func (r *ProductRepo) Update(ctx context.Context, product *entity.Product) error {
	return r.db.Update(ctx, product)
}

func (r *ProductRepo) Delete(ctx context.Context, id string) error {
	return r.db.Delete(ctx, &entity.Product{}, dbs.WithQuery(dbs.NewQuery("id = ?", id)))
}

func (r *ProductRepo) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.FindOne(ctx, &product,
		dbs.WithQuery(dbs.NewQuery("products.id = ?", id)),
		dbs.WithPreload([]string{"Variants", "NutritionalInfo", "Category"}),
	); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepo) GetBySlug(ctx context.Context, slug string) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.FindOne(ctx, &product,
		dbs.WithQuery(dbs.NewQuery("products.slug = ?", slug)),
		dbs.WithPreload([]string{"Variants", "NutritionalInfo", "Category"}),
	); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepo) List(ctx context.Context, opts ...dbs.FindOption) ([]entity.Product, error) {
	var products []entity.Product
	allOpts := append([]dbs.FindOption{
		dbs.WithPreload([]string{"Variants", "NutritionalInfo", "Category"}),
		dbs.WithOrder("created_at DESC"),
	}, opts...)

	if err := r.db.Find(ctx, &products, allOpts...); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepo) Count(ctx context.Context, opts ...dbs.FindOption) (int64, error) {
	var total int64
	if err := r.db.Count(ctx, &entity.Product{}, &total, opts...); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *ProductRepo) GetVariantByID(ctx context.Context, id string) (*entity.ProductVariant, error) {
	var variant entity.ProductVariant
	if err := r.db.FindOne(ctx, &variant, dbs.WithQuery(dbs.NewQuery("id = ?", id))); err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *ProductRepo) UpdateVariant(ctx context.Context, variant *entity.ProductVariant) error {
	return r.db.Update(ctx, variant)
}

func (r *ProductRepo) DeleteVariantsByProductID(ctx context.Context, productID string) error {
	return r.db.Delete(ctx, &entity.ProductVariant{}, dbs.WithQuery(dbs.NewQuery("product_id = ?", productID)))
}

func (r *ProductRepo) DeleteNutritionalInfoByProductID(ctx context.Context, productID string) error {
	return r.db.Delete(ctx, &entity.NutritionalInfo{}, dbs.WithQuery(dbs.NewQuery("product_id = ?", productID)))
}

func (r *ProductRepo) GetDB() *dbs.Database {
	return r.db
}
