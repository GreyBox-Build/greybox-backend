package serializers

import "time"

type KYCStatus string

// Define constants of the possible values of the enum
const (
	Pending  KYCStatus = "Pending"
	Approved KYCStatus = "Approved"
	Rejected KYCStatus = "Rejected"
)

type KYCRequest struct {
	Email         string    `json:"email"`
	IDType        string    `json:"id_type"`
	IdNumber      string    `json:"id_number"`
	IssueDate     string    `json:"issue_date"`
	ExpiryDate    string    `json:"expiry_date"`
	FrontPhoto    string    `json:"front_photo"`
	BackPhoto     string    `json:"back_photo"`
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

type KYCFilterRequest struct {
	ID              *uint      `json:"id,omitempty"`
	UserID          *uint      `json:"user_id,omitempty"`
	Status          *KYCStatus `json:"status,omitempty"`
	RejectionReason *string    `json:"rejection_reason,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	RejectedAt      *time.Time `json:"rejected_at,omitempty"`
}

type KYC struct {
	ID                   uint      `json:"id"`
	IDType               string    `json:"id_type"`
	TaxId                string    `json:"tax_id"`
	IdNumber             string    `json:"id_number"`
	DateOfBirth          string    `json:"date_of_birth"`
	IssueDate            string    `json:"issue_date"`
	ExpiryDate           string    `json:"expiry_date"`
	FrontPhoto           string    `json:"front_photo"`
	BackPhoto            string    `json:"back_photo"`
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
