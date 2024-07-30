package serializers

import (
	"github.com/google/uuid"
)

type Hmac struct {
	HmacSecret uuid.UUID `json:"hmacSecret"`
}

type Subscription struct {
	Type string           `json:"type"`
	Attr SubscriptionAttr `json:"attr"`
}

type SubscriptionAttr struct {
	Chain   string `json:"chain"`
	Address string `json:"address"`
	Url     string `json:"url"`
}
