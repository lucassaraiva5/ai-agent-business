package presenter

import (
	"ai-agent-business/internal/app/domain/health"
	"ai-agent-business/internal/app/transport/outbound"
)

func HealthResponse(read *health.Health, write *health.Health, redis *health.Health) *outbound.HealthResponse {
	response := &outbound.HealthResponse{}

	if read.Up {
		response.ReadDB = "UP"
	} else {
		response.ReadDB = "DOWN"
	}

	if write.Up {
		response.WriteDB = "UP"
	} else {
		response.WriteDB = "DOWN"
	}

	if redis.Up {
		response.RedisDB = "UP"
	} else {
		response.RedisDB = "DOWN"
	}

	response.Healthy = read.Up && write.Up && redis.Up

	return response
}
