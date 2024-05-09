package std

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ippoippo/go-servemux-lt/demo"
	"github.com/ippoippo/go-servemux-lt/demo/internal/slogg"
)

func writeErrorResponse(w http.ResponseWriter, code int, errResponse *demo.ErrorResponse) {
	encoder := json.NewEncoder(w)
	w.WriteHeader(code)
	if err := encoder.Encode(errResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeResponse(ctx context.Context, w http.ResponseWriter, result any, action string) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(result); err != nil {
		slog.ErrorContext(ctx, "unable to encode response", slogg.ErrorAttr(err))
		msg := fmt.Sprintf("unable to %s", action)
		writeErrorResponse(w, http.StatusInternalServerError, demo.ErrWithMsg(msg))
	}
}
