package entities

type LineWebhook struct {
	Destination string `json:"destination"`
	Events      []struct {
		Type    string `json:"type"`
		Message struct {
			Type       string `json:"type"`
			ID         string `json:"id"`
			QuoteToken string `json:"quoteToken"`
			Text       string `json:"text"`
		} `json:"message"`
		WebhookEventID  string `json:"webhookEventId"`
		DeliveryContext struct {
			IsRedelivery bool `json:"isRedelivery"`
		} `json:"deliveryContext"`
		Timestamp int64 `json:"timestamp"`
		Source    struct {
			Type    string `json:"type"`
			GroupID string `json:"groupId"`
			UserID  string `json:"userId"`
		} `json:"source"`
		ReplyToken string `json:"replyToken"`
		Mode       string `json:"mode"`
	} `json:"events"`
}
