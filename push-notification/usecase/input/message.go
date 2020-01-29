package input

import (
	"github.com/abyssparanoia/rapid-go/push-notification/domain/model"
)

// MessageRequest ... push notification message
type MessageRequest struct {
	Title   string            `json:"title"`
	Body    string            `json:"body"`
	Data    map[string]string `json:"data"`
	IOS     *MessageIOS       `json:"ios"`
	Android *MessageAndroid   `json:"android"`
	Web     *MessageWeb       `json:"web"`
}

// MessageIOS ... push notification message only ios
type MessageIOS struct {
	Badge int    `json:"badge,omitempty"`
	Sound string `json:"sound,omitempty"`
}

// OutputModel ... output model
func (e *MessageIOS) OutputModel() *model.MessageIOS {
	if e == nil {
		return nil
	}
	return model.NewMessageIOS(e.Badge, e.Sound)
}

// MessageAndroid ... push notification message only android
type MessageAndroid struct {
	ClickAction string `json:"click_action,omitempty"`
	Sound       string `json:"sound,omitempty"`
	Tag         string `json:"badge,omitempty"`
}

// OutputModel ... output model
func (e *MessageAndroid) OutputModel() *model.MessageAndroid {
	if e == nil {
		return nil
	}
	return model.NewMessageAndroid(e.ClickAction, e.Sound, e.Tag)
}

// MessageWeb ... push notification message only web
type MessageWeb struct {
	Icon string `json:"icon,omitempty"`
}

// OutputModel ... output model
func (e *MessageWeb) OutputModel() *model.MessageWeb {
	if e == nil {
		return nil
	}
	return model.NewMessageWeb(e.Icon)
}

// MessageSendToUser ...
type MessageSendToUser struct {
	AppID   string
	UserID  string
	Message *model.Message
}

// NewMessageSendToUser ... new message send to user input
func NewMessageSendToUser(appID, userID string, messsage *MessageRequest) *MessageSendToUser {
	return &MessageSendToUser{
		AppID:  appID,
		UserID: userID,
		Message: &model.Message{
			Title:   messsage.Title,
			Body:    messsage.Body,
			Data:    messsage.Data,
			IOS:     messsage.IOS.OutputModel(),
			Android: messsage.Android.OutputModel(),
			Web:     messsage.Web.OutputModel(),
		},
	}

}
