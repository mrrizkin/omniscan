package models

import (
	"time"

	"gorm.io/gorm"
)

type EStatement struct {
	ID               uint               `json:"id"            gorm:"primary_key"`
	CreatedAt        *time.Time         `json:"created_at"`
	UpdatedAt        *time.Time         `json:"updated_at"`
	DeletedAt        gorm.DeletedAt     `json:"deleted_at"    gorm:"index"`
	Filename         string             `json:"filename"      gorm:"index"`
	Bank             string             `json:"bank"`
	Produk           string             `json:"produk"`
	Rekening         string             `json:"rekening"`
	Periode          string             `json:"periode"`
	Expired          *time.Time         `json:"expired"       gorm:"index"`
	EStatementDetail []EStatementDetail `json:"e_statement_detail" gorm:"foreignKey:EStatementID;references:ID"`
}
