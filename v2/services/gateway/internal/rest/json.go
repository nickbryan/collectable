package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type jsonHandler[Req, Res any] struct {
	logger *zap.Logger
	action Action[Req, Res]
}

func (h *jsonHandler[Req, Res]) setLogger(l *zap.Logger) {
	h.logger = l
}

func (h *jsonHandler[Req, Res]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := Request[Req]{Request: r, PathValues: mux.Vars(r)}

	if r.Body != http.NoBody {
		if err := json.NewDecoder(r.Body).Decode(&req.Data); err != nil {
			h.logger.Error("unable to decode request data as json", zap.Error(err))

			return
		}
	}

	resp, err := h.action(req)
	if err != nil {
		// TODO: Handle both error types
	}

	w.WriteHeader(resp.statusCode)

	if _, ok := any(resp).(NoBody); !ok {
		if err := json.NewEncoder(w).Encode(resp.data); err != nil {
			h.logger.Error("unable to encode response data as json", zap.Error(err))
		}
	}
}

func JSON[Req, Res any](action Action[Req, Res]) http.Handler {
	return &jsonHandler[Req, Res]{action: action}
}
