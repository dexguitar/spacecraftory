package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/dexguitar/spacecraftory/assembly/internal/config"
	kafkaConverter "github.com/dexguitar/spacecraftory/assembly/internal/converter/kafka"
	decoder "github.com/dexguitar/spacecraftory/assembly/internal/converter/kafka/decoder"
	"github.com/dexguitar/spacecraftory/assembly/internal/service"
	assemblyConsumer "github.com/dexguitar/spacecraftory/assembly/internal/service/consumer/assembly_consumer"
	assemblyProducer "github.com/dexguitar/spacecraftory/assembly/internal/service/producer/assembly_producer"
	"github.com/dexguitar/spacecraftory/platform/pkg/closer"
	wrappedKafka "github.com/dexguitar/spacecraftory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/dexguitar/spacecraftory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/dexguitar/spacecraftory/platform/pkg/kafka/producer"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
	kafkaMiddleware "github.com/dexguitar/spacecraftory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	assemblyProducerService service.ProducerService
	assemblyConsumerService service.ConsumerService
	consumerGroup           sarama.ConsumerGroup
	orderPaidConsumer       wrappedKafka.Consumer
	orderPaidDecoder        kafkaConverter.OrderPaidDecoder
	syncProducer            sarama.SyncProducer
	shipAssembledProducer   wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AssemblyProducerService() service.ProducerService {
	if d.assemblyProducerService == nil {
		d.assemblyProducerService = assemblyProducer.NewService(d.ShipAssembledProducer())
	}

	return d.assemblyProducerService
}

func (d *diContainer) AssemblyConsumerService() service.ConsumerService {
	if d.assemblyConsumerService == nil {
		d.assemblyConsumerService = assemblyConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.AssemblyProducerService(),
		)
	}

	return d.assemblyConsumerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) ShipAssembledProducer() wrappedKafka.Producer {
	if d.shipAssembledProducer == nil {
		d.shipAssembledProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.shipAssembledProducer
}
