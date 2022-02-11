package model

import (
	"context"
	"errors"
)

var (
	ErrEntityNotFound = errors.New("entity not found")
)

// EntityService Interface contains the methods that are exposed to the transport layer.
type EntityService interface {
	GetEntity(ctx context.Context, id string) (*Entity, error)
	CreateEntity(ctx context.Context, entity *Entity) error
	UpdateEntity(ctx context.Context, user *Entity) error
	DeleteEntity(ctx context.Context, id string) error
}
