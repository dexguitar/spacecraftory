package assembly

import (
	srv "github.com/dexguitar/spacecraftory/assembly/internal/service"
)

type service struct {
	assemblyProducerService srv.ProducerService
	assemblyConsumerService srv.ConsumerService
}

func NewService(assemblyProducerService srv.ProducerService, assemblyConsumerService srv.ConsumerService) *service {
	return &service{
		assemblyProducerService: assemblyProducerService,
		assemblyConsumerService: assemblyConsumerService,
	}
}
