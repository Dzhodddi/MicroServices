package main

import (
	pb "commons/api"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	authClient pb.AuthServiceClient
}

func NewHandler(client pb.AuthServiceClient) *Handler {
	return &Handler{client}
}

func (h *Handler) ValidateToken(c echo.Context) error {
	tokenPath := c.Param("token")
	token, err := h.authClient.ValidateToken(c.Request().Context(), &pb.Token{Token: tokenPath})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, token.Expired)
}
