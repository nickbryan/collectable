package token

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nickbryan/collectable/services/gateway/internal/rest"
	"github.com/nickbryan/collectable/services/iam/token"
)

func NewCreateHandler(client token.TokenServiceClient, logger *zap.Logger) rest.Handler {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Token string `json:"token"`
	}

	return rest.Handler{
		Route: func(r *mux.Route) {
			r.Path("/auth/token").Methods(http.MethodPost)
		},
		Action: func(res rest.Responder, req *rest.Request) {
			var request request

			if err := req.Decode(&request); err != nil {
				res.Respond(http.StatusBadRequest).WithErrors(err)

				return
			}

			resp, err := client.CreateToken(req.Context(), &token.CreateTokenRequest{
				Email:    request.Email,
				Password: request.Password,
			})

			if err != nil {
				logger.Error("err from grpc client when calling token.CreateToken", zap.Error(err))
				res.Respond(http.StatusInternalServerError)

				return
			}

			res.Respond(http.StatusCreated).WithData(response{
				Token: resp.Token,
			})
		},
	}
}
