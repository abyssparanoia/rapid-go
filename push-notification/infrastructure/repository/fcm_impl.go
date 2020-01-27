package repository

import (
	"context"

	"firebase.google.com/go/messaging"
	"github.com/abyssparanoia/rapid-go/push-notification/domain/repository"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
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
		log.Errorm(ctx, "r.fCli.SubscribeToTopic", err)
		return err
	}
	if res.FailureCount > 0 {
		for _, rerr := range res.Errors {
			err = log.Warninge(ctx, "SubscribeToTopic index: %d, reason: %s", rerr.Index, rerr.Reason)
			return err
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
		log.Errorm(ctx, "r.fCli.UnsubscribeFromTopic", err)
		return err
	}
	if res.FailureCount > 0 {
		for _, rerr := range res.Errors {
			err = log.Warninge(ctx, "UnsubscribeFromTopic index: %d, reason: %s", rerr.Index, rerr.Reason)
			return err
		}
	}
	return nil
}

// NewFcm ... new fcm repository
func NewFcm(messagingClient *messaging.Client,
	serverKey string) repository.Fcm {
	return &fcm{messagingClient, serverKey}
}
