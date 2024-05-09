package demo

import (
	"time"

	"github.com/google/uuid"
)

type CreateNoteRequest struct {
	Text string `json:"text"`
}

type Note struct {
	Id        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	CreatedAt string    `json:"created_at"`
}

func NoteFromCreateRequest(id uuid.UUID, n CreateNoteRequest) *Note {
	return &Note{
		Id:        id,
		Text:      n.Text,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}
