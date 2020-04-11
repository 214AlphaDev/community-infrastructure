package community_infrastructure

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	cd "github.com/214alphadev/community-bl"
	vo "github.com/214alphadev/community-bl/value_objects"
)

type ConfirmationCode struct {
	ID               string `gorm:"unique;not null;primary_key"`
	Member           string `gorm:"not null"`
	EmailAddress     string `gorm:"not null"`
	ConfirmationCode string `gorm:"not null"`
	IssuedAt         int64  `gorm:"not null"`
	Used             bool   `gorm:"not null"`
}

type ConfirmationCodeRepository struct {
	db *gorm.DB
}

var confCodeModelToEntity = func(cc ConfirmationCode) (*cd.ConfirmationCode, error) {

	id, err := uuid.FromString(cc.ID)
	if err != nil {
		return nil, err
	}

	member, err := uuid.FromString(cc.Member)
	if err != nil {
		return nil, err
	}

	emailAddress, err := vo.NewEmailAddress(cc.EmailAddress)
	if err != nil {
		return nil, err
	}

	confirmationCode, err := vo.NewConfirmationCode(cc.ConfirmationCode)
	if err != nil {
		return nil, err
	}

	return &cd.ConfirmationCode{
		ID:               id,
		MemberIdentifier: member,
		EmailAddress:     emailAddress,
		ConfirmationCode: confirmationCode,
		IssuedAt:         cc.IssuedAt,
		Used:             cc.Used,
	}, nil

}

func (r *ConfirmationCodeRepository) Fetch(emailAddress vo.EmailAddress, confirmationCode vo.ConfirmationCode) (*cd.ConfirmationCode, error) {

	cc := &ConfirmationCode{}

	err := r.db.Find(cc, "email_address = ? AND confirmation_code = ?", emailAddress.String(), confirmationCode.String()).Error

	switch err {
	case nil:
		return confCodeModelToEntity(*cc)
	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		return nil, err
	}

}

func (r *ConfirmationCodeRepository) Save(cc *cd.ConfirmationCode) error {
	return r.db.Save(&ConfirmationCode{
		ID:               cc.ID.String(),
		Member:           cc.MemberIdentifier.String(),
		EmailAddress:     cc.EmailAddress.String(),
		ConfirmationCode: cc.ConfirmationCode.String(),
		IssuedAt:         cc.IssuedAt,
		Used:             cc.Used,
	}).Error
}

func (r *ConfirmationCodeRepository) Last(emailAddress vo.EmailAddress) (*cd.ConfirmationCode, error) {

	cc := &ConfirmationCode{}

	err := r.db.Find(cc, "email_address = ?", emailAddress.String()).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return nil, nil
	case nil:
		return confCodeModelToEntity(*cc)
	default:
		return nil, err
	}

}

func NewConfirmationCodeRepository(db *gorm.DB) (*ConfirmationCodeRepository, error) {

	if err := db.AutoMigrate(&ConfirmationCode{}).Error; err != nil {
		return nil, err
	}

	return &ConfirmationCodeRepository{
		db: db,
	}, nil
}
