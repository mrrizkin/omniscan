package models

import (
	"time"

	"gorm.io/gorm"
)

type EStatementDetail struct {
	ID              uint           `json:"id"                         gorm:"primary_key"`
	CreatedAt       *time.Time     `json:"created_at"`
	UpdatedAt       *time.Time     `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at"                 gorm:"index"`
	Date            time.Time      `json:"date,omitempty"             gorm:"index"`
	EStatementID    uint           `json:"e_statement_id,omitempty"   gorm:"index"`
	Description1    string         `json:"description1,omitempty"`
	Description2    string         `json:"description2,omitempty"`
	Branch          string         `json:"branch,omitempty"`
	Change          float64        `json:"change,omitempty"           gorm:"type:numeric"`
	TransactionType string         `json:"transaction_type,omitempty" gorm:"index"`
	Balance         float64        `json:"balance,omitempty"          gorm:"type:numeric"`
	EStatement      *EStatement    `json:"e_statement,omitempty"     gorm:"foreignKey:EStatementID;references:ID"`
}
