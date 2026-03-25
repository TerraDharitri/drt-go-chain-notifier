package common

import (
	"time"

	"github.com/TerraDharitri/drt-go-chain-notifier/data"
)

// StatusMetricsHandler defines the behavior of a component that handles status metrics
type StatusMetricsHandler interface {
	AddRequest(path string, duration time.Duration)
	GetAll() map[string]*data.EndpointMetricsResponse
	GetMetricsForPrometheus() string
	IsInterfaceNil() bool
}
