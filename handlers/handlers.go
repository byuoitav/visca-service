package handlers

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/byuoitav/visca"
	viscaservice "github.com/byuoitav/visca-service"
	"github.com/labstack/echo"
)

type CreateCameraFunc func(ctx context.Context, addr string) (*visca.Camera, error)

type Handlers struct {
	CreateCamera   CreateCameraFunc
	EventPublisher viscaservice.EventPublisher
	Resolver       net.Resolver
}

func (h *Handlers) getSourceIP(c echo.Context) (net.IP, error) {
	host := c.RealIP()
	var err error

	if strings.Contains(host, ":") {
		host, _, err = net.SplitHostPort(host)
		if err != nil {
			return nil, fmt.Errorf("unable to split host/port: %w", err)
		}
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return nil, fmt.Errorf("unable to parse remote ip %q", ip)
	}

	return ip, nil
}

func (h *Handlers) getCameraIP(ctx context.Context, addr string) (net.IP, error) {
	host := addr
	var err error

	if strings.Contains(host, ":") {
		host, _, err = net.SplitHostPort(host)
		if err != nil {
			return nil, fmt.Errorf("unable to split host/port: %w", err)
		}
	}

	// figure out if it's an ip or not
	ip := net.ParseIP(host)
	if ip == nil {
		addrs, err := h.Resolver.LookupHost(ctx, host)
		if err != nil {
			return nil, fmt.Errorf("unable to reverse lookup ip: %w", err)
		}

		if len(addrs) == 0 {
			return nil, errors.New("no camera IP addresses found")
		}

		for _, addr := range addrs {
			ip = net.ParseIP(addr)
			if ip != nil {
				break
			}
		}
	}

	return ip, nil
}

type genericAction func(context.Context, *visca.Camera) error

func (h *Handlers) generic(c echo.Context, actionName string, action genericAction) error {
	var err error
	info := viscaservice.RequestInfo{
		Action: actionName,
	}

	info.SourceIP, err = h.getSourceIP(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to get source ip: %s", err))
	}

	addr := c.Param("address")
	if len(addr) == 0 {
		err = fmt.Errorf("must include the address of the camera")
		go h.EventPublisher.Error(context.Background(), viscaservice.RequestError{
			RequestInfo: info,
			Error:       err,
			StatusCode:  http.StatusBadRequest,
		})

		return c.String(http.StatusBadRequest, err.Error())
	}

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

	if err := action(ctx, cam); err != nil {
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
