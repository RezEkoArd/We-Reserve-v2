package services

import (
	"errors"
	"strconv"
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

	// parse id int to string
	strID := strconv.Itoa(user.ID)

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.Email, user.Role, strID)
	if err != nil {
		return "", err
	}

	return token, nil
}


// CRUD USER MANAGEMENT

func (s *UserService) DeleteUser(id int) error {
	// Cek User Exist
	_, err := s.userRepo.GetUserByid(id)
	if err != nil {
		return errors.New("user not Found")
	}

	// Delete
	err = s.userRepo.DeleteUser(id)
	if err != nil{
		return err
	}

	return nil
}

// Get All user 
func (s *UserService) GetAllUser() ([]models.User, error) {
	users, err := s.userRepo.GetAllUser()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Get UserById
func (s *UserService) GetUserById(id int) (*models.User, error) {

	user, err := s.userRepo.GetUserByid(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update User

func (s *UserService) UpdateUser(id int, user models.User)  error {
	// Validasi min satu field yang disi
	if user.Name == "" && user.Email == "" && user.Password == "" {
		return errors.New("minimal satu field harus diisi")
	}

	//Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password),10)
		if err != nil {
			return err
		}
	user.Password = string(hashedPassword)
	
	
	err = s.userRepo.UpdateUser(id, &models.User{
		Name:      user.Name,
		Email:     user.Email,
		Password:  string(hashedPassword),
	})

	if err != nil {
		return err
	}

	return nil
}