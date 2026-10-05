package get_current_user

import "uuid"

type AuthView struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
