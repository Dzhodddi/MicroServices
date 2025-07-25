package main

import (
	"commons"
	pb "commons/api"
	"commons/shared_errors"
	"commons/shared_types"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type Handler struct {
	authClient pb.AuthServiceClient
}
type TokenCheck struct {
	Email string `json:"email" validate:"required,email"`
}

func NewHandler(client pb.AuthServiceClient) *Handler {
	return &Handler{client}
}

func (h *Handler) ValidateToken(c echo.Context) error {
	var payload TokenCheck
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, shared_errors.BadRequestPayload.Error())
	}

	if err := commons.Validate.Struct(&payload); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, shared_errors.BadRequestPayload.Error())
	}

	user, err := h.authClient.ValidateToken(c.Request().Context(), &pb.TokenRequest{
		Email: payload.Email,
	})
	if err != nil {
		switch {
		case errors.Is(err, status.Error(codes.NotFound, "not found")):
			return c.JSON(http.StatusUnauthorized, shared_errors.NotFoundError.Error())
		default:

			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("user: %v, err %v\n", user, err))
		}
	}
	return c.JSON(http.StatusOK, &shared_types.RedisUserInfo{
		Email: user.Email,
		TTL:   user.Ttl,
	})
}
