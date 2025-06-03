package borderless

import (
	"backend/models"
	"fmt"
)

func (hc Borderless) GetCustomerIdentity(email string, lastname string) (map[string]interface{}, error) {
	response, err := hc.MakeRequest("GET",
		fmt.Sprintf("%s/identities?emailPrefix=%s&namePrefix=%s&type=%s", hc.BaseUrl, email, lastname, "Personal"),
		nil,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (hc Borderless) CreateCustomerIdentity(identity models.BorderlessIdentity) (map[string]interface{}, error) {
	requestData := map[string]interface{}{
		"firstName":   identity.FirstName,
		"lastName":    identity.LastName,
		"taxId":       identity.TaxId,
		"dateOfBirth": identity.DateOfBirth,
		"email":       identity.Email,
		"phone":       identity.Phone,
		"address":     identity.Address,
	}
	response, err := hc.MakeRequest(
		"POST",
		fmt.Sprintf("%s/identities/personal", hc.BaseUrl),
		requestData,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (hc Borderless) UploadCustomerIdentityDocument(identityId string, kyc models.KYC, kycData models.KYCData) (map[string]interface{}, error) {
	requestData := map[string]interface{}{
		"issuingCountry": "US",
		"type":           kyc.IDType,
		"issuedDate":     kyc.IssueDate,
		"expiryDate":     kyc.ExpiryDate,
		"imageFront":     kycData.FrontPhoto,
		"imageBack":      kycData.BackPhoto,
	}

	response, err := hc.MakeRequest(
		"PUT",
		fmt.Sprintf("%s/identities/%s/documents", hc.BaseUrl, identityId),
		requestData,
	)

	if err != nil {
		return nil, err
	}

	return response, nil

}
