package models

import "time"

// PersonalInfo represents personal information provided by guardians
type PersonalInfo struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Content     string    `json:"content"`      // 보호자가 입력한 텍스트
	Category    string    `json:"category"`     // e.g., "medical", "contact", "emergency", "allergy"
	Importance  string    `json:"importance"`   // "high", "medium", "low"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PersonalInfoCreateRequest represents a request to create personal info
type PersonalInfoCreateRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Category   string `json:"category" binding:"required"`
	Importance string `json:"importance" binding:"required,oneof=high medium low"`
}

// PersonalInfoUpdateRequest represents a request to update personal info
type PersonalInfoUpdateRequest struct {
	Content    string `json:"content"`
	Category   string `json:"category"`
	Importance string `json:"importance" binding:"omitempty,oneof=high medium low"`
}

// PersonalInfoResponse represents a personal info response
type PersonalInfoResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Content    string `json:"content"`
	Category   string `json:"category"`
	Importance string `json:"importance"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// PersonalInfoListResponse represents a list of personal info items
type PersonalInfoListResponse struct {
	Items       []PersonalInfoResponse `json:"items"`
	Total       int                    `json:"total"`
	UserID      string                 `json:"user_id"`
}