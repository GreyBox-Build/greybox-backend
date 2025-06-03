package models

import (
	"backend/serializers"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type KYCStatus string

// Define constants of the possible values of the enum
const (
	Pending  KYCStatus = "Pending"
	Approved KYCStatus = "Approved"
	Rejected KYCStatus = "Rejected"
)

type KYC struct {
	gorm.Model
	UserID               uint      `gorm:"uniqueIndex" json:"user_id"`
	User                 User      `gorm:"foreignKey:UserID" json:"user"`
	IDType               string    `json:"id_type"`
	TaxId                string    `json:"tax_id"`
	IdNumber             string    `json:"id_number"`
	DateOfBirth          string    `json:"date_of_birth"`
	IssueDate            string    `json:"issue_date"`
	ExpiryDate           string    `json:"expiry_date"`
	Status               KYCStatus `gorm:"default:Pending" json:"status"`
	Phone                string    `json:"phone"`
	StreetAddress        string    `json:"street_address"`
	BorderlessIdentityId string    `json:"borderless_identity_id"`
	City                 string    `json:"city"`
	State                string    `json:"state"`
	PostalCode           string    `json:"postal_code"`
	Country              string    `json:"country"`
	RejectionReason      string    `json:"rejection_reason"`
	CreatedAt            time.Time `json:"created_at"`
	ApprovedAt           time.Time `json:"approved_at"`
	RejectedAt           time.Time `json:"rejected_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type KYCWithData struct {
	Kyc        KYC
	FrontPhoto string `json:"front_photo"`
	BackPhoto  string `json:"back_photo"`
}

type KYCData struct {
	gorm.Model
	UserID     uint   `gorm:"uniqueIndex" json:"user_id"`
	FrontPhoto string `json:"front_photo"`
	BackPhoto  string `json:"back_photo"`
}

type KYCRequest struct {
	UserID        uint      `json:"user_id"`
	IDType        string    `json:"id_type"`
	IdNumber      string    `json:"id_number"`
	IssueDate     string    `json:"issue_date"`
	ExpiryDate    string    `json:"expiry_date"`
	Status        KYCStatus `gorm:"default:Pending" json:"status"`
	DateOfBirth   string    `json:"date_of_birth"`
	TaxId         string    `json:"tax_id"`
	Phone         string    `json:"phone"`
	StreetAddress string    `json:"street_address"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	PostalCode    string    `json:"postal_code"`
	Country       string    `json:"country"`
}

type KYCDataRequest struct {
	FrontPhoto string `json:"front_photo"`
	BackPhoto  string `json:"back_photo"`
}

type ApproveOrDeleteKYC struct {
	UserID uint `json:"user_id"`
}

type RejectKYC struct {
	UserID          uint   `json:"user_id"`
	RejectionReason string `json:"rejection_reason"`
}

// Borderless Identity Address Structure
type BorderlessIdentityAddress struct {
	Street1    string `json:"street1"`
	Street2    string `json:"street2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`
}

// Borderless Identity Structure
type BorderlessIdentity struct {
	FirstName   string                    `json:"firstName"`
	LastName    string                    `json:"lastName"`
	MiddleName  string                    `json:"middleName,omitempty"`
	TaxId       string                    `json:"taxId"`
	DateOfBirth string                    `json:"dateOfBirth"`
	Email       string                    `json:"email"`
	Phone       string                    `json:"phone,omitempty"`
	Address     BorderlessIdentityAddress `json:"address"`
}

func (kyc *KYC) BeforeCreate(tx *gorm.DB) (err error) {
	kyc.CreatedAt = time.Now()
	kyc.UpdatedAt = time.Now()
	return
}

func (kyc *KYC) BeforeUpdate(tx *gorm.DB) (err error) {
	kyc.UpdatedAt = time.Now()
	return
}

func (kyc *KYC) BeforeSave(tx *gorm.DB) (err error) {
	kyc.UpdatedAt = time.Now()
	return
}

func (kyc *KYC) CreateKYC() error {
	err := db.Unscoped().Create(&kyc).Error // Use Unscoped() to include soft deleted records
	if err != nil {
		return err
	}

	return nil
}

func (kycData *KYCData) CreateKYCData() error {
	err := db.Unscoped().Create(&kycData).Error
	if err != nil {
		return err
	}

	return nil
}

func (kyc *KYC) UpdateKYC(data KYCRequest) error {
	// Only update if status is rejected
	if kyc.Status != Rejected {
		return fmt.Errorf("KYC can only be updated when status is rejected")
	}

	kyc.IDType = data.IDType
	kyc.IdNumber = data.IdNumber
	kyc.IssueDate = data.IssueDate
	kyc.ExpiryDate = data.ExpiryDate
	kyc.TaxId = data.TaxId
	kyc.Phone = data.Phone
	kyc.StreetAddress = data.StreetAddress
	kyc.City = data.City
	kyc.State = data.State
	kyc.PostalCode = data.PostalCode
	kyc.Country = data.Country
	kyc.DateOfBirth = data.DateOfBirth
	kyc.Status = Pending // set status to pending during update

	return db.Save(kyc).Error
}

func (kycData *KYCData) UpdateKYCData(status KYCStatus, data KYCDataRequest) error {
	// Only update if status of kyc is pending
	if status != Rejected {
		return fmt.Errorf("KYC can only be updated when status is rejected")
	}

	kycData.FrontPhoto = data.FrontPhoto
	kycData.BackPhoto = data.BackPhoto

	return db.Save(kycData).Error
}

func (kyc *KYC) DeleteKYC() error {
	err := db.Where("id = ?", kyc.ID).Delete(kyc).Error
	if err != nil {
		return err
	}
	return nil
}

func (kycData *KYCData) DeleteKYCData() error {
	err := db.Where("id = ?", kycData.ID).Delete(kycData).Error
	if err != nil {
		return err
	}

	return nil
}

func (kyc *KYC) ApproveKYC(borderlessIdentityId string) error {
	// Only update if status is pending
	if kyc.Status != Pending {
		return fmt.Errorf("KYC can only be approved when status is pending")
	}

	kyc.Status = Approved
	kyc.BorderlessIdentityId = borderlessIdentityId
	kyc.ApprovedAt = time.Now()
	return db.Save(kyc).Error
}

func (k *KYC) RejectKYC(rejectionReason string) error {
	// Only update if status is pending
	if k.Status != Pending {
		return fmt.Errorf("KYC can only be rejected when status is pending")
	}

	k.Status = Rejected
	k.RejectionReason = rejectionReason
	k.RejectedAt = time.Now()
	return db.Save(k).Error
}

func GetKYCByUserID(id uint) (*KYC, error) {
	var kyc KYC
	err := db.Where("user_id = ?", id).First(&kyc).Error
	return &kyc, err
}

func GetKYCDataByUserId(id uint) (*KYCData, error) {
	var kycData KYCData
	err := db.Where("user_id = ?", id).First(&kycData).Error
	return &kycData, err
}

func GetKYCByID(id uint) (*KYC, error) {
	var kyc KYC
	err := db.Where("id = ?", id).First(&kyc).Error
	return &kyc, err
}

func GetKycByUserIDWithPhotos(id uint) (*KYCWithData, error) {
	var kyc KYC
	err := db.
		Where("user_id = ?", id).
		First(&kyc).Error

	kycData, err := GetKYCDataByUserId(kyc.UserID)
	if err != nil {
		return nil, err
	}

	fullKyc := KYCWithData{
		Kyc:        kyc,
		FrontPhoto: kycData.FrontPhoto,
		BackPhoto:  kycData.BackPhoto,
	}

	return &fullKyc, err
}

func FilterKYC(filter serializers.KYCFilterRequest) ([]KYC, error) {
	var kycs []KYC
	query := db

	// Build the query dynamically based on provided filters
	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.RejectionReason != nil {
		query = query.Where("rejection_reason ILIKE ?", "%"+*filter.RejectionReason+"%")
	}

	if filter.CreatedAt != nil {
		query = query.Where("created_at = ?", *filter.CreatedAt)
	}

	if filter.ApprovedAt != nil {
		query = query.Where("approved_at = ?", *filter.ApprovedAt)
	}

	if filter.RejectedAt != nil {
		query = query.Where("rejected_at = ?", *filter.RejectedAt)
	}

	// Execute the query
	err := query.Find(&kycs).Error
	return kycs, err
}

// Helper function to convert serializers.KYCRequest to models.KYCRequest
func KYCRequestFromSerializer(
	userId uint,
	status KYCStatus,
	kycRequest serializers.KYCRequest,
) KYCRequest {
	return KYCRequest{
		UserID:        userId,
		IdNumber:      kycRequest.IdNumber,
		IDType:        kycRequest.IDType,
		IssueDate:     kycRequest.IssueDate,
		ExpiryDate:    kycRequest.ExpiryDate,
		Status:        status,
		DateOfBirth:   kycRequest.DateOfBirth,
		TaxId:         kycRequest.TaxId,
		Phone:         kycRequest.Phone,
		StreetAddress: kycRequest.StreetAddress,
		City:          kycRequest.City,
		State:         kycRequest.State,
		Country:       kycRequest.Country,
		PostalCode:    kycRequest.PostalCode,
	}
}
