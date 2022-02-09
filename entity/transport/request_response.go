package transport

type GetEntityByIDRequest struct {
	ID string `json:"entity_id"`
}

type GetEntityByIDResponse struct {
	ID        string `json:"entity_id"`
	Name      string `json:"entity_name"`
	CreatedAt int64  `json:"created_at"`
}
