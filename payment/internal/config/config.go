package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/dexguitar/spacecraftory/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	Tracing     TracingConfig
	PaymentGRPC PaymentGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	tracingCfg, err := env.NewPaymentTracingConfig()
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:      loggerCfg,
		Tracing:     tracingCfg,
		PaymentGRPC: paymentGRPCCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
