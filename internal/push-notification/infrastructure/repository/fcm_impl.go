package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"firebase.google.com/go/messaging"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/repository"

	"github.com/abyssparanoia/rapid-go/internal/push-notification/infrastructure/internal/entity"
)

type fcm struct {
	messagingClient *messaging.Client
	serverKey       string
}

func (r *fcm) SubscribeTopic(
	ctx context.Context,
	topic string,
	tokens []string) error {

	res, err := r.messagingClient.SubscribeToTopic(ctx, tokens, topic)
	if err != nil {
		return err
	}
	if res.FailureCount > 0 {
		for _, rerr := range res.Errors {
			return errors.New(fmt.Sprintf("SubscribeToTopic index: %d, reason: %s", rerr.Index, rerr.Reason))
		}
	}
	return nil
}

func (r *fcm) Unsubscribe(
	ctx context.Context,
	topic string,
	tokens []string) error {

	res, err := r.messagingClient.UnsubscribeFromTopic(ctx, tokens, topic)
	if err != nil {
		return err
	}
	if res.FailureCount > 0 {
		for _, rerr := range res.Errors {
			return errors.New(fmt.Sprintf("UnsubscribeFromTopic index: %d, reason: %s", rerr.Index, rerr.Reason))
		}
	}
	return nil
}

func (r *fcm) SendMessageByTokens(
	ctx context.Context,
	appID string,
	tokens []string,
	message *model.Message) error {

	messageEntity := entity.NewMessageFromModel(message, r.serverKey)

	multiMessage := &messaging.MulticastMessage{
		Tokens:       tokens,
		Notification: messageEntity.Notification,
		Data:         messageEntity.Data,
		APNS:         messageEntity.APNS,
		Android:      messageEntity.Android,
		Webpush:      messageEntity.Webpush,
	}

	_, err := r.messagingClient.SendMulticast(ctx, multiMessage)
	if err != nil {
		return err
	}
	return nil

}

func (r *fcm) SendMessageByTopic(
	ctx context.Context,
	appID string,
	topic string,
	message *model.Message) error {

	messageEntity := entity.NewMessageFromModel(message, r.serverKey)
	messageEntity.Topic = topic
	_, err := r.messagingClient.Send(ctx, messageEntity)
	if err != nil {
		return err
	}
	return nil
}

// NewFcm ... new fcm repository
func NewFcm(messagingClient *messaging.Client,
	serverKey string) repository.Fcm {
	return &fcm{messagingClient, serverKey}
}
