package repository

import (
	"context"

	"ecom/internal/cart/entity"
	"ecom/pkg/dbs"
)

type CartRepo struct {
	db *dbs.Database
}

func NewCartRepository(db *dbs.Database) *CartRepo {
	return &CartRepo{db: db}
}

func (r *CartRepo) GetOrCreateCart(ctx context.Context, userID string) (*entity.Cart, error) {
	var cart entity.Cart
	query := dbs.NewQuery("user_id = ?", userID)
	err := r.db.FindOne(ctx, &cart,
		dbs.WithQuery(query),
		dbs.WithPreload([]string{"Items", "Items.Product", "Items.Variant"}),
	)
	if err != nil {
		// Create new cart
		cart = entity.Cart{UserID: userID}
		if createErr := r.db.Create(ctx, &cart); createErr != nil {
			return nil, createErr
		}
		cart.Items = []entity.CartItem{}
	}
	return &cart, nil
}

func (r *CartRepo) GetCart(ctx context.Context, userID string) (*entity.Cart, error) {
	var cart entity.Cart
	query := dbs.NewQuery("user_id = ?", userID)
	err := r.db.FindOne(ctx, &cart,
		dbs.WithQuery(query),
		dbs.WithPreload([]string{"Items", "Items.Product", "Items.Variant"}),
	)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepo) AddItem(ctx context.Context, item *entity.CartItem) error {
	return r.db.Create(ctx, item)
}

func (r *CartRepo) GetCartItem(ctx context.Context, itemID string) (*entity.CartItem, error) {
	var item entity.CartItem
	if err := r.db.FindOne(ctx, &item, dbs.WithQuery(dbs.NewQuery("id = ?", itemID))); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CartRepo) GetCartItemByProductVariant(ctx context.Context, cartID, productID, variantID string) (*entity.CartItem, error) {
	var item entity.CartItem
	err := r.db.FindOne(ctx, &item,
		dbs.WithQuery(
			dbs.NewQuery("cart_id = ? AND product_id = ? AND variant_id = ?", cartID, productID, variantID),
		),
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CartRepo) UpdateItem(ctx context.Context, item *entity.CartItem) error {
	return r.db.Update(ctx, item)
}

func (r *CartRepo) RemoveItem(ctx context.Context, itemID string) error {
	return r.db.Delete(ctx, &entity.CartItem{}, dbs.WithQuery(dbs.NewQuery("id = ?", itemID)))
}

func (r *CartRepo) ClearCart(ctx context.Context, cartID string) error {
	return r.db.Delete(ctx, &entity.CartItem{}, dbs.WithQuery(dbs.NewQuery("cart_id = ?", cartID)))
}
