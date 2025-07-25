package main

import (
	"auth/internal/service"
	"auth/shared"
	"commons"
	"commons/shared_errors"
	"errors"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

type handler struct {
	service *service.Service
	server  *echo.Echo
}

func (h *handler) mount() {
	v1 := h.server.Group("/v1")
	v1.GET("/health", h.healthCheckHandler)
	auth := v1.Group("/auth")
	auth.POST("/register", h.registerNewUser)
	auth.POST("/login", h.loginUser)
	auth.GET("/activate", h.activateUserHandler)
	v1.GET("/swagger/*", echoSwagger.WrapHandler)

}

// HealthCheckAPI  godoc
//
//	@Summary		Health check
//	@Description	Health check
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	map[string]any
//	@Router			/health [get]
func (h *handler) healthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status":  "OK",
		"env":     env,
		"version": version,
	})
}

// activate user handler
//
// @Summary activate new user
// @Description activate new user via email and html page
// @Tags auth
// @Accept json
// @Produce html
// @Param	token query	string true "Invitation token"
// @Success 201	{object} nil "User activated"
// @Failure 400 {object} error
// @Failure 422 {object} error
// @Failure 500 {object} error
// @Router /auth/activate [get]
func (h *handler) activateUserHandler(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, shared_errors.BadRequestPayload.Error())
	}
	err := h.service.IUserService.Activate(c.Request().Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, shared_errors.NotFoundError):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.HTML(http.StatusOK, "activated")
}

// register user handler
//
// @Summary Register new user
// @Description Register new user with payload as struct RegisterNewUser
// @Tags auth
// @Accept json
// @Produce json
// @Param	payload body shared.RegisterNewUser true "User credentials"
// @Success 201	{object} shared.User "User registered"
// @Failure 400 {object} error
// @Failure 422 {object} error
// @Failure 500 {object} error
// @Router /auth/register [post]
func (h *handler) registerNewUser(c echo.Context) error {
	var payload shared.RegisterNewUser
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, shared_errors.BadRequestPayload.Error())
	}

	if err := commons.Validate.Struct(&payload); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, shared_errors.BadRequestPayload.Error())
	}

	userPayload := shared.RegisterNewUser{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	user, err := h.service.IUserService.RegisterNewUser(c.Request().Context(), userPayload)
	if err != nil {
		switch {
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

// login user handler
//
// @Summary Login user
// @Description Login user with payload as struct RegisterNewUser
// @Tags auth
// @Accept json
// @Produce json
// @Param	payload body shared.LoginUser true "User credentials"
// @Success 201	{object} shared.User "User log in successfully"
// @Failure 422 {object} error
// @Failure 500 {object} error
// @Router /auth/login [post]
func (h *handler) loginUser(c echo.Context) error {
	var payload shared.LoginUser
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, shared_errors.BadRequestPayload.Error())
	}

	if err := commons.Validate.Struct(&payload); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, shared_errors.BadRequestPayload.Error())
	}
	user, token, err := h.service.IUserService.Login(c.Request().Context(), payload)
	if err != nil {
		switch {
		case errors.Is(err, shared_errors.ServerError):
			return c.JSON(http.StatusInternalServerError, err.Error())
		case errors.Is(err, shared_errors.NotFoundError):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+token)
	return c.JSON(http.StatusCreated, struct {
		Token string `json:"token"`
		User  *shared.User
	}{
		Token: token,
		User:  user,
	})

}
