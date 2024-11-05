package repo

import (
	"database/sql"
	"errors"
	"time"

	appEntity "github.com/ghulammuzz/backend-parkerin/internal/applicants/entity"
)

type ApplicationRepository interface {
	Apply(userID, storeID int, isDirectHire bool) error
	GetApplicationsByStore(storeID int) ([]appEntity.ApplicationResponse, error)
	GetApplicationsByUser(userID int, isDirectHire bool) ([]appEntity.ApplicationUserResponse, error)
	UpdateApplicationStatusUser(appID, userID int, status string) error
	UpdateApplicationStatusStore(appID, storeID int, status string) error
	CheckApplicantsAlreadyExist(userID, storeID int) (bool, error)
	RejectedAllApplicantsByStoreID(storeID int) error
}

type applicationRepository struct {
	db *sql.DB
}

func (r *applicationRepository) RejectedAllApplicantsByStoreID(storeID int) error {
	query := `
		UPDATE applications
		SET status = 'rejected'
		WHERE store_id = $1 AND status != 'rejected'
	`

	result, err := r.db.Exec(query, storeID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no applications to reject or already rejected")
	}

	return nil
}

func (r *applicationRepository) CheckApplicantsAlreadyExist(userID, storeID int) (bool, error) {
	query := "SELECT 1 FROM applications WHERE tukang_id = $1 AND store_id = $2 LIMIT 1"

	var exists int
	err := r.db.QueryRow(query, userID, storeID).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// mult (apply)
func (r *applicationRepository) Apply(userID, storeID int, isDirectHire bool) error {
	query := `
		INSERT INTO applications (tukang_id, store_id, status, applied_at, updated_at, is_direct_hire)
		SELECT $1, $2, 'sent', $3, $4, $5
		FROM stores
		WHERE id = $2 AND is_hiring = true
		RETURNING id
	`

	var applicationID int
	err := r.db.QueryRow(query, userID, storeID, time.Now().Unix(), time.Now().Unix(), isDirectHire).Scan(&applicationID)

	if err != nil {
		return errors.New("cannot apply: store is not hiring or application exists")
	}
	return nil
}

// store (list app by store)
func (r *applicationRepository) GetApplicationsByStore(storeID int) ([]appEntity.ApplicationResponse, error) {
	query := `
		SELECT a.id, u.name, a.status
		FROM applications a
		JOIN users u ON a.tukang_id = u.id
		WHERE a.store_id = $1
	`
	rows, err := r.db.Query(query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []appEntity.ApplicationResponse
	for rows.Next() {
		var app appEntity.ApplicationResponse
		if err := rows.Scan(&app.ID, &app.UserName, &app.Status); err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}
	return applications, nil
}

func (r *applicationRepository) GetApplicationsByUser(userID int, isDirectHire bool) ([]appEntity.ApplicationUserResponse, error) {
	query := `
		SELECT a.id, s.store_name, s.address, a.status
		FROM applications a
		JOIN stores s ON a.store_id = s.id
		WHERE a.tukang_id = $1 and a.is_direct_hire = $2
	`
	rows, err := r.db.Query(query, userID, isDirectHire)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []appEntity.ApplicationUserResponse
	for rows.Next() {
		var app appEntity.ApplicationUserResponse
		if err := rows.Scan(&app.ID, &app.StoreName, &app.Address, &app.Status); err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return applications, nil
}

// mult
func (r *applicationRepository) UpdateApplicationStatusUser(appID, userID int, status string) error {
	query := `UPDATE applications SET status = $1 WHERE id = $2 AND tukang_id = $3`
	_, err := r.db.Exec(query, status, appID, userID)
	return err
}

func (r *applicationRepository) UpdateApplicationStatusStore(appID, storeID int, status string) error {
	query := `UPDATE applications SET status = $1 WHERE id = $2 AND store_id = $3`
	_, err := r.db.Exec(query, status, appID, storeID)
	return err
}

func NewApplicationRepository(db *sql.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}
