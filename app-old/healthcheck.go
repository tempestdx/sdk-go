package app

import "context"

//go:generate go run golang.org/x/tools/cmd/stringer -type=HealthCheckStatus -linecomment

type HealthCheckResponse struct {
	Status  HealthCheckStatus
	Message string
}

type HealthCheckStatus int

const (
	HealthCheckStatusUnknown   HealthCheckStatus = iota // unknown
	HealthCheckStatusHealthy                            // healthy
	HealthCheckStatusDegraded                           // degraded
	HealthCheckStatusDisrupted                          // disrupted
)

type HealthCheckFunc func(context.Context) (*HealthCheckResponse, error)
