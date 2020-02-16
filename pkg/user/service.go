package user

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/LuD1161/restructuring-tnbt/pkg/middlewares/auth"
	hashing "github.com/LuD1161/restructuring-tnbt/pkg/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Service : UserService
type Service interface {
	Login(username, password string) (string, error) // returns JWToken
	BeforeSave(*User) error                          // TBD later : Not sure if this is needed
	Prepare(*User)                                   // TBD later : Not sure how this is needed, if it can be incorporated in UpdateUser
	CreateUser(*User) (*User, error)
	UpdateUser(*User) (*User, error)
	DeleteUser(uint64) (int64, error)
	GetUserByID(uint64) (*User, error)
}

type service struct {
	repo Repository
	log  *logrus.Logger
}

// NewService creates a listing service with the necessary dependencies
func NewService(repo Repository, log *logrus.Logger) Service {
	return &service{
		repo,
		log,
	}
}

// Prepare : Prepare the user-data to be updated; Invoked on update user and login
func (s *service) Prepare(u *User) {
	// FIXME : Check whethere it should be u.ID = u.ID or 0
	u.ID = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt // Diff : Since initially both should be same
}

// BeforeSave : Operations to be performed before saving the
// user to database will be done here
func (s *service) BeforeSave(u *User) error {
	hashedPassword, err := hashing.Hash(u.Password)
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}
	u.Password = string(hashedPassword)
	return nil

}

// CreateUser : Creates the user in database
func (s *service) CreateUser(u *User) (*User, error) {
	s.BeforeSave(u)
	return s.repo.CreateUser(u)
}

// UpdateUser : Update user details
func (s *service) UpdateUser(u *User) (*User, error) {
	s.BeforeSave(u)
	return s.repo.UpdateUser(u)
}

// GetUserByID : Finds a user by ID
func (s *service) GetUserByID(uid uint64) (*User, error) {
	s.log.Info("in pkg.user.service.GetUserByID")
	return s.repo.GetUserByID(uid)
}

// DeleteAUser : Deletes a user from the database
func (s *service) DeleteUser(uid uint64) (int64, error) {
	return s.repo.DeleteUser(uid)
}

// Login : Returns JWT for login verification
func (s *service) Login(username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)

	if err != nil {
		logrus.WithField("username", username).Error("Unable to fetch account")
		return "", err
	}

	if user == nil {
		return "", errors.New("Invalid Username")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.WithFields(logrus.Fields{"username": username, "error": err.Error()}).Error("Invalid login")
		return "", err
	}

	token, err := auth.CreateToken(user.ID)

	if err != nil {
		logrus.WithFields(logrus.Fields{"username": username, "error": err}).Error("Unable to generate token")
	}

	return token, nil
}
