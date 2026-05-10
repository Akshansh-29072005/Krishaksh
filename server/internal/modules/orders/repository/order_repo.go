package repository

import (
	"context"
	"encoding/json"

	"github.com/aarcsx/krishaksh-backend/internal/database"
	"github.com/aarcsx/krishaksh-backend/internal/models"
	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrderWithItems(ctx context.Context, order *models.Order, items []*models.OrderItem) error
	GetOrdersByUser(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
	GetOrderByID(ctx context.Context, orderID, userID uuid.UUID) (*models.Order, error)
	GetOrderByIDUnscoped(ctx context.Context, orderID uuid.UUID) (*models.Order, error)
	SetOrderStatus(ctx context.Context, orderID uuid.UUID, from []string, to string) (bool, error)
	LinkPayment(ctx context.Context, orderID, paymentID uuid.UUID) error
}

type orderRepoImpl struct{ db *database.DB }

func NewOrderRepository(db *database.DB) OrderRepository { return &orderRepoImpl{db: db} }

func (r *orderRepoImpl) CreateOrderWithItems(ctx context.Context, order *models.Order, items []*models.OrderItem) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	q := `INSERT INTO orders (id, user_id, status, currency, subtotal, tax_amount, shipping_amount, discount_amount, grand_total, shipping_metadata, notes, created_at, updated_at)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW(),NOW())`
	_, err = tx.Exec(ctx, q, order.ID, order.UserID, order.Status, order.Currency, order.Subtotal, order.TaxAmount, order.ShippingAmount, order.DiscountAmount, order.GrandTotal, order.ShippingMetadata, order.Notes)
	if err != nil {
		return err
	}
	for _, it := range items {
		_, err = tx.Exec(ctx, `INSERT INTO order_items (id, order_id, product_id, quantity, unit_price, line_total, created_at) VALUES ($1,$2,$3,$4,$5,$6,NOW())`, it.ID, it.OrderID, it.ProductID, it.Quantity, it.UnitPrice, it.LineTotal)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *orderRepoImpl) GetOrdersByUser(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT id, user_id, status, currency, subtotal, tax_amount, shipping_amount, discount_amount, grand_total, shipping_metadata, payment_id, notes, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Order
	for rows.Next() {
		o := &models.Order{}
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.Currency, &o.Subtotal, &o.TaxAmount, &o.ShippingAmount, &o.DiscountAmount, &o.GrandTotal, &o.ShippingMetadata, &o.PaymentID, &o.Notes, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, nil
}

func (r *orderRepoImpl) GetOrderByID(ctx context.Context, orderID, userID uuid.UUID) (*models.Order, error) {
	o := &models.Order{}
	err := r.db.Pool.QueryRow(ctx, `SELECT id, user_id, status, currency, subtotal, tax_amount, shipping_amount, discount_amount, grand_total, shipping_metadata, payment_id, notes, created_at, updated_at FROM orders WHERE id = $1 AND user_id = $2`, orderID, userID).Scan(&o.ID, &o.UserID, &o.Status, &o.Currency, &o.Subtotal, &o.TaxAmount, &o.ShippingAmount, &o.DiscountAmount, &o.GrandTotal, &o.ShippingMetadata, &o.PaymentID, &o.Notes, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Pool.Query(ctx, `SELECT id, order_id, product_id, quantity, unit_price, line_total, created_at FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		it := &models.OrderItem{}
		if err := rows.Scan(&it.ID, &it.OrderID, &it.ProductID, &it.Quantity, &it.UnitPrice, &it.LineTotal, &it.CreatedAt); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, it)
	}
	return o, nil
}

func (r *orderRepoImpl) GetOrderByIDUnscoped(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	o := &models.Order{}
	err := r.db.Pool.QueryRow(ctx, `SELECT id, user_id, status, currency, subtotal, tax_amount, shipping_amount, discount_amount, grand_total, shipping_metadata, payment_id, notes, created_at, updated_at FROM orders WHERE id = $1`, orderID).Scan(&o.ID, &o.UserID, &o.Status, &o.Currency, &o.Subtotal, &o.TaxAmount, &o.ShippingAmount, &o.DiscountAmount, &o.GrandTotal, &o.ShippingMetadata, &o.PaymentID, &o.Notes, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (r *orderRepoImpl) SetOrderStatus(ctx context.Context, orderID uuid.UUID, from []string, to string) (bool, error) {
	_, err := r.db.Pool.Exec(ctx, `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2 AND status = ANY($3::text[])`, to, orderID, from)
	if err != nil {
		return false, err
	}
	ct, err := r.db.Pool.Query(ctx, `SELECT 1 FROM orders WHERE id = $1 AND status = $2`, orderID, to)
	if err != nil {
		return false, err
	}
	defer ct.Close()
	return ct.Next(), nil
}

func (r *orderRepoImpl) LinkPayment(ctx context.Context, orderID, paymentID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE orders SET payment_id = $1, updated_at = NOW() WHERE id = $2`, paymentID, orderID)
	return err
}

var _ = json.RawMessage{}
