package models

import (
	"time"

	"gorm.io/gorm"
)

type Caisse struct {
	UUID      string `gorm:"type:varchar(255);primary_key;not null" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	AppartmentUUID string     `gorm:"type:varchar(255);not null" json:"appartment_uuid"`
	Appartment     Appartment `gorm:"foreignKey:AppartmentUUID;references:UUID" json:"appartment"`

	Type      string  `gorm:"type:varchar(20);not null" json:"type"` // Entrees et Sorties (Income/Expense)
	DeviceCDF float64 `gorm:"default:0" json:"device_cdf"`
	DeviceUSD float64 `gorm:"default:0" json:"device_usd"`

	Signature string `gorm:"not null" json:"signature"` // Pour savoir qui q fait des entrees et des sorties
}

// ValidateType validates that the Type field contains only allowed values
func (c *Caisse) ValidateType() bool {
	return c.Type == "Income" || c.Type == "Expense"
}
