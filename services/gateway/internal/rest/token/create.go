package token

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nickbryan/collectable/proto/iam/token/service/v1"
	"github.com/nickbryan/collectable/services/gateway/internal/rest"
)

func CreateHandler(client token.TokenServiceClient, logger *zap.Logger) rest.Handler {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Token string `json:"token"`
	}

	return rest.Handler{
		Route: func(r *mux.Route) {
			r.Path("/api/auth/token").Methods(http.MethodPost)
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

			st, ok := status.FromError(err)
			if !ok {
				logger.Error("err from grpc client when calling token.CreateToken", zap.Error(err))
				res.Respond(http.StatusInternalServerError)

				return
			}

			switch st.Code() {
			case codes.OK:
				res.Respond(http.StatusCreated).WithData(response{Token: resp.Token})
			case codes.NotFound:
				res.Respond(http.StatusNotFound)
			default:
				logger.Error("unexpected status code from grpc se when calling token.CreateToken", zap.Error(err))
				res.Respond(http.StatusInternalServerError)
			}
		},
	}
}
