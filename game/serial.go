package game

import (
	"github.com/google/uuid"
)

type Serial struct {
	ID        uuid.UUID `json:"id"`
	ClassName string    `json:"class_name"`
	Action    string    `json:"action"`
}
