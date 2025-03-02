package health

import (
	"os"

	"github.com/gin-gonic/gin"
)

type HealthProbe struct {
	Hostname string `json:"hostname"`

	Status string `json:"status"`
}

func newHealthProbe() *HealthProbe {
	hostname, _ := os.Hostname()

	return &HealthProbe{
		Status:   "ok",
		Hostname: hostname,
	}
}

func NewHealthProbeWithCustomGin(engine *gin.Engine) *HealthProbe {
	healthProbe := newHealthProbe()

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, healthProbe)
	})

	return healthProbe
}
