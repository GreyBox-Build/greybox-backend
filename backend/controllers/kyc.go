package controllers

import (
	"backend/apis/borderless"
	"backend/models"
	"backend/serializers"
	"backend/utils"
	"backend/utils/tokens"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUserKYC(c *gin.Context) {
	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	includePhotos := c.GetHeader("include_photos") == "true"

	if includePhotos {
		// Struct returned includes photo data
		kycWithPhotos, err := models.GetKycByUserIDWithPhotos(user.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		kycResponse := serializers.KYC{
			ID:              kycWithPhotos.Kyc.ID,
			IDType:          kycWithPhotos.Kyc.IDType,
			IssueDate:       kycWithPhotos.Kyc.IssueDate,
			ExpiryDate:      kycWithPhotos.Kyc.ExpiryDate,
			TaxId:           kycWithPhotos.Kyc.TaxId,
			IdNumber:        kycWithPhotos.Kyc.IdNumber,
			DateOfBirth:     kycWithPhotos.Kyc.DateOfBirth,
			StreetAddress:   kycWithPhotos.Kyc.StreetAddress,
			City:            kycWithPhotos.Kyc.City,
			State:           kycWithPhotos.Kyc.State,
			PostalCode:      kycWithPhotos.Kyc.PostalCode,
			Country:         kycWithPhotos.Kyc.Country,
			Phone:           kycWithPhotos.Kyc.Phone,
			Status:          serializers.KYCStatus(kycWithPhotos.Kyc.Status),
			RejectionReason: kycWithPhotos.Kyc.RejectionReason,
			CreatedAt:       kycWithPhotos.Kyc.CreatedAt,
			UpdatedAt:       kycWithPhotos.Kyc.UpdatedAt,
			RejectedAt:      kycWithPhotos.Kyc.RejectedAt,
			ApprovedAt:      kycWithPhotos.Kyc.ApprovedAt,
			FrontPhoto:      kycWithPhotos.FrontPhoto,
			BackPhoto:       kycWithPhotos.BackPhoto,
		}

		c.JSON(http.StatusOK, gin.H{"kyc": kycResponse})
		return
	}

	// Fallback: basic KYC struct
	kyc, err := models.GetKYCByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	kycResponse := serializers.KYC{
		ID:              kyc.ID,
		IDType:          kyc.IDType,
		IssueDate:       kyc.IssueDate,
		ExpiryDate:      kyc.ExpiryDate,
		TaxId:           kyc.TaxId,
		IdNumber:        kyc.IdNumber,
		DateOfBirth:     kyc.DateOfBirth,
		StreetAddress:   kyc.StreetAddress,
		City:            kyc.City,
		State:           kyc.State,
		PostalCode:      kyc.PostalCode,
		Country:         kyc.Country,
		Phone:           kyc.Phone,
		Status:          serializers.KYCStatus(kyc.Status),
		RejectionReason: kyc.RejectionReason,
		CreatedAt:       kyc.CreatedAt,
		UpdatedAt:       kyc.UpdatedAt,
		RejectedAt:      kyc.RejectedAt,
		ApprovedAt:      kyc.ApprovedAt,
	}

	c.JSON(http.StatusOK, gin.H{"kyc": kycResponse})
}

func GetKYCS(c *gin.Context) {
	var request serializers.KYCFilterRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// if email is provided then find user Id by the email
	if request.Email != nil {
		// find the user by their email
		user, exists := models.FindUserByEmail(*request.Email)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		// update request with user ID
		request.UserID = &user.ID
	}

	m_kycs, err := models.FilterKYC(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// convert kyc models to serializers
	var kycs []serializers.KYC
	for _, kyc := range m_kycs {
		kycs = append(kycs, serializers.KYC{
			ID:              kyc.ID,
			IDType:          kyc.IDType,
			IssueDate:       kyc.IssueDate,
			ExpiryDate:      kyc.ExpiryDate,
			TaxId:           kyc.TaxId,
			IdNumber:        kyc.IdNumber,
			DateOfBirth:     kyc.DateOfBirth,
			StreetAddress:   kyc.StreetAddress,
			City:            kyc.City,
			State:           kyc.State,
			PostalCode:      kyc.PostalCode,
			Country:         kyc.Country,
			Phone:           kyc.Phone,
			Status:          serializers.KYCStatus(kyc.Status),
			RejectionReason: kyc.RejectionReason,
			CreatedAt:       kyc.CreatedAt,
			UpdatedAt:       kyc.UpdatedAt,
			RejectedAt:      kyc.RejectedAt,
			ApprovedAt:      kyc.ApprovedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"kycs": kycs,
	})
}

func CreateKYC(c *gin.Context) {
	// Max file size: 5MB
	const maxSize = 5 << 20

	// Accepted MIME types
	validMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	// Bind form data
	var request struct {
		Email         string `form:"email" binding:"required,email"`
		IDType        string `form:"id_type" binding:"required"`
		IssueDate     string `form:"issue_date" binding:"required"`
		ExpiryDate    string `form:"expiry_date" binding:"required"`
		TaxId         string `form:"tax_id" binding:"required"`
		IdNumber      string `form:"id_number" binding:"required"`
		DateOfBirth   string `form:"date_of_birth" binding:"required"`
		Phone         string `form:"phone" binding:"required"`
		StreetAddress string `form:"street_address" binding:"required"`
		City          string `form:"city" binding:"required"`
		State         string `form:"state" binding:"required"`
		PostalCode    string `form:"postal_code" binding:"required"`
		Country       string `form:"country" binding:"required"`
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Make sure country is a valid ISO 3166-1 alpha-2 code
	validCodes := utils.CreateValidCountryCodes()
	country := strings.ToUpper(request.Country)

	if !validCodes[country] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid country code, must be a valid alpha-2 country code",
		})
		return
	}

	// Handle front photo
	frontFile, err := c.FormFile("front_photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Front photo is required",
		})
		return
	}

	// Handle back photo
	backFile, err := c.FormFile("back_photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Back photo is required",
		})
		return
	}

	// Validate front photo
	if frontFile.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Front photo exceeds 5MB limit",
		})
		return
	}

	// Validate back photo
	if backFile.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Back photo exceeds 5MB limit",
		})
		return
	}

	// Open and validate front photo
	frontSrc, err := frontFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error processing front photo",
		})
		return
	}
	defer frontSrc.Close()

	// Detect front photo MIME type
	frontBuffer := make([]byte, 512)
	_, err = frontSrc.Read(frontBuffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error reading front photo",
		})
		return
	}
	frontMimeType := http.DetectContentType(frontBuffer)
	if !validMimeTypes[frontMimeType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Front photo must be JPG, JPEG, or PNG",
		})
		return
	}

	// Open and validate back photo
	backSrc, err := backFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error processing back photo",
		})
		return
	}
	defer backSrc.Close()

	// Detect back photo MIME type
	backBuffer := make([]byte, 512)
	_, err = backSrc.Read(backBuffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error reading back photo",
		})
		return
	}
	backMimeType := http.DetectContentType(backBuffer)
	if !validMimeTypes[backMimeType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Back photo must be JPG, JPEG, or PNG",
		})
		return
	}

	// Convert front photo to base64
	frontSrc.Seek(0, 0)
	frontBytes, err := io.ReadAll(frontSrc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error processing front photo",
		})
		return
	}
	frontBase64 := fmt.Sprintf("data:%s;base64,%s", frontMimeType, base64.StdEncoding.EncodeToString(frontBytes))

	// Convert back photo to base64
	backSrc.Seek(0, 0)
	backBytes, err := io.ReadAll(backSrc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error processing back photo",
		})
		return
	}
	backBase64 := fmt.Sprintf("data:%s;base64,%s", backMimeType, base64.StdEncoding.EncodeToString(backBytes))

	// Check if the user exists in the database
	user, ok := models.FindUserByEmail(request.Email)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	log.Printf("Found user data: %+v", user)

	// Check if the user has already submitted a KYC request
	_, err = models.GetKYCByUserID(user.ID)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "KYC request already submitted",
		})
		return
	}

	// Construct KYC model object
	kyc := &models.KYC{
		UserID:        user.ID,
		IDType:        request.IDType,
		IssueDate:     request.IssueDate,
		ExpiryDate:    request.ExpiryDate,
		Status:        models.Pending,
		IdNumber:      request.IdNumber,
		TaxId:         request.TaxId,
		DateOfBirth:   request.DateOfBirth,
		Phone:         request.Phone,
		StreetAddress: request.StreetAddress,
		City:          request.City,
		State:         request.State,
		PostalCode:    request.PostalCode,
		Country:       request.Country,
	}

	// Construct KYCdata model object
	kycData := &models.KYCData{
		UserID:     user.ID,
		FrontPhoto: frontBase64,
		BackPhoto:  backBase64,
	}

	if err := kyc.CreateKYC(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := kycData.CreateKYCData(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "KYC request submitted successfully",
	})
}

