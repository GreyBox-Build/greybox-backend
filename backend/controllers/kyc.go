package controllers

import (
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils/tokens"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strconv"

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

	c.JSON(http.StatusOK, gin.H{
		"kyc": kycResponse,
	})
}

func GetKYCS(c *gin.Context) {
	var request serializers.KYCFilterRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
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
		FrontPhoto    string `form:"front_photo" binding:"required"`
		BackPhoto     string `form:"back_photo" binding:"required"`
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
	frontBase64 := base64.StdEncoding.EncodeToString(frontBytes)

	// Convert back photo to base64
	backSrc.Seek(0, 0)
	backBytes, err := io.ReadAll(backSrc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error processing back photo",
		})
		return
	}
	backBase64 := base64.StdEncoding.EncodeToString(backBytes)

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
		FrontPhoto:    frontBase64,
		BackPhoto:     backBase64,
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

	if err := kyc.CreateKYC(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
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
		FrontPhoto    string `form:"front_photo" binding:"required"`
		BackPhoto     string `form:"back_photo" binding:"required"`
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

	// Initialize update data with existing values
	frontBase64 := existingKyc.FrontPhoto
	backBase64 := existingKyc.BackPhoto

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
		frontBase64 = base64.StdEncoding.EncodeToString(frontBytes)
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
		backBase64 = base64.StdEncoding.EncodeToString(backBytes)
	}

	// Construct updated KYC model
	updatedKyc := models.KYCRequest{
		IDType:        request.IDType,
		IssueDate:     request.IssueDate,
		ExpiryDate:    request.ExpiryDate,
		FrontPhoto:    frontBase64,
		BackPhoto:     backBase64,
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

	if err := existingKyc.UpdateKYC(updatedKyc); err != nil {
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

	// Find the user this kyc belongs to
	user, err := models.GetUserByID(existingKyc.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// make call to Borderless to create individual identity
	borderless := apis.NewBorderless()

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

	response, err := borderless.CreateCustomerIdentity(borderlessIdentity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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

	// Upload documents to Borderless Identity
	response, err = borderless.UploadCustomerIdentityDocument(borderlessID, *existingKyc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Extract id field from the response
	borderlessID, ok = response["id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid response format",
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

	// Update the KYC request
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
