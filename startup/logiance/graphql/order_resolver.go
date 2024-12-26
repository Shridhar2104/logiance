package main

import (
    "context"
    "time"
    "github.com/Shridhar2104/logilo/graphql/models"
)

type orderResolver struct {
    server *Server
}

func (r *orderResolver) CancelledAt(ctx context.Context, obj *models.Order) (*string, error) {
    if obj.CancelledAt == nil {
        return nil, nil
    }
    str := obj.CancelledAt.Format(time.RFC3339)
    return &str, nil
}

func (r *orderResolver) ClosedAt(ctx context.Context, obj *models.Order) (*string, error) {
    if obj.ClosedAt == nil {
        return nil, nil
    }
    str := obj.ClosedAt.Format(time.RFC3339)
    return &str, nil
}

func (r *orderResolver) CreatedAt(ctx context.Context, obj *models.Order) (string, error) {
    return obj.CreatedAt.Format(time.RFC3339), nil
}

func (r *orderResolver) LineItems(ctx context.Context, obj *models.Order) ([]*models.OrderLineItem, error) {
    return obj.LineItems, nil
}

func (r *orderResolver) ProcessedAt(ctx context.Context, obj *models.Order) (*string, error) {
    if obj.ProcessedAt == nil {
        return nil, nil
    }
    str := obj.ProcessedAt.Format(time.RFC3339)
    return &str, nil
}

func (r *orderResolver) UpdatedAt(ctx context.Context, obj *models.Order) (string, error) {
    return obj.UpdatedAt.Format(time.RFC3339), nil
}

func (r *orderResolver) Customer(ctx context.Context, obj *models.Order) (*Customer, error) {
    if obj.Customer == nil {
        return nil, nil
    }
    return &Customer{
        ID:        &obj.Customer.ID,
        Email:     &obj.Customer.Email,
        FirstName: &obj.Customer.FirstName,
        LastName:  &obj.Customer.LastName,
        Phone:     &obj.Customer.Phone,
    }, nil
}