func UpdateKYC(c *gin.Context) {
	// Max file size: 5MB
	const maxSize = 5 << 20

	// Accepted MIME types
	validMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}

	// Bind form data
	var request struct {
		Email         string `form:"email" binding:"required,email"`
		IDType        string `form:"id_type" binding:"required"`
		IssueDate     string `form:"issue_date" binding:"required"`
		ExpiryDate    string `form:"expiry_date" binding:"required"`
		TaxId         string `form:"tax_id" binding:"required"`
		IdNumber      string `form:"id_number" binding:"required"`
		DateOfBirth   string `form:"date_of_birth" binding:"required"`
		Phone         string `form:"phone" binding:"required"`
		StreetAddress string `form:"street_address" binding:"required"`
		City          string `form:"city" binding:"required"`
		State         string `form:"state" binding:"required"`
		PostalCode    string `form:"postal_code" binding:"required"`
		Country       string `form:"country" binding:"required"`
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if the user exists in the database
	user, ok := models.FindUserByEmail(request.Email)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Check if the user has already submitted a KYC request
	existingKyc, err := models.GetKYCByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KYC Not Found",
		})
		return
	}

	// Find corresponding KYC Data
	existingKycData, err := models.GetKYCDataByUserId(user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KYC Data Not Found",
		})
		return
	}

	// Initialize update data with existing values
	frontBase64 := existingKycData.FrontPhoto
	backBase64 := existingKycData.BackPhoto

	// Process front photo if provided
	if frontFile, err := c.FormFile("front_photo"); err == nil {
		if frontFile.Size > maxSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Front photo exceeds 5MB limit",
			})
			return
		}

		frontSrc, err := frontFile.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error processing front photo",
			})
			return
		}
		defer frontSrc.Close()

		// Detect MIME type
		frontBuffer := make([]byte, 512)
		_, err = frontSrc.Read(frontBuffer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error reading front photo",
			})
			return
		}
		frontMimeType := http.DetectContentType(frontBuffer)
		if !validMimeTypes[frontMimeType] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Front photo must be JPG, JPEG, or PNG",
			})
			return
		}

		// Convert to base64
		frontSrc.Seek(0, 0)
		frontBytes, err := io.ReadAll(frontSrc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error processing front photo",
			})
			return
		}
		frontBase64 = fmt.Sprintf("data:%s;base64,%s", frontMimeType, base64.StdEncoding.EncodeToString(frontBytes))
	}

	// Process back photo if provided
	if backFile, err := c.FormFile("back_photo"); err == nil {
		if backFile.Size > maxSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Back photo exceeds 5MB limit",
			})
			return
		}

		backSrc, err := backFile.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error processing back photo",
			})
			return
		}
		defer backSrc.Close()

		// Detect MIME type
		backBuffer := make([]byte, 512)
		_, err = backSrc.Read(backBuffer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error reading back photo",
			})
			return
		}
		backMimeType := http.DetectContentType(backBuffer)
		if !validMimeTypes[backMimeType] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Back photo must be JPG, JPEG, or PNG",
			})
			return
		}

		// Convert to base64
		backSrc.Seek(0, 0)
		backBytes, err := io.ReadAll(backSrc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error processing back photo",
			})
			return
		}
		backBase64 = fmt.Sprintf("data:%s;base64,%s", backMimeType, base64.StdEncoding.EncodeToString(backBytes))
	}

	// Construct updated KYC model
	updatedKyc := models.KYCRequest{
		IDType:        request.IDType,
		IssueDate:     request.IssueDate,
		ExpiryDate:    request.ExpiryDate,
		IdNumber:      request.IdNumber,
		TaxId:         request.TaxId,
		DateOfBirth:   request.DateOfBirth,
		Phone:         request.Phone,
		StreetAddress: request.StreetAddress,
		City:          request.City,
		State:         request.State,
		PostalCode:    request.PostalCode,
		Country:       request.Country,
	}

	// Construct updated KYCData model
	updatedKycData := models.KYCDataRequest{
		FrontPhoto: frontBase64,
		BackPhoto:  backBase64,
	}

	if err := existingKyc.UpdateKYC(updatedKyc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := existingKycData.UpdateKYCData(existingKyc.Status, updatedKycData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "KYC request updated successfully",
	})
}

