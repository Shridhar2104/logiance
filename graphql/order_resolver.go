package main

import (
	"context"
	"github.com/Shridhar2104/logilo/graphql/models"
)

// Add this new resolver struct and implementation
type orderResolver struct {
    server *Server
}

func (r *orderResolver) LineItems(ctx context.Context, obj *models.Order) ([]models.OrderLineItem, error) {
    // Implement line items resolution for an order
    return obj.LineItems, nil
}

func(r *orderResolver) CreatedAt(ctx context.Context, obj *models.Order) (string, error){
	return "", nil
}
