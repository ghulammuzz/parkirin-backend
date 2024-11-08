// repo/store_repository.go
package repo

import (
	"database/sql"
	"fmt"

	storeEntity "github.com/ghulammuzz/backend-parkerin/internal/store/entity"
)

type StoreRepository interface {
	List(page, limit int, isHiring bool) (storeEntity.ListStoreResponse, error)
	Detail(id int) (*storeEntity.DetailStoreResponse, error)
	DetailByUserID(id int) (*storeEntity.DetailStoreResponse, error)
	GetStoreIDByUserID(userID int) (int, error)
	UpdateIsHiring(isHiring bool, storeID int) error
	IsStoreIDValid(storeID int) (bool, error)
	UploadStoreIMG(storeID int, img string) error
	VerifiedStore(storeID int) error
}

type storeRepository struct {
	db *sql.DB
}

func (r *storeRepository) VerifiedStore(storeID int) error {
	query := `
		UPDATE users
		SET is_verified = true
		FROM stores
		WHERE stores.id = $1
		AND stores.user_id = users.id
	`
	_, err := r.db.Exec(query, storeID)
	if err != nil {
		return fmt.Errorf("failed to update user verification status: %w", err)
	}
	return nil
}


func (r *storeRepository) UploadStoreIMG(storeID int, img string) error {
	query := `UPDATE stores SET url_image = $1 WHERE id = $2`
	_, err := r.db.Exec(query, img, storeID)
	if err != nil {
		return fmt.Errorf("failed to update store image: %w", err)
	}
	return nil
}

func (r *storeRepository) IsStoreIDValid(storeID int) (bool, error) {
	query := "SELECT 1 FROM stores WHERE id = $1 LIMIT 1"

	row := r.db.QueryRow(query, storeID)
	var exists int
	err := row.Scan(&exists)
	// log.Debug("deb 1")

	if err != nil {
		if err == sql.ErrNoRows {
			// log.Debug("deb 2")
			return false, nil
		}
		// log.Debug("deb 3")
		return false, err
	}
	// log.Debug("deb 4")
	return true, nil
}

func (r *storeRepository) UpdateIsHiring(isHiring bool, storeID int) error {
	query := "UPDATE stores SET is_hiring = $1 WHERE id = $2"
	_, err := r.db.Exec(query, isHiring, storeID)
	return err
}

func (s *storeRepository) GetStoreIDByUserID(userID int) (int, error) {
	var intUserID int
	query := `SELECT id from stores where user_id = $1`
	err := s.db.QueryRow(query, userID).Scan(&intUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user with ID %d not found", userID)
		}
		return 0, err
	}

	return intUserID, nil
}

func (s *storeRepository) DetailByUserID(id int) (*storeEntity.DetailStoreResponse, error) {
	query := `
		SELECT s.id, s.user_id, s.store_name, s.address, s.latitude, s.longitude, 
		       s.working_hours, s.is_hiring, s.is_paid, u.is_verified, s.created_at
		FROM stores s
		JOIN users u ON s.user_id = u.id
		WHERE s.user_id = $1
	`

	storeDetail := &storeEntity.DetailStoreResponse{}
	err := s.db.QueryRow(query, id).Scan(
		&storeDetail.ID,
		&storeDetail.UserID,
		&storeDetail.StoreName,
		&storeDetail.Address,
		&storeDetail.Latitude,
		&storeDetail.Longitude,
		&storeDetail.WorkingHours,
		&storeDetail.IsHiring,
		&storeDetail.IsPaid,
		&storeDetail.IsVerified,
		&storeDetail.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, err
	}

	return storeDetail, nil
}

func (s *storeRepository) Detail(id int) (*storeEntity.DetailStoreResponse, error) {
	query := `
		SELECT id, user_id, store_name, address, latitude, longitude, working_hours, is_hiring, is_paid, created_at
		FROM stores
		WHERE id = $1
	`

	storeDetail := &storeEntity.DetailStoreResponse{}
	err := s.db.QueryRow(query, id).Scan(
		&storeDetail.ID,
		&storeDetail.UserID,
		&storeDetail.StoreName,
		&storeDetail.Address,
		&storeDetail.Latitude,
		&storeDetail.Longitude,
		&storeDetail.WorkingHours,
		&storeDetail.IsHiring,
		&storeDetail.IsPaid,
		&storeDetail.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("store with ID %d not found", id)
		}
		return nil, err
	}

	return storeDetail, nil
}

func (s *storeRepository) List(page, limit int, isHiring bool) (storeEntity.ListStoreResponse, error) {
	offset := (page - 1) * limit
	query := `
		SELECT id, user_id, store_name, address, working_hours, is_hiring, is_paid
		FROM stores WHERE is_hiring = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(query, isHiring, limit, offset)
	if err != nil {
		return storeEntity.ListStoreResponse{}, err
	}
	defer rows.Close()

	stores := []storeEntity.ListStoreSubResponse{}
	for rows.Next() {
		store := storeEntity.ListStoreSubResponse{}
		if err := rows.Scan(
			&store.ID,
			&store.UserID,
			&store.StoreName,
			&store.Address,
			&store.WorkingHours,
			&store.IsHiring,
			&store.IsPaid,
		); err != nil {
			return storeEntity.ListStoreResponse{}, err
		}
		stores = append(stores, store)
	}

	if err := rows.Err(); err != nil {
		return storeEntity.ListStoreResponse{}, err
	}

	return storeEntity.ListStoreResponse{
		Stores: stores,
		Page:   page,
		Limit:  limit,
	}, nil
}

func NewStoreRepository(db *sql.DB) StoreRepository {
	return &storeRepository{db: db}
}
