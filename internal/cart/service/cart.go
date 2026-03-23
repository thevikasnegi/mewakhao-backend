package service

import (
	"context"
	"errors"

	"ecom/internal/cart/dto"
	"ecom/internal/cart/entity"
	"ecom/internal/cart/repository"

	"github.com/quangdangfit/gocommon/logger"
)

const (
	shippingThreshold = 50.0
	shippingCost      = 5.99
)

type CartService struct {
	repo *repository.CartRepo
}

func NewCartService(repo *repository.CartRepo) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) GetCart(ctx context.Context, userID string) (*dto.CartRes, error) {
	cart, err := s.repo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.mapCartToRes(cart), nil
}

func (s *CartService) AddItem(ctx context.Context, userID string, req *dto.AddItemReq) (*dto.CartRes, error) {
	cart, err := s.repo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if item already exists in cart
	existingItem, _ := s.repo.GetCartItemByProductVariant(ctx, cart.ID, req.ProductID, req.VariantID)
	if existingItem != nil {
		existingItem.Quantity += req.Quantity
		if err := s.repo.UpdateItem(ctx, existingItem); err != nil {
			logger.Errorf("CartService.AddItem update fail, error: %s", err)
			return nil, err
		}
	} else {
		item := &entity.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			VariantID: req.VariantID,
			Quantity:  req.Quantity,
		}
		if err := s.repo.AddItem(ctx, item); err != nil {
			logger.Errorf("CartService.AddItem create fail, error: %s", err)
			return nil, err
		}
	}

	// Re-fetch cart with all relations
	updatedCart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.mapCartToRes(updatedCart), nil
}

func (s *CartService) UpdateItem(ctx context.Context, userID, itemID string, req *dto.UpdateItemReq) (*dto.CartRes, error) {
	cart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	item, err := s.repo.GetCartItem(ctx, itemID)
	if err != nil || item.CartID != cart.ID {
		return nil, errors.New("item not found in cart")
	}

	item.Quantity = req.Quantity
	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return nil, err
	}

	updatedCart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.mapCartToRes(updatedCart), nil
}

func (s *CartService) RemoveItem(ctx context.Context, userID, itemID string) (*dto.CartRes, error) {
	cart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	item, err := s.repo.GetCartItem(ctx, itemID)
	if err != nil || item.CartID != cart.ID {
		return nil, errors.New("item not found in cart")
	}

	if err := s.repo.RemoveItem(ctx, itemID); err != nil {
		return nil, err
	}

	updatedCart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.mapCartToRes(updatedCart), nil
}

func (s *CartService) ClearCart(ctx context.Context, userID string) error {
	cart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return errors.New("cart not found")
	}
	return s.repo.ClearCart(ctx, cart.ID)
}

func (s *CartService) mapCartToRes(cart *entity.Cart) *dto.CartRes {
	res := &dto.CartRes{
		ID:    cart.ID,
		Items: []dto.CartItemRes{},
	}

	subtotal := 0.0
	for _, item := range cart.Items {
		itemRes := dto.CartItemRes{
			ID:       item.ID,
			Quantity: item.Quantity,
		}

		if item.Product != nil {
			itemRes.Product = dto.ProductRef{
				ID:               item.Product.ID,
				Name:             item.Product.Name,
				Slug:             item.Product.Slug,
				ShortDescription: item.Product.ShortDescription,
				Images:           item.Product.Images,
				BasePrice:        item.Product.BasePrice,
			}
			if itemRes.Product.Images == nil {
				itemRes.Product.Images = []string{}
			}
		}

		if item.Variant != nil {
			itemRes.Variant = dto.VariantRef{
				ID:     item.Variant.ID,
				Weight: item.Variant.Weight,
				Price:  item.Variant.Price,
				Stock:  item.Variant.Stock,
			}
			subtotal += item.Variant.Price * float64(item.Quantity)
		}

		res.Items = append(res.Items, itemRes)
	}

	res.Subtotal = subtotal
	if subtotal >= shippingThreshold {
		res.Shipping = 0
	} else {
		res.Shipping = shippingCost
	}
	res.Total = subtotal + res.Shipping

	return res
}
