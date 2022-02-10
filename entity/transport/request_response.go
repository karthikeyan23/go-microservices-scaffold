package transport

import "time"

type GetEntityByIDRequest struct {
	ID string `json:"entity_id"`
}

type GetEntityByIDResponse struct {
	ID        string    `json:"entity_id"`
	Name      string    `json:"entity_name"`
	CreatedAt time.Time `json:"created_at"`
}
