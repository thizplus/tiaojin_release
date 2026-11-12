package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BannerStatus string

const (
	BannerStatusPublish BannerStatus = "publish"
	BannerStatusDraft   BannerStatus = "draft"
)

type BannerPosition string

const (
	BannerPositionInternal BannerPosition = "internal"
	BannerPositionPublic   BannerPosition = "public"
)

type Banner struct {
	ID           uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title        string         `gorm:"not null;size:200" json:"title"`
	Image        string         `gorm:"not null;size:500;column:img" json:"image"`
	LinkURL      *string        `gorm:"size:500" json:"link_url"`
	Alt          *string        `gorm:"size:200" json:"alt"`
	DisplayOrder int            `gorm:"not null;default:0" json:"display_order"`
	Position     BannerPosition `gorm:"type:varchar(20);not null;default:'internal'" json:"position"`
	Status       BannerStatus   `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Banner) TableName() string {
	return "banners"
}

func (b *Banner) IsPublished() bool {
	return b.Status == BannerStatusPublish
}

func (b *Banner) IsDraft() bool {
	return b.Status == BannerStatusDraft
}