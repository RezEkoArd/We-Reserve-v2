package services

import (
	"errors"
	"wereserve/middleware"
	"wereserve/models"
	"wereserve/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// RegisterUser
func (s *UserService) RegisterUser( user *models.User ) ( error) {
	//Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password),10)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	//Set default role to 'customer if not provided
	if user.Role == "" {
		user.Role = "customer"
	}

	// Validate user role
	if !s.userRepo.IsValidUserRole(user.Role) {
		return errors.New("invalid user role")
	}

	//Check if email already use
	emailExist, err := s.userRepo.IsEmailExists(user.Email)
	if err != nil {
		return err
	}

	if emailExist {
		return errors.New("email sudah terdaftar")
	}

	// Create user
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) LoginUser(email, password string) (string, error) {
	// ngambil data 
	var user models.User
	err := s.userRepo.DB.Raw("SELECT id, name, email, password, role FROM users WHERE email = ?", email).Scan(&user).Error
	if err != nil {
		return "", nil
	}

	if user.ID == 0 {
		return "", errors.New("user tidak ditemukan")
	}

	//compare
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
