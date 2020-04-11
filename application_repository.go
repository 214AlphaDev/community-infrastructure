package community_infrastructure

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	cd "github.com/214alphadev/community-bl"
	"time"
)

type Application struct {
	InternalID      uint      `gorm:"AUTO_INCREMENT"`
	ID              string    `gorm:"unique;non null;primary_key"`
	Member          string    `gorm:"non null"`
	ApplicationText string    `gorm:"non null"`
	State           string    `gorm:"non null"`
	RejectionReason string    `gorm:"non null"`
	CreatedAt       time.Time `gorm:"non null"`
	RejectedAt      *time.Time
	ApprovedAt      *time.Time
	RejectedBy      *string
	ApprovedBy      *string
}

type ApplicationRepository struct {
	db *gorm.DB
}

var mapApplicationToEntity = func(application Application) (*cd.ApplicationEntity, error) {

	applicationID, err := uuid.FromString(application.ID)
	if err != nil {
		return nil, err
	}

	memberID, err := uuid.FromString(application.Member)
	if err != nil {
		return nil, err
	}

	applicationEntity := &cd.ApplicationEntity{
		ID:              applicationID,
		MemberID:        memberID,
		ApplicationText: application.ApplicationText,
		State:           cd.ApplicationState(application.State),
		RejectionReason: application.RejectionReason,
		CreatedAt:       application.CreatedAt,
		RejectedAt:      application.RejectedAt,
		ApprovedAt:      application.ApprovedAt,
	}

	if application.RejectedBy != nil {
		rejectedBy, err := uuid.FromString(*application.RejectedBy)
		if err != nil {
			return nil, err
		}
		applicationEntity.RejectedBy = &rejectedBy
	}

	if application.ApprovedBy != nil {
		approvedBy, err := uuid.FromString(*application.ApprovedBy)
		if err != nil {
			return nil, err
		}
		applicationEntity.ApprovedBy = &approvedBy
	}

	return applicationEntity, nil

}

func (r *ApplicationRepository) FetchLast(member cd.MemberIdentifier) (*cd.ApplicationEntity, error) {

	a := &Application{}

	err := r.db.Find(a, "member = ?", member.String()).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return nil, nil
	case nil:
		return mapApplicationToEntity(*a)
	default:
		return nil, err
	}

}

func (r *ApplicationRepository) Save(applicationEntity cd.ApplicationEntity) error {

	if !applicationEntity.State.Valid() {
		return fmt.Errorf("invalid application state: %s", applicationEntity.State)
	}

	application := &Application{
		ID:              applicationEntity.ID.String(),
		ApplicationText: applicationEntity.ApplicationText,
		State:           string(applicationEntity.State),
		RejectionReason: applicationEntity.RejectionReason,
		CreatedAt:       applicationEntity.CreatedAt,
		RejectedAt:      applicationEntity.RejectedAt,
		ApprovedAt:      applicationEntity.ApprovedAt,
		Member:          applicationEntity.MemberID.String(),
	}

	if applicationEntity.ApprovedBy != nil {
		approvedBy := applicationEntity.ApprovedBy.String()
		application.ApprovedBy = &approvedBy
	}

	if applicationEntity.RejectedBy != nil {
		rejectedBy := applicationEntity.RejectedBy.String()
		application.RejectedBy = &rejectedBy
	}

	return r.db.Save(application).Error

}

func (r *ApplicationRepository) FetchByID(applicationID cd.ApplicationID) (*cd.ApplicationEntity, error) {

	a := &Application{}

	err := r.db.Find(a, "id = ?", applicationID.String()).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return nil, nil
	case nil:
		return mapApplicationToEntity(*a)
	default:
		return nil, err
	}

}

func (r *ApplicationRepository) FetchByQuery(query cd.ApplicationsQuery) ([]cd.ApplicationEntity, error) {

	var mapFilteredApplications = func(applications []Application) ([]cd.ApplicationEntity, error) {

		filteredApplications := []cd.ApplicationEntity{}

		for _, application := range applications {
			a, err := mapApplicationToEntity(application)
			if err != nil {
				return nil, err
			}
			filteredApplications = append(filteredApplications, *a)
		}

		return filteredApplications, nil

	}

	switch query.Position {
	case nil:
		applications := []Application{}
		err := r.db.Where("state = ?", query.State).Limit(query.Next).Find(&applications).Error
		if err != nil {
			return nil, err
		}
		return mapFilteredApplications(applications)
	default:
		applications := []Application{}
		startApplication := Application{}
		if err := r.db.Find(&startApplication, "id = ?", query.Position.String()).Error; err != nil {
			return nil, err
		}
		err := r.db.Where("internal_id > ? AND state = ?", startApplication.InternalID, query.State).Limit(query.Next).Find(&applications).Error
		if err != nil {
			return nil, err
		}
		return mapFilteredApplications(applications)
	}

}

func NewApplicationRepository(db *gorm.DB) (*ApplicationRepository, error) {

	if err := db.AutoMigrate(&Application{}).Error; err != nil {
		return nil, err
	}

	return &ApplicationRepository{
		db: db,
	}, nil

}
