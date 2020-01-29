package model

// Message ... push notification message
type Message struct {
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

// NewMessageIOS ... new message IOS
func NewMessageIOS(badge int, sound string) *MessageIOS {
	return &MessageIOS{
		Badge: badge,
		Sound: sound,
	}
}

// MessageAndroid ... push notification message only android
type MessageAndroid struct {
	ClickAction string `json:"click_action,omitempty"`
	Sound       string `json:"sound,omitempty"`
	Tag         string `json:"badge,omitempty"`
}

// NewMessageAndroid ... new message android
func NewMessageAndroid(clickAction, sound, tag string) *MessageAndroid {
	return &MessageAndroid{
		ClickAction: clickAction,
		Sound:       sound,
		Tag:         tag,
	}
}

// MessageWeb ... push notification message only web
type MessageWeb struct {
	Icon string `json:"icon,omitempty"`
}

// NewMessageWeb ... new message web
func NewMessageWeb(icon string) *MessageWeb {
	return &MessageWeb{
		Icon: icon,
	}
}
