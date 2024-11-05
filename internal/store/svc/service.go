package svc

import (
	"errors"

	appRepo "github.com/ghulammuzz/backend-parkerin/internal/applicants/repo"
	"github.com/ghulammuzz/backend-parkerin/internal/store/entity"
	storeRepo "github.com/ghulammuzz/backend-parkerin/internal/store/repo"
	userRepo "github.com/ghulammuzz/backend-parkerin/internal/users/repo"
	"github.com/ghulammuzz/backend-parkerin/pkg/log"
)

type StoreService interface {
	ListStores(page, limit int, isHiring bool) (entity.ListStoreResponse, error)
	GetStoreDetail(id int) (*entity.DetailStoreResponse, error)
	DashboardStore(userId int) (*entity.DashboardStoreResponse, error)
	GetStoreIDByUserID(userID int) (int, error)
	UpdateIsHiring(isHiring bool, storeID int) error
	CheckStoreID(storeID int) (bool, error)
}

type storeService struct {
	storeRepo storeRepo.StoreRepository
	userRepo  userRepo.UserRepository
	appRepo   appRepo.ApplicationRepository
}

// CheckStoreID implements StoreService.
func (s *storeService) CheckStoreID(storeID int) (bool, error) {
	return s.storeRepo.IsStoreIDValid(storeID)
}

func (s *storeService) UpdateIsHiring(isHiring bool, storeID int) error {

	if !isHiring {
		err := s.appRepo.RejectedAllApplicantsByStoreID(storeID)
		if err != nil {
			log.Error("Failed to reject all applicants:", err)
			return errors.New("failed to reject all applicants")
		}
	}
	return s.storeRepo.UpdateIsHiring(isHiring, storeID)
}

func (s *storeService) GetStoreIDByUserID(userID int) (int, error) {
	return s.storeRepo.GetStoreIDByUserID(userID)
}

func (s *storeService) ListStores(page, limit int, isHiring bool) (entity.ListStoreResponse, error) {
	stores, err := s.storeRepo.List(page, limit, isHiring)
	if err != nil {
		return entity.ListStoreResponse{}, err
	}
	return stores, nil
}

func (s *storeService) GetStoreDetail(id int) (*entity.DetailStoreResponse, error) {
	return s.storeRepo.Detail(id)
}

func (s *storeService) DashboardStore(id int) (*entity.DashboardStoreResponse, error) {
	store, err := s.storeRepo.DetailByUserID(id)
	if err != nil {
		return nil, errors.New("error repo detail store by user id")
	}
	user, err := s.userRepo.Detail(id)
	if err != nil {
		return nil, errors.New("error repo detail by user id")
	}

	response := &entity.DashboardStoreResponse{
		ID:           store.ID,
		User:         *user,
		StoreName:    store.StoreName,
		Address:      store.Address,
		Latitude:     store.Latitude,
		Longitude:    store.Longitude,
		WorkingHours: store.WorkingHours,
		IsHiring:     store.IsHiring,
		IsPaid:       store.IsPaid,
		CreatedAt:    store.CreatedAt,
	}

	return response, nil
}

func NewStoreService(storeRepo storeRepo.StoreRepository, userRepo userRepo.UserRepository, appRepo appRepo.ApplicationRepository) StoreService {
	return &storeService{storeRepo: storeRepo, userRepo: userRepo, appRepo: appRepo}
}
