package handlers

import (
	"context"

	"github.com/byuoitav/visca"
	"github.com/labstack/echo"
)

func (h *Handlers) ZoomIn(c echo.Context) error {
	return h.generic(c, "ZoomIn", func(ctx context.Context, cam *visca.Camera) error {
		return cam.ZoomTele(ctx)
	})
}

func (h *Handlers) ZoomOut(c echo.Context) error {
	return h.generic(c, "ZoomOut", func(ctx context.Context, cam *visca.Camera) error {
		return cam.ZoomWide(ctx)
	})
}

func (h *Handlers) ZoomStop(c echo.Context) error {
	return h.generic(c, "ZoomStop", func(ctx context.Context, cam *visca.Camera) error {
		return cam.ZoomStop(ctx)
	})
}
