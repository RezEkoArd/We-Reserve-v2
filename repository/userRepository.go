package repository

import (
	"errors"
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


func (r *UserRepository) GetAllUser() ([]models.User, error) {
	var users []models.User
	err := r.DB.Raw("SELECT id, name, email, role, created_at, updated_at from users").Scan(&users).Error
	if err != nil {
		return nil, err // Kembalikan nil dan error jika terjadi kesalahan
	}
	return users, nil 
}

func (r *UserRepository) GetUserByid(id int) (models.User, error) {
	var user models.User
	err := r.DB.Raw("SELECT id, name, email, role, created_at, updated_at FROM users WHERE id = $1", id).Scan(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
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


// Update USer
func (r *UserRepository) UpdateUser (id int, user *models.User) error {
	var currentUser models.User
	err := r.DB.Raw("SELECT email FROM users WHERE id = ?", id).Scan(&currentUser).Error
	if err != nil {
		return err
	}

	//
	if user.Email != currentUser.Email {
		emailExists, err := r.IsEmailExists(user.Email)
		if err != nil {
			return err
		}

		if emailExists {
			return errors.New("email sudah terdaftar")
		}
	}

	sqlQuery := `UPDATE users SET name = $1, email = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	err = r.DB.Exec(sqlQuery, user.Name, user.Email, id).Error
	if err != nil {
		return err 
	}
	return nil
}


// Delete User 
func (r *UserRepository) DeleteUser (id int, user *models.User) error {
	sqlQuery := `DELETE FROM users WHERE id = $1`
	err := r.DB.Exec(sqlQuery, id).Error
	if err != nil {
		return err
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