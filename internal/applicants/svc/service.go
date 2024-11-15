// service/application_service.go
package service

import (
	"errors"

	appEntity "github.com/ghulammuzz/backend-parkerin/internal/applicants/entity"
	appRepo "github.com/ghulammuzz/backend-parkerin/internal/applicants/repo"
	storeRepo "github.com/ghulammuzz/backend-parkerin/internal/store/repo"
	userRepo "github.com/ghulammuzz/backend-parkerin/internal/users/repo"

	"github.com/ghulammuzz/backend-parkerin/pkg/log"
)

type ApplicationService interface {
	CreateApply(userID, storeID int, isDirectHire bool) error
	ReviewApplications(storeID int) ([]appEntity.ApplicationResponse, error)
	ReviewApplicationsUser(userID int, isDirectHire bool) ([]appEntity.ApplicationUserResponse, error)
	AcceptApplicationUser(appID, userID int) error
	RejectApplicationUser(appID, userID int) error
	AcceptApplicationStore(appID, storeID int) error
	RejectApplicationStore(appID, storeID int) error
	DeleteAppsInUser(userID, appID int) error
}

type applicationService struct {
	appRepo   appRepo.ApplicationRepository
	storeRepo storeRepo.StoreRepository
	userRepo  userRepo.UserRepository
}

func (s *applicationService) DeleteAppsInUser(userID, appID int) error {
	return s.appRepo.DeleteApplicantsByUserIDAppsID(userID, appID)
}

func (s *applicationService) ReviewApplicationsUser(userID int, isDirectHire bool) ([]appEntity.ApplicationUserResponse, error) {
	return s.appRepo.GetApplicationsByUser(userID, isDirectHire)
}

func (s *applicationService) CreateApply(userID, storeID int, isDirectHire bool) error {
	storeExists, err := s.storeRepo.IsStoreIDValid(storeID)
	if err != nil {
		return err
	}

	if !storeExists {
		return errors.New("invalid store ID")
	}

	userExists, err := s.userRepo.IsUserIDValid(userID)
	if err != nil {
		return err
	}

	if !userExists {
		return errors.New("invalid user ID")
	}

	appExists, err := s.appRepo.CheckApplicantsAlreadyExist(userID, storeID)
	if err != nil {
		log.Error("Database error:", err)
		return err
	}
	if appExists {
		return errors.New("application already exists for this user and store")
	}

	err = s.appRepo.Apply(userID, storeID, isDirectHire)
	if err != nil {
		return err
	}
	// log.Debug(fmt.Sprint(isDirectHire))

	return nil
}

func (s *applicationService) ReviewApplications(storeID int) ([]appEntity.ApplicationResponse, error) {
	return s.appRepo.GetApplicationsByStore(storeID)
}

func (s *applicationService) AcceptApplicationUser(appID, userID int) error {
	return s.appRepo.UpdateApplicationStatusUser(appID, userID, "accepted")
}

func (s *applicationService) RejectApplicationUser(appID, userID int) error {
	return s.appRepo.UpdateApplicationStatusUser(appID, userID, "rejected")
}

func (s *applicationService) AcceptApplicationStore(appID, storeID int) error {
	return s.appRepo.UpdateApplicationStatusStore(appID, storeID, "accepted")
}

func (s *applicationService) RejectApplicationStore(appID, storeID int) error {
	return s.appRepo.UpdateApplicationStatusStore(appID, storeID, "rejected")
}

func NewApplicationService(appRepo appRepo.ApplicationRepository, storeRepo storeRepo.StoreRepository, userRepo userRepo.UserRepository) ApplicationService {
	return &applicationService{
		appRepo:   appRepo,
		storeRepo: storeRepo,
		userRepo:  userRepo,
	}
}
