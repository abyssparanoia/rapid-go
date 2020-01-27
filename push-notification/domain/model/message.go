package model

// Message ... push notification message
type Message struct {
	Title   string            `json:"title"   firestore:"title"`
	Body    string            `json:"body"    firestore:"body"`
	Data    map[string]string `json:"data"    firestore:"data"`
	IOS     *MessageIOS       `json:"ios"     firestore:"ios"`
	Android *MessageAndroid   `json:"android" firestore:"android"`
	Web     *MessageWeb       `json:"web"     firestore:"web"`
}

// MessageIOS ... push notification message only ios
type MessageIOS struct {
	Badge int    `json:"badge,omitempty" firestore:"badge"`
	Sound string `json:"sound,omitempty" firestore:"sound"`
}

// MessageAndroid ... push notification message only android
type MessageAndroid struct {
	ClickAction string `json:"click_action,omitempty" firestore:"click_action"`
	Sound       string `json:"sound,omitempty"        firestore:"sound"`
	Tag         string `json:"badge,omitempty"        firestore:"badge"`
}

// MessageWeb ... push notification message only web
type MessageWeb struct {
	Icon string `json:"icon,omitempty" firestore:"icon"`
}
