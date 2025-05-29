package serializers

import "encoding/json"

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

type EventObject struct {
	Type               string      `json:"type"`
	ID                 string      `json:"id"`
	Partner            string      `json:"partner"`
	CustomerName       string      `json:"customer_name"`
	CollectionCurrency string      `json:"collection_currency"`
	CollectionRail     string      `json:"collection_rail"`
	CollectionAmount   json.Number `json:"collection_amount"`
	BlockchainNetwork  string      `json:"blockchain_network"`
	BlockchainToken    string      `json:"blockchain_token"`
	BlockchainProof    string      `json:"blockchain_proof"`
	TokenAmount        json.Number `json:"token_amount"`
	Description        string      `json:"description"`
}

type Event struct {
	APIVersion     string      `json:"api_version"`
	EventID        string      `json:"event_id"`
	EventCategory  string      `json:"event_category"`
	EventType      string      `json:"event_type"`
	EventObject    EventObject `json:"event_object"`
	EventCreatedAt string      `json:"event_created_at"`
}
