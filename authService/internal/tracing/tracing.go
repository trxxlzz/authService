package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"log"
)

func Init(service string) {
	// Настроим Jaeger
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const", // фиксированный тип, всегда отправлять данные
			Param: 1,       // 100% сэмплирование
		},
	}

	_, err := cfg.InitGlobalTracer(service)
	if err != nil {
		log.Fatal("Failed to init tracing", zap.Error(err))
	}
}
