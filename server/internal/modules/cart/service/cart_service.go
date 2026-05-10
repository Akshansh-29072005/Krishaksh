package service

import (
	"context"
	"errors"

	cartRepo "github.com/aarcsx/krisho-backend/internal/modules/cart/repository"
	productRepo "github.com/aarcsx/krisho-backend/internal/modules/products/repository"
	"github.com/google/uuid"
)

type CartService interface {
	AddItem(ctx context.Context, userID, productID uuid.UUID, qty int) error
	GetCart(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	UpdateItemQuantity(ctx context.Context, userID, itemID, productID uuid.UUID, qty int) error
	RemoveItem(ctx context.Context, userID, itemID uuid.UUID) error
}

type cartServiceImpl struct {
	cartRepo    cartRepo.CartRepository
	productRepo productRepo.ProductRepository
}

func NewCartService(c cartRepo.CartRepository, p productRepo.ProductRepository) CartService {
	return &cartServiceImpl{cartRepo: c, productRepo: p}
}

func (s *cartServiceImpl) AddItem(ctx context.Context, userID, productID uuid.UUID, qty int) error {
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if !p.IsActive || !p.StockAvailable || p.StockQuantity < qty {
		return errors.New("product unavailable or insufficient stock")
	}
	cart, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}
	return s.cartRepo.UpsertCartItem(ctx, cart.ID, productID, qty, p.Price)
}

func (s *cartServiceImpl) GetCart(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	cart, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	items, err := s.cartRepo.GetCartItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}
	subtotal := 0.0
	for _, i := range items {
		subtotal += i.LineTotal
	}
	return map[string]interface{}{"cart_id": cart.ID, "items": items, "subtotal": subtotal, "currency": "INR"}, nil
}

func (s *cartServiceImpl) RemoveItem(ctx context.Context, userID, itemID uuid.UUID) error {
	cart, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}
	return s.cartRepo.DeleteCartItem(ctx, cart.ID, itemID)
}

func (s *cartServiceImpl) UpdateItemQuantity(ctx context.Context, userID, itemID, productID uuid.UUID, qty int) error {
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if !p.IsActive || !p.StockAvailable || p.StockQuantity < qty {
		return errors.New("product unavailable or insufficient stock")
	}
	cart, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}
	return s.cartRepo.UpdateCartItemQuantity(ctx, cart.ID, itemID, qty, p.Price)
}
