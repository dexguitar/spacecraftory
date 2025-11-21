package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	"github.com/dexguitar/spacecraftory/notification/internal/client/http"
	tgClient "github.com/dexguitar/spacecraftory/notification/internal/client/http/telegram"
	"github.com/dexguitar/spacecraftory/notification/internal/config"
	kafkaConverter "github.com/dexguitar/spacecraftory/notification/internal/converter/kafka"
	decoder "github.com/dexguitar/spacecraftory/notification/internal/converter/kafka/decoder"
	"github.com/dexguitar/spacecraftory/notification/internal/service"
	"github.com/dexguitar/spacecraftory/notification/internal/service/consumer/order_assembled_consumer"
	"github.com/dexguitar/spacecraftory/notification/internal/service/consumer/order_paid_consumer"
	tgService "github.com/dexguitar/spacecraftory/notification/internal/service/telegram"
	"github.com/dexguitar/spacecraftory/platform/pkg/closer"
	wrappedKafka "github.com/dexguitar/spacecraftory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/dexguitar/spacecraftory/platform/pkg/kafka/consumer"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
	kafkaMiddleware "github.com/dexguitar/spacecraftory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	orderPaidConsumer      wrappedKafka.Consumer
	orderPaidDecoder       kafkaConverter.OrderPaidDecoder
	orderPaidConsumerGroup sarama.ConsumerGroup

	orderAssembledConsumer      wrappedKafka.Consumer
	orderAssembledDecoder       kafkaConverter.OrderAssembledDecoder
	orderAssembledConsumerGroup sarama.ConsumerGroup

	telegramClient  http.TelegramClient
	telegramBot     *bot.Bot
	telegramService service.TelegramService

	orderPaidConsumerService      service.ConsumerService
	orderAssembledConsumerService service.ConsumerService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = order_paid_consumer.NewService(d.OrderPaidConsumer(), d.OrderPaidDecoder(), d.TelegramService(ctx))
	}
	return d.orderPaidConsumerService
}

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = order_assembled_consumer.NewService(d.OrderAssembledConsumer(), d.OrderAssembledDecoder(), d.TelegramService(ctx))
	}

	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		orderPaidConsumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.orderPaidConsumerGroup.Close()
		})

		d.orderPaidConsumerGroup = orderPaidConsumerGroup
	}

	return d.orderPaidConsumerGroup
}

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		orderAssembledConsumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.orderAssembledConsumerGroup.Close()
		})

		d.orderAssembledConsumerGroup = orderAssembledConsumerGroup
	}

	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
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

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder
}

func (d *diContainer) TelegramClient() http.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = tgClient.NewClient(d.telegramBot)
	}

	return d.telegramClient
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = tgService.NewService(d.TelegramClient())
	}

	return d.telegramService
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().TelegramBot.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}
