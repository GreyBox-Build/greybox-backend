package serializers

import (
	"github.com/google/uuid"
)

type Hmac struct {
	HmacSecret uuid.UUID `json:"hmacSecret"`
}
