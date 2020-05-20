package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func (h *Handlers) ZoomIn(c echo.Context) error {
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

	if err := cam.ZoomTele(ctx); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) ZoomOut(c echo.Context) error {
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

	if err := cam.ZoomWide(ctx); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) ZoomStop(c echo.Context) error {
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

	if err := cam.ZoomStop(ctx); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
