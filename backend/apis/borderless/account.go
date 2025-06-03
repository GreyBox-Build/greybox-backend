package borderless

import "fmt"

func (hc Borderless) CreateBorderlessAccount(name string, identityId string) (map[string]interface{}, error) {
	requestData := map[string]interface{}{
		"name":       name,
		"identityId": identityId,
	}

	response, err := hc.MakeRequest(
		"POST",
		fmt.Sprintf("%s/accounts", hc.BaseUrl),
		requestData,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}
