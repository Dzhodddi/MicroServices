package main

import (
	"auth/internal/service"
	"commons/shared_errors"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
	service *service.Service
	server  *echo.Echo
}

func (h *handler) mount() {
	v1 := h.server.Group("/v1")
	v1.GET("/health", h.healthCheckHandler)
	users := v1.Group("/users")
	users.POST("/activate", h.activateUserHandler)
	auth := v1.Group("/auth")
	auth.POST("/register", h.registerNewUser)
	auth.POST("/token", nil)
}

func (h *handler) healthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status":  "OK",
		"env":     env,
		"version": version,
	})
}

func (h *handler) activateUserHandler(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, shared_errors.BadRequestPayload.Error())
	}
	return nil
}

func (h *handler) registerNewUser(c echo.Context) error {
	user, err := h.service.RegisterNewUser(c)
	if err != nil {
		switch {
		case errors.Is(err, shared_errors.BadRequestPayload):
			return c.JSON(http.StatusBadRequest, err.Error())
		case errors.Is(err, shared_errors.ValidationError):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, shared_errors.ServerError):
			return c.JSON(http.StatusInternalServerError, err.Error())
		case errors.Is(err, shared_errors.ViolatePK):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated, user)
}
