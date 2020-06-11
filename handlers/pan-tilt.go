package handlers

import (
	"context"

	"github.com/byuoitav/visca"
	"github.com/labstack/echo"
)

func (h *Handlers) TiltUp(c echo.Context) error {
	return h.generic(c, "TiltUp", func(ctx context.Context, cam *visca.Camera) error {
		return cam.TiltUp(ctx, 0x0e)
	})
}

func (h *Handlers) TiltDown(c echo.Context) error {
	return h.generic(c, "TiltDown", func(ctx context.Context, cam *visca.Camera) error {
		return cam.TiltDown(ctx, 0x0e)
	})
}

func (h *Handlers) PanLeft(c echo.Context) error {
	return h.generic(c, "PanLeft", func(ctx context.Context, cam *visca.Camera) error {
		return cam.PanLeft(ctx, 0x0b)
	})
}

func (h *Handlers) PanRight(c echo.Context) error {
	return h.generic(c, "PanRight", func(ctx context.Context, cam *visca.Camera) error {
		return cam.PanRight(ctx, 0x0b)
	})
}

func (h *Handlers) PanTiltStop(c echo.Context) error {
	return h.generic(c, "PanTiltStop", func(ctx context.Context, cam *visca.Camera) error {
		return cam.PanTiltStop(ctx)
	})
}
