package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aarcsx/krishaksh-backend/internal/models"
	cartRepo "github.com/aarcsx/krishaksh-backend/internal/modules/cart/repository"
	orderRepo "github.com/aarcsx/krishaksh-backend/internal/modules/orders/repository"
	productRepo "github.com/aarcsx/krishaksh-backend/internal/modules/products/repository"
	"github.com/google/uuid"
)

var AllowedTransitions = map[string][]string{
	"pending":         {"payment_pending", "cancelled", "failed"},
	"payment_pending": {"paid", "failed", "cancelled"},
	"paid":            {"processing", "refunded"},
	"processing":      {"shipped", "cancelled"},
	"shipped":         {"delivered", "refunded"},
	"delivered":       {"refunded"},
	"cancelled":       {}, "refunded": {}, "failed": {},
}

type OrderService interface {
	CreateFromCart(ctx context.Context, userID uuid.UUID, shipping map[string]interface{}, notes *string) (*models.Order, error)
	GetOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
	GetOrder(ctx context.Context, userID, orderID uuid.UUID) (*models.Order, error)
	Transition(ctx context.Context, orderID uuid.UUID, from, to string) error
}

type orderServiceImpl struct {
	orders   orderRepo.OrderRepository
	carts    cartRepo.CartRepository
	products productRepo.ProductRepository
}

func NewOrderService(o orderRepo.OrderRepository, c cartRepo.CartRepository, p productRepo.ProductRepository) OrderService {
	return &orderServiceImpl{orders: o, carts: c, products: p}
}

func (s *orderServiceImpl) CreateFromCart(ctx context.Context, userID uuid.UUID, shipping map[string]interface{}, notes *string) (*models.Order, error) {
	cart, err := s.carts.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	cartItems, err := s.carts.GetCartItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}
	order := &models.Order{ID: uuid.New(), UserID: userID, Status: "payment_pending", Currency: "INR", ShippingAmount: 0, TaxAmount: 0, DiscountAmount: 0, Notes: notes}
	if shipping == nil {
		shipping = map[string]interface{}{}
	}
	meta, _ := json.Marshal(shipping)
	order.ShippingMetadata = meta
	var items []*models.OrderItem
	for _, ci := range cartItems {
		p, err := s.products.GetByID(ctx, ci.ProductID)
		if err != nil {
			return nil, err
		}
		if !p.IsActive || !p.StockAvailable || p.StockQuantity < ci.Quantity {
			return nil, errors.New("stock changed; revalidate cart")
		}
		it := &models.OrderItem{ID: uuid.New(), OrderID: order.ID, ProductID: p.ID, Quantity: ci.Quantity, UnitPrice: p.Price, LineTotal: p.Price * float64(ci.Quantity)}
		items = append(items, it)
		order.Subtotal += it.LineTotal
	}
	order.GrandTotal = order.Subtotal + order.TaxAmount + order.ShippingAmount - order.DiscountAmount
	if err := s.orders.CreateOrderWithItems(ctx, order, items); err != nil {
		return nil, err
	}
	_ = s.carts.ClearCart(ctx, cart.ID)
	return order, nil
}
func (s *orderServiceImpl) GetOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	return s.orders.GetOrdersByUser(ctx, userID)
}
func (s *orderServiceImpl) GetOrder(ctx context.Context, userID, orderID uuid.UUID) (*models.Order, error) {
	return s.orders.GetOrderByID(ctx, orderID, userID)
}
func (s *orderServiceImpl) Transition(ctx context.Context, orderID uuid.UUID, from, to string) error {
	allowed := false
	for _, n := range AllowedTransitions[from] {
		if n == to {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.New("invalid order state transition")
	}
	_, err := s.orders.SetOrderStatus(ctx, orderID, []string{from}, to)
	return err
}
