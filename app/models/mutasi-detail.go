package models

import (
	"time"

	"gorm.io/gorm"
)

type MutasiDetail struct {
	ID              uint           `json:"id"                         gorm:"primary_key"`
	CreatedAt       *time.Time     `json:"created_at"`
	UpdatedAt       *time.Time     `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at"                 gorm:"index"`
	Date            time.Time      `json:"date,omitempty"             gorm:"index"`
	MutasiID        uint           `json:"mutasi_id"`
	Description1    string         `json:"description1,omitempty"`
	Description2    string         `json:"description2,omitempty"`
	Branch          string         `json:"branch,omitempty"`
	Change          float64        `json:"change,omitempty"           gorm:"type:numeric"`
	TransactionType string         `json:"transaction_type,omitempty" gorm:"index"`
	Balance         float64        `json:"balance,omitempty"          gorm:"type:numeric"`
	Mutasi          *Mutasi        `json:"mutasi,omitempty"           gorm:"foreignKey:MutasiID;references:ID"`
}
