package repository

import (
	"context"

	"github.com/aarcsx/krisho-backend/internal/database"
	"github.com/aarcsx/krisho-backend/internal/models"
	"github.com/google/uuid"
)

type CartRepository interface {
	GetOrCreateCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	GetCartItems(ctx context.Context, cartID uuid.UUID) ([]*models.CartItem, error)
	UpsertCartItem(ctx context.Context, cartID, productID uuid.UUID, qty int, unitPrice float64) error
	UpdateCartItemQuantity(ctx context.Context, cartID, itemID uuid.UUID, qty int, unitPrice float64) error
	DeleteCartItem(ctx context.Context, cartID, itemID uuid.UUID) error
	ClearCart(ctx context.Context, cartID uuid.UUID) error
}

type cartRepoImpl struct{ db *database.DB }

func NewCartRepository(db *database.DB) CartRepository { return &cartRepoImpl{db: db} }

func (r *cartRepoImpl) GetOrCreateCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	q := `INSERT INTO carts (user_id) VALUES ($1) ON CONFLICT (user_id) DO UPDATE SET updated_at = NOW() RETURNING id, user_id, created_at, updated_at`
	c := &models.Cart{}
	err := r.db.Pool.QueryRow(ctx, q, userID).Scan(&c.ID, &c.UserID, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *cartRepoImpl) GetCartItems(ctx context.Context, cartID uuid.UUID) ([]*models.CartItem, error) {
	q := `SELECT id, cart_id, product_id, quantity, unit_price, line_total, created_at, updated_at FROM cart_items WHERE cart_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, q, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*models.CartItem
	for rows.Next() {
		i := &models.CartItem{}
		if err := rows.Scan(&i.ID, &i.CartID, &i.ProductID, &i.Quantity, &i.UnitPrice, &i.LineTotal, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (r *cartRepoImpl) UpsertCartItem(ctx context.Context, cartID, productID uuid.UUID, qty int, unitPrice float64) error {
	q := `INSERT INTO cart_items (cart_id, product_id, quantity, unit_price, line_total)
	      VALUES ($1,$2,$3,$4,$5)
	      ON CONFLICT (cart_id, product_id) DO UPDATE
	      SET quantity = EXCLUDED.quantity, unit_price = EXCLUDED.unit_price, line_total = EXCLUDED.line_total, updated_at = NOW()`
	_, err := r.db.Pool.Exec(ctx, q, cartID, productID, qty, unitPrice, unitPrice*float64(qty))
	return err
}

func (r *cartRepoImpl) DeleteCartItem(ctx context.Context, cartID, itemID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM cart_items WHERE id = $1 AND cart_id = $2`, itemID, cartID)
	return err
}

func (r *cartRepoImpl) UpdateCartItemQuantity(ctx context.Context, cartID, itemID uuid.UUID, qty int, unitPrice float64) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE cart_items SET quantity=$1, unit_price=$2, line_total=$3, updated_at=NOW() WHERE cart_id=$4 AND id=$5`, qty, unitPrice, unitPrice*float64(qty), cartID, itemID)
	return err
}

func (r *cartRepoImpl) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM cart_items WHERE cart_id = $1`, cartID)
	return err
}
