package entity

import (
	"fmt"

	"firebase.google.com/go/messaging"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/model"
)

// NewMessageFromModel ... new message from model
func NewMessageFromModel(m *model.Message, serverKey string) *messaging.Message {

	e := &messaging.Message{}

	if m.IOS == nil {
		m.IOS = &model.MessageIOS{
			Badge: 1,
		}
	}

	if m.Android == nil {
		m.Android = &model.MessageAndroid{}
	}

	if m.Web == nil {
		m.Web = &model.MessageWeb{}
	}

	e.Notification = &messaging.Notification{
		Title: m.Title,
		Body:  m.Body,
	}

	e.Data = m.Data

	e.APNS = &messaging.APNSConfig{
		Headers: map[string]string{
			"apns-priority": "10",
		},
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				Badge: &m.IOS.Badge,
				Sound: m.IOS.Sound,
			},
			CustomData: map[string]interface{}{
				"notification_foreground": true,
			},
		}}

	e.Android = &messaging.AndroidConfig{
		Notification: &messaging.AndroidNotification{
			ClickAction: m.Android.ClickAction,
			Sound:       m.Android.Sound,
			Tag:         m.Android.Tag,
		}}

	e.Webpush = &messaging.WebpushConfig{
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", serverKey),
		},
		Notification: &messaging.WebpushNotification{
			Icon: m.Web.Icon,
		},
	}

	return e
}
