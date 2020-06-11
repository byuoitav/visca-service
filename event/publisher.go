package event

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	viscaservice "github.com/byuoitav/visca-service"
)

type Publisher struct {
	GeneratingSystem string

	Resolver net.Resolver
	URL      string
}

type event struct {
	GeneratingSystem string      `json:"generating-system"`
	Timestamp        time.Time   `json:"timestamp"`
	Tags             []string    `json:"event-tags"`
	TargetDevice     deviceInfo  `json:"target-device"`
	AffectedRoom     roomInfo    `json:"affected-room"`
	Key              string      `json:"key"`
	Value            string      `json:"value"`
	User             string      `json:"user"`
	Data             interface{} `json:"data,omitempty"`
}

type roomInfo struct {
	BuildingID string `json:"buildingID,omitempty"`
	RoomID     string `json:"roomID,omitempty"`
}

type deviceInfo struct {
	roomInfo
	DeviceID string `json:"deviceID,omitempty"`
}

func (p *Publisher) Publish(ctx context.Context, info viscaservice.RequestInfo) error {
	event := event{
		GeneratingSystem: p.GeneratingSystem,
		Timestamp:        time.Now(),
		User:             info.SourceIP.String(),
		Key:              info.Action,
		Value:            info.CameraIP.String(),
		Tags: []string{
			"cameraControl",
			"viscaService",
		},
	}

	event = p.handleIPs(ctx, info, event)

	return p.publish(ctx, event)
}

func (p *Publisher) Error(ctx context.Context, err viscaservice.RequestError) error {
	event := event{
		GeneratingSystem: p.GeneratingSystem,
		Timestamp:        time.Now(),
		User:             err.SourceIP.String(),
		Key:              err.Action,
		Value:            err.CameraIP.String(),
		Tags: []string{
			"cameraControl",
			"viscaService",
			"error",
		},
		Data: struct {
			StatusCode int    `json:"statusCode"`
			Error      string `json:"error"`
		}{
			StatusCode: err.StatusCode,
			Error:      err.Error.Error(),
		},
	}

	event = p.handleIPs(ctx, err.RequestInfo, event)
	return p.publish(ctx, event)
}

func (p *Publisher) handleIPs(ctx context.Context, info viscaservice.RequestInfo, event event) event {
	// lookup hostname for source
	if info.SourceIP != nil {
		sources, _ := p.Resolver.LookupAddr(ctx, info.SourceIP.String())
		for _, source := range sources {
			trimmed := strings.TrimSuffix(source, ".byu.edu.")
			split := strings.SplitN(trimmed, "-", 3)
			if len(split) == 3 {
				event.User = trimmed
			}
		}
	}

	if info.CameraIP != nil {
		cameras, _ := p.Resolver.LookupAddr(ctx, info.CameraIP.String())
		for _, camera := range cameras {
			trimmed := strings.TrimSuffix(camera, ".byu.edu.")
			split := strings.SplitN(trimmed, "-", 3)
			if len(split) == 3 {
				event.TargetDevice.BuildingID = split[0]
				event.TargetDevice.RoomID = event.TargetDevice.BuildingID + "-" + split[1]
				event.TargetDevice.DeviceID = event.TargetDevice.RoomID + "-" + split[2]

				event.AffectedRoom = event.TargetDevice.roomInfo
			}
		}
	}

	return event
}

func (p *Publisher) publish(ctx context.Context, event event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("unable to marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("unable to build request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 2 {
		return fmt.Errorf("got a %v response from event url", resp.StatusCode)
	}

	return nil
}
