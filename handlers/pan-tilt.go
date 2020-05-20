package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func (h *Handlers) TiltUp(c echo.Context) error {
	addr := c.Param("address")
	if len(addr) == 0 {
		return c.String(http.StatusBadRequest, "must include the address of the camera")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	cam, err := h.CreateCamera(ctx, addr)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to create camera: %s", err))
	}

	if err := cam.TiltUp(ctx, 0x0e); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) TiltDown(c echo.Context) error {
	addr := c.Param("address")
	if len(addr) == 0 {
		return c.String(http.StatusBadRequest, "must include the address of the camera")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	cam, err := h.CreateCamera(ctx, addr)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to create camera: %s", err))
	}

	if err := cam.TiltDown(ctx, 0x0e); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) PanLeft(c echo.Context) error {
	addr := c.Param("address")
	if len(addr) == 0 {
		return c.String(http.StatusBadRequest, "must include the address of the camera")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	cam, err := h.CreateCamera(ctx, addr)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to create camera: %s", err))
	}

	if err := cam.PanLeft(ctx, 0x0b); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) PanRight(c echo.Context) error {
	addr := c.Param("address")
	if len(addr) == 0 {
		return c.String(http.StatusBadRequest, "must include the address of the camera")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	cam, err := h.CreateCamera(ctx, addr)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to create camera: %s", err))
	}

	if err := cam.PanRight(ctx, 0x0b); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) PanTiltStop(c echo.Context) error {
	addr := c.Param("address")
	if len(addr) == 0 {
		return c.String(http.StatusBadRequest, "must include the address of the camera")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	cam, err := h.CreateCamera(ctx, addr)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to create camera: %s", err))
	}

	if err := cam.PanTiltStop(ctx); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
