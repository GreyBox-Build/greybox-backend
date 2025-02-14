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
	Email      string    `json:"email"`
	IDType     string    `json:"id_type"`
	IssueDate  string    `json:"issue_date"`
	ExpiryDate string    `json:"expiry_date"`
	FrontPhoto string    `json:"front_photo"`
	BackPhoto  string    `json:"back_photo"`
	Status     KYCStatus `gorm:"default:Pending" json:"status"`
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
	ID              uint      `json:"id"`
	IDType          string    `json:"id_type"`
	IssueDate       string    `json:"issue_date"`
	ExpiryDate      string    `json:"expiry_date"`
	Status          KYCStatus `gorm:"default:Pending" json:"status"`
	RejectionReason string    `json:"rejection_reason"`
	CreatedAt       time.Time `json:"created_at"`
	ApprovedAt      time.Time `json:"approved_at"`
	RejectedAt      time.Time `json:"rejected_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
