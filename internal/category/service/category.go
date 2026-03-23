package service

import (
	"context"
	"errors"
	"strings"

	"ecom/internal/category/dto"
	"ecom/internal/category/entity"
	"ecom/internal/category/repository"
	"ecom/pkg/utils"

	"github.com/quangdangfit/gocommon/logger"
)

type CategoryService struct {
	repo *repository.CategoryRepo
}

func NewCategoryService(repo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, req *dto.CreateCategoryReq) (*entity.Category, error) {
	slug := generateSlug(req.Name)

	// Check for duplicate slug
	existing, _ := s.repo.GetBySlug(ctx, slug)
	if existing != nil {
		return nil, errors.New("category with this name already exists")
	}

	var category entity.Category
	utils.Copy(&category, req)
	category.Slug = slug

	if err := s.repo.Create(ctx, &category); err != nil {
		logger.Errorf("CategoryService.Create fail, error: %s", err)
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) Update(ctx context.Context, id string, req *dto.UpdateCategoryReq) (*entity.Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if req.Name != "" {
		category.Name = req.Name
		category.Slug = generateSlug(req.Name)
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Image != "" {
		category.Image = req.Image
	}

	if err := s.repo.Update(ctx, category); err != nil {
		logger.Errorf("CategoryService.Update fail, error: %s", err)
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("category not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *CategoryService) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) List(ctx context.Context) ([]entity.Category, error) {
	return s.repo.List(ctx)
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
