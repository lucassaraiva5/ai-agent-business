package handler

import (
	"net/http"
	"testing"

	"ai-agent-business/test"
)

func TestShouldHealth(t *testing.T) {
	defer test.StopApplication()
	test.InitLogger()
	test.ComposeUp(t)
	test.StartApplication()

	_ = test.Request().
		Get("/health").
		Expect(t).
		Status(http.StatusOK).
		Type("json").
		JSON(map[string]interface{}{
			"healthy":  true,
			"read_db":  "UP",
			"write_db": "UP",
			"redis_db": "UP",
		}).Done()
}
