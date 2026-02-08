package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OAuthAccount struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Provider       string     `gorm:"type:varchar(50);not null;index" json:"provider"`
	ProviderUserID string     `gorm:"not null;size:255;index" json:"provider_user_id"`
	ProviderEmail  string     `gorm:"size:255" json:"provider_email"`
	AccessToken    string     `gorm:"size:500" json:"-"`
	RefreshToken   string     `gorm:"size:500" json:"-"`
	ExpiresAt      *time.Time `gorm:"index" json:"expires_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (*OAuthAccount) TableName() string {
	return "oauth_accounts"
}

func (o *OAuthAccount) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}

	o.Provider = strings.ToLower(o.Provider)

	return nil
}

func (o *OAuthAccount) BeforeUpdate(tx *gorm.DB) error {
	o.Provider = strings.ToLower(o.Provider)
	return nil
}

func (o *OAuthAccount) IsTokenExpired() bool {
	if o.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*o.ExpiresAt)
}

func (o *OAuthAccount) NeedsRefresh() bool {
	if o.ExpiresAt == nil {
		return false
	}
	return time.Now().Add(5 * time.Minute).After(*o.ExpiresAt)
}
