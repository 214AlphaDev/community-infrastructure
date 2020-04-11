package community_infrastructure

import (
	"github.com/jinzhu/gorm"
	cd "github.com/214alphadev/community-bl"
)

type AccessToken struct {
	ExpiresAt int64  `gorm:"not null"`
	ID        string `gorm:"unique;not null;primary_key"`
	IssuedAt  int64  `gorm:"not null"`
	Subject   string `gorm:"not null"`
}

type AccessTokenRepository struct {
	db *gorm.DB
}

func (r *AccessTokenRepository) Save(accessToken *cd.MemberAccessTokenEntity) error {
	return r.db.Create(&AccessToken{
		ExpiresAt: accessToken.ExpiresAt,
		ID:        accessToken.ID.String(),
		IssuedAt:  accessToken.IssuedAt,
		Subject:   accessToken.Subject.String(),
	}).Error
}

func NewAccessTokenRepository(db *gorm.DB) (*AccessTokenRepository, error) {

	if err := db.AutoMigrate(&AccessToken{}).Error; err != nil {
		return nil, err
	}

	return &AccessTokenRepository{
		db: db,
	}, nil
}
