package community_infrastructure

import (
	"encoding/hex"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	cd "github.com/214alphadev/community-bl"
	vo "github.com/214alphadev/community-bl/value_objects"
	"time"
)

type Member struct {
	ID                    string    `gorm:"unique;not null;primary_key"`
	CreatedAt             time.Time `gorm:"non null"`
	Admin                 bool      `gorm:"non null"`
	EmailAddress          string    `gorm:"unique;non null"`
	Username              string    `gorm:"unique;non null"`
	VerifiedEmailAddress  bool      `gorm:"non null"`
	Metadata              string    `gorm:"non null"`
	MemberAccessPublicKey *string
	AccessTokenID         *string
	Verified              bool `gorm:"non null"`
}
type Metadata struct {
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	ProfileImage *string `json:"profile_image,omitempty"`
}

type MemberRepository struct {
	db *gorm.DB
}

var mapMemberToEntity = func(member Member) (*cd.MemberEntity, error) {

	username, err := vo.NewUsername(member.Username)
	if err != nil {
		return nil, err
	}

	parsedEmailAddress, err := vo.NewEmailAddress(member.EmailAddress)
	if err != nil {
		return nil, err
	}

	memberID, err := uuid.FromString(member.ID)
	if err != nil {
		return nil, err
	}

	metadata := Metadata{}
	if err := json.Unmarshal([]byte(member.Metadata), &metadata); err != nil {
		return nil, err
	}
	metadataProperName, err := vo.NewProperName(metadata.FirstName, metadata.LastName)
	if err != nil {
		return nil, err
	}

	memberEntity := &cd.MemberEntity{
		ID:                   memberID,
		VerifiedEmailAddress: member.VerifiedEmailAddress,
		Username:             username,
		EmailAddress:         parsedEmailAddress,
		Metadata: cd.MetadataEntity{
			ProperName: metadataProperName,
		},
		CreatedAt: member.CreatedAt,
		Admin:     member.Admin,
		Verified:  member.Verified,
	}

	if metadata.ProfileImage != nil {
		metadataProfileImage, err := vo.NewBase64String(*metadata.ProfileImage)
		if err != nil {
			return nil, err
		}
		memberEntity.Metadata.ProfileImage = &metadataProfileImage
	}

	if member.MemberAccessPublicKey != nil {
		memberAccessPublicKeyBytes, err := hex.DecodeString(*member.MemberAccessPublicKey)
		if err != nil {
			return nil, err
		}
		memberAccessPublicKey, err := vo.NewMemberAccessPublicKey(memberAccessPublicKeyBytes)
		if err != nil {
			return nil, err
		}
		memberEntity.MemberAccessPublicKey = &memberAccessPublicKey
	}

	if member.AccessTokenID != nil {
		accessTokenID, err := uuid.FromString(*member.AccessTokenID)
		if err != nil {
			return nil, err
		}
		memberEntity.AccessTokenID = &accessTokenID
	}

	return memberEntity, nil

}

func (r *MemberRepository) FetchByID(memberID cd.MemberIdentifier) (*cd.MemberEntity, error) {

	m := &Member{}

	err := r.db.First(m, "id = ?", memberID.String()).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return nil, nil
	case nil:
		return mapMemberToEntity(*m)
	default:
		return nil, err
	}

}

func (r *MemberRepository) Save(memberEntity cd.MemberEntity) error {

	metadata := Metadata{
		FirstName: memberEntity.Metadata.ProperName.FirstName(),
		LastName:  memberEntity.Metadata.ProperName.LastName(),
	}

	if memberEntity.Metadata.ProfileImage != nil {
		i := memberEntity.Metadata.ProfileImage.String()
		metadata.ProfileImage = &i
	}

	metadataStr, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	m := &Member{
		ID:                   memberEntity.ID.String(),
		CreatedAt:            memberEntity.CreatedAt,
		Admin:                memberEntity.Admin,
		Verified:             memberEntity.Verified,
		VerifiedEmailAddress: memberEntity.VerifiedEmailAddress,
		Username:             memberEntity.Username.String(),
		EmailAddress:         memberEntity.EmailAddress.String(),
		Metadata:             string(metadataStr),
	}

	if memberEntity.AccessTokenID != nil {
		i := memberEntity.AccessTokenID.String()
		m.AccessTokenID = &i
	}

	if memberEntity.MemberAccessPublicKey != nil {
		i := hex.EncodeToString(memberEntity.MemberAccessPublicKey.Key())
		m.MemberAccessPublicKey = &i
	}

	return r.db.Save(m).Error
}

func (r *MemberRepository) IsUsernameTaken(username vo.Username) (bool, error) {

	err := r.db.Find(&Member{}, "username = ?", username.String()).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}

}

func (r *MemberRepository) IsEmailAddressTaken(emailAddress vo.EmailAddress) (bool, error) {

	err := r.db.Find(&Member{}, "email_address = ?", emailAddress.String()).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}

}

func (r *MemberRepository) FetchByEmailAddress(emailAddress vo.EmailAddress) (*cd.MemberEntity, error) {

	member := &Member{}

	err := r.db.Find(member, "email_address = ?", emailAddress.String()).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return nil, nil
	case nil:
		return mapMemberToEntity(*member)
	default:
		return nil, err
	}

}

func NewMemberRepository(db *gorm.DB) (*MemberRepository, error) {

	if err := db.AutoMigrate(&Member{}).Error; err != nil {
		return nil, err
	}

	return &MemberRepository{
		db: db,
	}, nil

}
