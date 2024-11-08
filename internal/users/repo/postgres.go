package repo

import (
	"database/sql"
	"errors"
	"time"

	userEntity "github.com/ghulammuzz/backend-parkerin/internal/users/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	List(page, limit int) (*userEntity.UserListResponse, error)
	Create(user *userEntity.UserRegisterRequest) error
	Detail(userID int) (*userEntity.UserDetailResponse, error)
	LoginUser(user *userEntity.UserLoginRequest) (*userEntity.UserJWT, error)
	LoginStore(user *userEntity.UserLoginRequest) (*userEntity.StoreJWT, error)
	IsPhoneNumberExists(phoneNumber string) (bool, error)
	IsUserIDValid(userID int) (bool, error)
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) IsUserIDValid(userID int) (bool, error) {
	query := "SELECT 1 FROM users WHERE id = $1 LIMIT 1"

	row := r.db.QueryRow(query, userID)
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

func (r *userRepository) List(page int, limit int) (*userEntity.UserListResponse, error) {
	offset := (page - 1) * limit
	query := `
		SELECT id, phone_number, name
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []userEntity.UserListSubResponse{}
	for rows.Next() {
		user := userEntity.UserListSubResponse{}
		if err := rows.Scan(
			&user.ID,
			&user.PhoneNumber,
			&user.Name,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &userEntity.UserListResponse{
		Users: users,
		Page:  page,
		Limit: limit,
	}, nil
}

func (r *userRepository) LoginStore(user *userEntity.UserLoginRequest) (*userEntity.StoreJWT, error) {
	query := `SELECT id, phone_number, name, password, role FROM users WHERE phone_number = $1`
	storeUser := &userEntity.StoreJWT{}

	err := r.db.QueryRow(query, user.PhoneNumber).Scan(&storeUser.ID, &storeUser.PhoneNumber, &storeUser.Name, &storeUser.Password, &storeUser.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storeUser.Password), []byte(user.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	if storeUser.Role == "store" {
		storeQuery := `SELECT store_name, address, latitude, longitude FROM stores WHERE user_id = $1`
		err = r.db.QueryRow(storeQuery, storeUser.ID).Scan(&storeUser.StoreName, &storeUser.Address, &storeUser.Latitude, &storeUser.Longitude)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("store data not found")
			}
			return nil, err
		}
	} else {
		return nil, errors.New("user is not a store")
	}

	return storeUser, nil
}

func (r *userRepository) LoginUser(user *userEntity.UserLoginRequest) (*userEntity.UserJWT, error) {
	query := `SELECT id, phone_number, name, password, role FROM users WHERE phone_number = $1`
	dbUser := &userEntity.UserJWT{}

	err := r.db.QueryRow(query, user.PhoneNumber).Scan(&dbUser.ID, &dbUser.PhoneNumber, &dbUser.Name, &dbUser.Password, &dbUser.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return dbUser, nil
}
func (r *userRepository) Create(user *userEntity.UserRegisterRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `INSERT INTO users (phone_number, name, password, role, is_verified, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	user.CreatedAt = time.Now().UnixMilli()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = tx.QueryRow(query, user.PhoneNumber, user.Name, string(hashedPassword), user.Role, user.IsVerified, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return err
	}

	if user.Role == "store" {
		storeQuery := `INSERT INTO stores (user_id, store_name, address, latitude, longitude, created_at, working_hours) VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := tx.Exec(storeQuery, user.ID, user.StoreName, user.Address, user.Latitude, user.Longitude, user.CreatedAt, user.WorkingHours)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepository) Detail(userID int) (*userEntity.UserDetailResponse, error) {
	user := &userEntity.UserDetailResponse{}
	query := `SELECT id, name, phone_number, role FROM users WHERE id = $1`
	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) IsPhoneNumberExists(phoneNumber string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1)`
	err := r.db.QueryRow(query, phoneNumber).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}
