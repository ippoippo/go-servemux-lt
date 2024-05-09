package std

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/ippoippo/go-servemux-lt/demo"
	"github.com/ippoippo/go-servemux-lt/demo/internal/slogg"
)

func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Fake get notes

	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()

	note1 := demo.NoteFromCreateRequest(id1, demo.CreateNoteRequest{Text: "Some fake note 1"})
	note2 := demo.NoteFromCreateRequest(id2, demo.CreateNoteRequest{Text: "Some fake note 2"})
	notes := []*demo.Note{note1, note2}

	// Encode the response as JSON and send it
	writeResponse(ctx, w, notes, "get all notes")
}

func GetNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idPathValue := r.PathValue("id")

	if idPathValue == "unexpected" {
		// Fake something unexpected
		slog.ErrorContext(ctx, "faking something unexpected broke")
		http.Error(w, errors.New("something unexpected").Error(), http.StatusInternalServerError)
		return
	}

	id, err := uuid.Parse(idPathValue)
	if err != nil {
		msg := fmt.Sprintf("invalid id: [%v]", idPathValue)
		writeErrorResponse(w, http.StatusBadRequest, demo.ErrWithMsg(msg))
		return
	}

	if idPathValue != "018f7617-2706-7a87-afb9-ee40594dbab5" {
		// Fake not found
		writeErrorResponse(w, http.StatusNotFound, demo.ErrWithMsg("unable to get note"))
		return
	}

	// Fake looking up note etc
	nr := &demo.CreateNoteRequest{Text: "Some fake note 1"}
	note := demo.NoteFromCreateRequest(id, *nr)
	note.Id = id

	/// Encode the response as JSON and send it
	writeResponse(ctx, w, note, "get note")
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "creating note")
	req := &demo.CreateNoteRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		slog.ErrorContext(ctx, "unable to decode request body", slogg.ErrorAttr(err))
		writeErrorResponse(w, http.StatusBadRequest, demo.ErrWithMsg("unable to decode request body"))
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, demo.ErrWithMsg("unable to create note"))
		return
	}

	note := demo.NoteFromCreateRequest(id, *req)

	// Normally do stuff here

	// Encode the response as JSON and send it
	writeResponse(ctx, w, note, "create note")
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Would normally do something where

	slog.InfoContext(ctx, "deleted note")
	w.WriteHeader(http.StatusNoContent)
}
