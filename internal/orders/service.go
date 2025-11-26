package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	repo "github.com/wind1102/ecom/internal/adapters/postgresql/sqlc"
)

var (
	ErrorProductNotFound = errors.New("product not found")
	ErrorProductNoStock  = errors.New("product has not enough stock")
)

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
}

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	if tempOrder.CustomerId == 0 {
		return repo.Order{}, fmt.Errorf("CustomerId is required")

	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("Items is required")
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// create an order
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerId)

	if err != nil {
		return repo.Order{}, err
	}

	for _, item := range tempOrder.Items {
		product, err := s.repo.FindProductById(ctx, item.ProductId)
		if err != nil {
			return repo.Order{}, ErrorProductNotFound
		}
		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrorProductNoStock
		}

		_, err = qtx.CreateOrderItem(ctx, order.ID, item.ProductId, item.Quantity, product.PriceInCenters)
		if err != nil {
			return repo.Order{}, err
		}

	}
	tx.Commit(ctx)
	return order, nil
}
