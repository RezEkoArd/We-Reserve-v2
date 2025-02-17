package repository

import (
	"errors"
	"fmt"
	"time"

	"wereserve/models"

	"gorm.io/gorm"
)


type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}


func (r *UserRepository) IsValidUserRole(role string) bool {
	validRoles := []string{"admin","customer"}
	for _, valid := range validRoles{
		if role == valid {
			return true
		}
	}
	return false
}

// Cek email apakah sudah ada didatabase
func (r *UserRepository) IsEmailExists(email string) (bool, error) {
	var exists bool
	err := r.DB.Raw("SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) IsUserExists(id int) (bool, error) {
	var existId	bool
	err := r.DB.Raw("SELECT EXISTS (SELECT 1 FROM users WHERE id = ?)", id).Scan(&existId).Error
	if err != nil {
		return false, err
	}
	return existId, nil
}


func (r *UserRepository) GetAllUser() ([]models.User, error) {
	var users []models.User
	err := r.DB.Raw("SELECT id, name, email, role, created_at, updated_at from users").Scan(&users).Error
	if err != nil {
		return nil, err // Kembalikan nil dan error jika terjadi kesalahan
	}
	return users, nil 
}

func (r *UserRepository) GetUserByid(id int) (*models.User, error) {
	var user models.User

	query := `SELECT id, name, email, role, created_at, updated_at FROM users WHERE id = $1`

	err := r.DB.Raw(query, id).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with Id %d not found", id)
		}
		return nil, err
	}
	
	return &models.User{
		ID:       id,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

// Create user
func (r *UserRepository) CreateUser(user *models.User) error {
	emailExists, err := r.IsEmailExists(user.Email)
	if err != nil {
		return err
	}

	if emailExists {
		return errors.New("email sudah terdaftar")
	}

	//Validasi role Pengguna
	if !r.IsValidUserRole(user.Role) {
		return errors.New("invalid user role")
	}

	// membuat user baru menggunakan GORM
	err = r.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}


// Update User
func (r *UserRepository) UpdateUser (id int, user *models.User) error {
	var currentUser models.User
	err := r.DB.First(&currentUser, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user tidak ditemukan")
		}
		return err
	}

	// Validasi email jika diubah
	if user.Email != "" && user.Email != currentUser.Email {
		emailExists, err := r.IsEmailExists(user.Email)
		if err != nil {
			return err
		}

		if emailExists {
			return errors.New("email sudah terdaftar")
		}
	}
	
	// Update hanya field yang diisi
	updates := make(map[string]interface{})
	if user.Name != "" {
		updates["name"] = user.Name
	}
	if user.Email != "" {
		updates["email"] = user.Email
	}
	if user.Password != "" {
		updates["password"] = user.Password
	}
	updates["updated_at"] = time.Now()

	//eksekusi query Update
	err = r.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return err
	}

	return nil 
}


// Delete User 
func (r *UserRepository) DeleteUser (id int) error {
	sqlQuery := `DELETE FROM users WHERE id = $1`
	result := r.DB.Exec(sqlQuery, id)

	if result.Error != nil {
		return result.Error
	}

	//validate if error 
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}


//
func (r *UserRepository) GetAllUserByRole (role string) ([]models.User, error) {
	var users []models.User

	//Eksekusi query
	err := r.DB.Raw("SELECT id, name, email, role FROM users WHERE role = ?", role).Scan(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}