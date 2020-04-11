package community_infrastructure

import (
	"encoding/hex"

	"github.com/jinzhu/gorm"
	vo "github.com/214alphadev/community-bl/value_objects"
)

type MemberAccessPublicKey struct {
	Key string `gorm:"unique;non null"`
}

type MemberAccessPublicKeyRepository struct {
	db *gorm.DB
}

func (r *MemberAccessPublicKeyRepository) AlreadyUsed(memberAccessPublicKey vo.MemberAccessPublicKey) (bool, error) {

	err := r.db.Find(&MemberAccessPublicKey{}, "key = ?", hex.EncodeToString(memberAccessPublicKey.Key())).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}

}

func (r *MemberAccessPublicKeyRepository) Save(memberAccessPublicKey vo.MemberAccessPublicKey) error {
	return r.db.Create(&MemberAccessPublicKey{
		Key: hex.EncodeToString(memberAccessPublicKey.Key()),
	}).Error
}

func NewMemberAccessPublicKeyRepository(db *gorm.DB) (*MemberAccessPublicKeyRepository, error) {

	if err := db.AutoMigrate(&MemberAccessPublicKey{}).Error; err != nil {
		return nil, err
	}

	return &MemberAccessPublicKeyRepository{
		db: db,
	}, nil
}
