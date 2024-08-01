package serializers

type Hmac struct {
	HmacSecret string `json:"hmacSecret"`
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

type Webhook struct {
	Address          string  `json:"address"`
	Amount           string  `json:"amount"`
	CounterAddress   string  `json:"counterAddress"`
	Asset            string  `json:"asset"`
	BlockNumber      int     `json:"blockNumber"`
	TxID             string  `json:"txId"`
	Type             string  `json:"type"`
	TokenID          *string `json:"tokenId"`
	Chain            string  `json:"chain"`
	SubscriptionType string  `json:"subscriptionType"`
}
