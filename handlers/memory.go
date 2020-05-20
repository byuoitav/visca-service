package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func (h *Handlers) MemoryRecall(c echo.Context) error {
	addr := c.Param("address")
	channel, err := strconv.Atoi(c.Param("channel"))
	switch {
	case len(addr) == 0:
		return c.String(http.StatusBadRequest, "must include the address of the camera")
	case err != nil:
		return c.String(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	cam, err := h.CreateCamera(ctx, addr)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to create camera: %s", err))
	}

	if err := cam.MemoryRecall(ctx, byte(channel)); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
