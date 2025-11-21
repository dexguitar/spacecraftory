package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"
	"time"

	"go.uber.org/zap"

	"github.com/dexguitar/spacecraftory/notification/internal/client/http"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
)

const chatID = 2407852

//go:embed templates/order_assembled_notification.tmpl templates/order_paid_notification.tmpl
var templateFS embed.FS

type orderPaidTemplateData struct {
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
	RegisteredAt    time.Time
}

type orderAssembledTemplateData struct {
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
	RegisteredAt time.Time
}

var (
	orderPaidTemplate      = template.Must(template.ParseFS(templateFS, "templates/order_paid_notification.tmpl"))
	orderAssembledTemplate = template.Must(template.ParseFS(templateFS, "templates/order_assembled_notification.tmpl"))
)

type service struct {
	telegramClient http.TelegramClient
}

// NewService создает новый Telegram сервис
func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

// SendUFONotification отправляет уведомление о новом наблюдении UFO
func (s *service) SendOrderPaidNotification(ctx context.Context, orderUUID, userUUID, paymentMethod, transactionUUID string) error {
	message, err := s.buildOrderPaidMessage(orderUUID, userUUID, paymentMethod, transactionUUID)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) SendOrderAssembledNotification(ctx context.Context, orderUUID, userUUID string, buildTimeSec int64) error {
	message, err := s.buildOrderAssembledMessage(orderUUID, userUUID, buildTimeSec)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

// buildOrderPaidMessage создает сообщение о платеже заказа из шаблона
func (s *service) buildOrderPaidMessage(orderUUID, userUUID, paymentMethod, transactionUUID string) (string, error) {
	data := orderPaidTemplateData{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PaymentMethod:   paymentMethod,
		TransactionUUID: transactionUUID,
		RegisteredAt:    time.Now(),
	}

	var buf bytes.Buffer
	err := orderPaidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// buildOrderAssembledMessage создает сообщение о сборке заказа из шаблона
func (s *service) buildOrderAssembledMessage(orderUUID, userUUID string, buildTimeSec int64) (string, error) {
	data := orderAssembledTemplateData{
		OrderUUID:    orderUUID,
		UserUUID:     userUUID,
		BuildTimeSec: buildTimeSec,
		RegisteredAt: time.Now(),
	}

	var buf bytes.Buffer
	err := orderAssembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
