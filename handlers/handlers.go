package handlers

import (
	"context"

	"github.com/byuoitav/visca"
)

type CreateCameraFunc func(ctx context.Context, addr string) (*visca.Camera, error)

type Handlers struct {
	CreateCamera CreateCameraFunc
}
