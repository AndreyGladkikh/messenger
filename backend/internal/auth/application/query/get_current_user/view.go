package get_current_user

import "github.com/google/uuid"

type AuthView struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
