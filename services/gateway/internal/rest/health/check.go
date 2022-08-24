package health

import (
	"github.com/gorilla/mux"
	"github.com/nickbryan/collectable/services/gateway/internal/rest"
	"net/http"
)

func CheckHandler() rest.Handler {
	return rest.Handler{
		Route: func(r *mux.Route) {
			r.Path("/api/health").Methods(http.MethodGet)
		},
		Action: func(res rest.Responder, _ *rest.Request) {
			res.Respond(http.StatusOK)
		},
	}
}
