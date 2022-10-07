package identity

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	identity "github.com/nickbryan/collectable/proto/iam/identity/service/v1"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nickbryan/collectable/services/gateway/internal/rest"
)

func CreateHandler(client identity.IdentityServiceClient, logger *zap.Logger) rest.Handler {
	type request struct {
		Email                string `json:"email"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"passwordConfirmation"`
	}

	type response struct {
		ID string `json:"id"`
	}

	return rest.Handler{
		Route: func(r *mux.Route) {
			r.Path("/api/auth/identity").Methods(http.MethodPost)
		},
		Action: func(res rest.Responder, req *rest.Request) {
			var request request

			if err := req.Decode(&request); err != nil {
				res.Respond(http.StatusBadRequest).WithErrors(err)

				return
			}

			if err := validation.ValidateStruct(
				&request,
				validation.Field(&request.Email, validation.Required, is.Email),
				validation.Field(&request.Password, validation.Required, validation.Length(8, 256)),
				validation.Field(&request.PasswordConfirmation, validation.Required, validation.Length(8, 256), validation.By(func(_ interface{}) error {
					if request.Password != request.PasswordConfirmation {
						return errors.New("password does mot match password confirmation")
					}

					return nil
				})),
			); err != nil {
				res.Respond(http.StatusBadRequest).WithErrors(err)

				return
			}

			resp, err := client.CreateIdentity(req.Context(), &identity.CreateIdentityRequest{
				Email:                request.Email,
				Password:             request.Password,
				PasswordConfirmation: request.PasswordConfirmation,
			})

			st, ok := status.FromError(err)
			if !ok {
				logger.Error("err from grpc client when calling identity.CreateIdentityRequest", zap.Error(err))
				res.Respond(http.StatusInternalServerError)

				return
			}

			switch st.Code() {
			case codes.OK:
				res.Respond(http.StatusCreated).WithData(response{ID: resp.Id})
			default:
				logger.Error("unexpected status code from grpc se when calling identity.CreateIdentityRequest", zap.Error(err))
				res.Respond(http.StatusInternalServerError)
			}
		},
	}
}