func ApproveKYC(c *gin.Context) {
	id := c.Param("id")
	idUint64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format - must be a positive integer",
		})
		return
	}
	idUint := uint(idUint64)

	// Find the KYC
	existingKyc, err := models.GetKYCByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KYC Not Found",
		})
		return
	}

	// Find kyc data
	existingKycData, err := models.GetKYCDataByUserId(existingKyc.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KYC Data Not Found",
		})
		return
	}

	// Find the user this kyc belongs to
	user, err := models.GetUserByID(existingKyc.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// if user is verified, return error
	if user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User is already verified",
		})
		return
	}

	// make call to Borderless to create individual identity
	borderless := borderless.NewBorderless()

	borderlessIdentityAddress := models.BorderlessIdentityAddress{
		Street1:    existingKyc.StreetAddress,
		City:       existingKyc.City,
		State:      existingKyc.State,
		PostalCode: existingKyc.PostalCode,
		Country:    existingKyc.Country,
	}

	borderlessIdentity := models.BorderlessIdentity{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Phone:       existingKyc.Phone,
		TaxId:       existingKyc.TaxId,
		DateOfBirth: existingKyc.DateOfBirth,
		Address:     borderlessIdentityAddress,
	}

	// first we try to check if the customer already has an identity
	response, err := borderless.GetCustomerIdentity(borderlessIdentity.Email, borderlessIdentity.LastName)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	var borderlessID string

	// check length of data field in response
	// if customer has no identity then proceed to create one
	if len(response["data"].([]interface{})) < 1 {
		response, err := borderless.CreateCustomerIdentity(borderlessIdentity)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": err.Error(),
			})
			return
		}

		borderlessID = response["id"].(string)
	} else {
		borderlessID = response["data"].([]interface{})[0].(map[string]interface{})["id"].(string)
	}

	// Upload documents to Borderless Identity
	response, err = borderless.UploadCustomerIdentityDocument(borderlessID, *existingKyc, *existingKycData)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Extract id field from the response
	borderlessID, ok := response["id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid response format",
		})
		return
	}

	// set user to isVerified only if not already verified
	user.IsVerified = true
	if err := user.UpdateUserWithErrors(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update the KYC request
	if err := existingKyc.ApproveKYC(borderlessID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "KYC request approved successfully",
	})

}

func RejectKYC(c *gin.Context) {
	id := c.Param("id")
	idUint64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format - must be a positive integer",
		})
		return
	}
	idUint := uint(idUint64)

	var request struct {
		RejectionReason string `json:"rejection_reason"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Find the KYC
	existingKyc, err := models.GetKYCByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KYC Not Found",
		})
		return
	}

	// Update the KYC request
	if err := existingKyc.RejectKYC(request.RejectionReason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "KYC request rejected successfully",
	})

}

func DeleteKYC(c *gin.Context) {
	id := c.Param("id")
	idUint64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format - must be a positive integer",
		})
		return
	}
	idUint := uint(idUint64)

	// Check if the user has already submitted a KYC request
	existingKyc, err := models.GetKYCByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KYC Not Found",
		})
		return
	}

	// Find kyc data
	existingKycData, err := models.GetKYCDataByUserId(existingKyc.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KYC Data Not Found",
		})
		return
	}

	// Delete the KYC Data
	if err := existingKycData.DeleteKYCData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete KYC data",
		})
		return
	}

	// Delete the KYC request
	if err := existingKyc.DeleteKYC(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete KYC request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "KYC request deleted successfully",
	})
}
