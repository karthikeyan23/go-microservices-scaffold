package entity

import (
	"context"
	"time"
)

// Entity represents an object within the Domain or Business layer.
type Entity struct {
	ID        string    `json:"entity_id"`
	Name      string    `json:"entity_name"`
	CreatedAt time.Time `json:"created_at"`
}

// Repository is a generic interface for DB operations
type Repository interface {
	// Create creates a new entity
	Create(ctx context.Context, entity *Entity) error
	// Get reads an entity by ID
	Get(ctx context.Context, id string) (*Entity, error)
	// GetAll Gets a list of Entities
	GetAll(ctx context.Context) ([]*Entity, error)
	// Update updates an entity
	Update(ctx context.Context, entity *Entity) error
	// Delete deletes an entity
	Delete(ctx context.Context, id string) error
}
