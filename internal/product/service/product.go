package service

import (
	"context"
	"errors"
	"strings"

	"ecom/internal/product/dto"
	"ecom/internal/product/entity"
	"ecom/internal/product/repository"
	"ecom/pkg/dbs"

	"github.com/quangdangfit/gocommon/logger"
)

type ProductService struct {
	repo *repository.ProductRepo
}

func NewProductService(repo *repository.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, req *dto.CreateProductReq) (*entity.Product, error) {
	slug := generateSlug(req.Name)

	// Check duplicate slug
	existing, _ := s.repo.GetBySlug(ctx, slug)
	if existing != nil {
		return nil, errors.New("product with this name already exists")
	}

	product := entity.Product{
		Name:             req.Name,
		Slug:             slug,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		CategoryID:       req.CategoryID,
		Images:           entity.StringArray(req.Images),
		BasePrice:        req.BasePrice,
		Featured:         req.Featured,
		BestSeller:       req.BestSeller,
	}

	// Build variants
	totalStock := 0
	for _, v := range req.Variants {
		product.Variants = append(product.Variants, entity.ProductVariant{
			Weight: v.Weight,
			Price:  v.Price,
			Stock:  v.Stock,
		})
		totalStock += v.Stock
	}
	product.Stock = totalStock

	// Build nutritional info
	if req.NutritionalInfo != nil {
		product.NutritionalInfo = &entity.NutritionalInfo{
			Calories: req.NutritionalInfo.Calories,
			Protein:  req.NutritionalInfo.Protein,
			Fat:      req.NutritionalInfo.Fat,
			Carbs:    req.NutritionalInfo.Carbs,
			Fiber:    req.NutritionalInfo.Fiber,
		}
	}

	if err := s.repo.Create(ctx, &product); err != nil {
		logger.Errorf("ProductService.Create fail, error: %s", err)
		return nil, err
	}

	// Re-fetch with all relations
	created, err := s.repo.GetByID(ctx, product.ID)
	if err != nil {
		return &product, nil
	}
	return created, nil
}

func (s *ProductService) Update(ctx context.Context, id string, req *dto.UpdateProductReq) (*entity.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if req.Name != "" {
		product.Name = req.Name
		product.Slug = generateSlug(req.Name)
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.ShortDescription != "" {
		product.ShortDescription = req.ShortDescription
	}
	if req.CategoryID != "" {
		product.CategoryID = req.CategoryID
	}
	if len(req.Images) > 0 {
		product.Images = entity.StringArray(req.Images)
	}
	if req.BasePrice > 0 {
		product.BasePrice = req.BasePrice
	}
	if req.Featured != nil {
		product.Featured = *req.Featured
	}
	if req.BestSeller != nil {
		product.BestSeller = *req.BestSeller
	}

	// If variants are provided, replace all existing ones
	if len(req.Variants) > 0 {
		_ = s.repo.DeleteVariantsByProductID(ctx, id)
		product.Variants = nil
		totalStock := 0
		for _, v := range req.Variants {
			product.Variants = append(product.Variants, entity.ProductVariant{
				ProductID: id,
				Weight:    v.Weight,
				Price:     v.Price,
				Stock:     v.Stock,
			})
			totalStock += v.Stock
		}
		product.Stock = totalStock
	}

	// If nutritional info is provided, replace
	if req.NutritionalInfo != nil {
		_ = s.repo.DeleteNutritionalInfoByProductID(ctx, id)
		product.NutritionalInfo = &entity.NutritionalInfo{
			ProductID: id,
			Calories:  req.NutritionalInfo.Calories,
			Protein:   req.NutritionalInfo.Protein,
			Fat:       req.NutritionalInfo.Fat,
			Carbs:     req.NutritionalInfo.Carbs,
			Fiber:     req.NutritionalInfo.Fiber,
		}
	}

	if err := s.repo.Update(ctx, product); err != nil {
		logger.Errorf("ProductService.Update fail, error: %s", err)
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("product not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *ProductService) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) GetBySlug(ctx context.Context, slug string) (*entity.Product, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *ProductService) List(ctx context.Context, categoryID, search string, page, limit int) ([]entity.Product, int64, error) {
	var opts []dbs.FindOption

	if categoryID != "" {
		opts = append(opts, dbs.WithQuery(dbs.NewQuery("category_id = ?", categoryID)))
	}
	if search != "" {
		opts = append(opts, dbs.WithQuery(dbs.NewQuery("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")))
	}

	total, err := s.repo.Count(ctx, opts...)
	if err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	opts = append(opts, dbs.WithOffset(offset), dbs.WithLimit(limit))

	products, err := s.repo.List(ctx, opts...)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (s *ProductService) UpdateInventory(ctx context.Context, id string, req *dto.UpdateInventoryReq) error {
	variant, err := s.repo.GetVariantByID(ctx, req.VariantID)
	if err != nil {
		return errors.New("variant not found")
	}
	if variant.ProductID != id {
		return errors.New("variant does not belong to this product")
	}

	variant.Stock = req.Stock
	if err := s.repo.UpdateVariant(ctx, variant); err != nil {
		return err
	}

	// Recalculate total stock
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil
	}
	totalStock := 0
	for _, v := range product.Variants {
		totalStock += v.Stock
	}
	product.Stock = totalStock
	return s.repo.Update(ctx, product)
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
