package borderless

import (
	"fmt"

	"github.com/google/uuid"
)

func (hc Borderless) CreateBorderlessVirtualAccount(
	accountId string,
	fiat string,
	asset string,
	countryCode string,
	identityId string) (map[string]interface{}, error) {
	requestData := map[string]interface{}{
		"fiat":                   fiat,
		"country":                countryCode,
		"asset":                  asset,
		"counterPartyIdentityId": identityId,
	}

	idempotencyKey := uuid.New()
	hc.Headers["idempotency-key"] = idempotencyKey.String()
	response, err := hc.MakeRequest(
		"POST",
		fmt.Sprintf("%s/accounts/%s/virtual-accounts", hc.BaseUrl, accountId),
		requestData,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}
