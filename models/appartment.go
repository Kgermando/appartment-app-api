package models

import (
	"time"

	"gorm.io/gorm"
)

type Appartment struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name   string `gorm:"not null" json:"name"`   // Locateur Ex. okapi
	Number string `gorm:"not null" json:"number"` // Numero appartement Ex. 1201

	// Caractéristiques physiques
	Surface   float64 `gorm:"default:0" json:"surface"`       // Surface en m²
	Rooms     int     `gorm:"default:1" json:"rooms"`         // Nombre de chambres
	Bathrooms int     `gorm:"default:1" json:"bathrooms"`     // Nombre de salles de bain
	Balcony   bool    `gorm:"default:false" json:"balcony"`   // Présence d'un balcon
	Furnished bool    `gorm:"default:false" json:"furnished"` // Meublé ou non

	// Informations financières
	MonthlyRent   float64 `gorm:"not null;default:0" json:"monthly_rent" validate:"required,gt=0"`     // Loyer mensuel
	GarantieMonth float64 `gorm:"not null;default:2" json:"garantie_month" validate:"required,gt=0"`   // Nombre de mois de garantie
	Garantie      float64 `gorm:"not null;default:0" json:"garantie_montant" validate:"required,gt=0"` // Montant de la garantie

	// Date d'échéance
	Echeance time.Time `json:"echeance"` // Date de paiement loyer

	// Statut et disponibilité
	Status    string `gorm:"default:'available'" json:"status"` // available, occupied, maintenance, unavailable

	// Gestionnaire/Agent responsable
	ManagerUUID string `gorm:"type:varchar(255)" json:"manager_uuid"`
	Manager     User   `gorm:"foreignKey:ManagerUUID;references:UUID" json:"manager"`

	// Relations inverses
	Caisses []Caisse `gorm:"foreignKey:AppartmentUUID;references:UUID" json:"caisses,omitempty"`
}
