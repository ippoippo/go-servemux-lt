package echo

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/ippoippo/go-servemux-lt/demo"
	"github.com/labstack/echo/v4"
)

func GetAllNotes(c echo.Context) error {
	// Fake get notes

	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()

	note1 := demo.NoteFromCreateRequest(id1, demo.CreateNoteRequest{Text: "Some fake note 1"})
	note2 := demo.NoteFromCreateRequest(id2, demo.CreateNoteRequest{Text: "Some fake note 2"})
	notes := []*demo.Note{note1, note2}

	return c.JSON(http.StatusOK, notes)
}

func GetNote(c echo.Context) error {
	ctx := c.Request().Context()
	idParam := c.Param("id")

	if idParam == "unexpected" {
		// Fake something unexpected
		slog.ErrorContext(ctx, "faking something unexpected broke")
		return errors.New("something unexpected")
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		msg := fmt.Sprintf("invalid id: [%v]", idParam)
		return c.JSON(http.StatusBadRequest, demo.ErrWithMsg(msg))
	}

	if idParam != "018f7617-2706-7a87-afb9-ee40594dbab5" {
		// Fake not found
		return c.JSON(http.StatusNotFound, demo.ErrWithMsg("unable to get note"))
	}

	// Fake looking up note etc
	nr := &demo.CreateNoteRequest{Text: "Some fake note 1"}
	note := demo.NoteFromCreateRequest(id, *nr)
	note.Id = id

	return c.JSON(http.StatusOK, note)
}

func CreateNote(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "creating note")

	req := &demo.CreateNoteRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, demo.ErrWithMsg("unable to parse request"))
	}

	id, err := uuid.NewV7()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, demo.ErrWithMsg("unable to create note"))
	}

	note := demo.NoteFromCreateRequest(id, *req)

	// Normally do stuff here

	return c.JSON(http.StatusCreated, note)
}

func DeleteNote(c echo.Context) error {
	ctx := c.Request().Context()
	// Would normally do something where

	slog.InfoContext(ctx, "deleted note")

	return c.NoContent(http.StatusNoContent)
}
