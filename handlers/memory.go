package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	viscaservice "github.com/byuoitav/visca-service"
	"github.com/labstack/echo"
)

func (h *Handlers) MemoryRecall(c echo.Context) error {
	var err error
	info := viscaservice.RequestInfo{
		Action: "MemoryRecall",
	}

	info.SourceIP, err = h.getSourceIP(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to get source ip: %s", err))
	}

	addr := c.Param("address")
	channel, err := strconv.Atoi(c.Param("channel"))
	switch {
	case len(addr) == 0:
		err = fmt.Errorf("must include the address of the camera")
		go h.EventPublisher.Error(context.Background(), viscaservice.RequestError{
			RequestInfo: info,
			Error:       err,
			StatusCode:  http.StatusBadRequest,
		})

		return c.String(http.StatusBadRequest, err.Error())
	case err != nil:
		go h.EventPublisher.Error(context.Background(), viscaservice.RequestError{
			RequestInfo: info,
			Error:       err,
			StatusCode:  http.StatusBadRequest,
		})

		return c.String(http.StatusBadRequest, err.Error())
	}

	info.Action += strconv.Itoa(channel)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	info.CameraIP, err = h.getCameraIP(ctx, addr)
	if err != nil {
		err = fmt.Errorf("unable to get camera ip: %w", err)
		go h.EventPublisher.Error(context.Background(), viscaservice.RequestError{
			RequestInfo: info,
			Error:       err,
			StatusCode:  http.StatusInternalServerError,
		})

		return c.String(http.StatusInternalServerError, err.Error())
	}

	cam, err := h.CreateCamera(ctx, addr)
	if err != nil {
		err = fmt.Errorf("unable to create camera: %w", err)
		go h.EventPublisher.Error(context.Background(), viscaservice.RequestError{
			RequestInfo: info,
			Error:       err,
			StatusCode:  http.StatusInternalServerError,
		})

		return c.String(http.StatusInternalServerError, err.Error())
	}

	if err := cam.MemoryRecall(ctx, byte(channel)); err != nil {
		go h.EventPublisher.Error(context.Background(), viscaservice.RequestError{
			RequestInfo: info,
			Error:       err,
			StatusCode:  http.StatusInternalServerError,
		})

		return c.String(http.StatusInternalServerError, err.Error())
	}

	go h.EventPublisher.Publish(context.Background(), info)
	return c.NoContent(http.StatusOK)
}
