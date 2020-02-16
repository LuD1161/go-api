package postgres

import (
	"errors"
	"time"

	"github.com/LuD1161/restructuring-tnbt/pkg/user"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type userRepository struct {
	db *gorm.DB
}

// NewPostgresUserRepository : To create new postgres repository connection
func NewPostgresUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) CreateUser(user *user.User) (*user.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(u *user.User) (*user.User, error) {
	var err error
	user := new(user.User)
	logrus.Info("Inside UpdateUser (in repo) : ", u.ID)
	db := r.db.Model(user).Where("id = ?", u.ID).Updates(
		map[string]interface{}{
			"password":   u.Password,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return user, db.Error
	}
	err = r.db.Where("id = ?", u.ID).First(&u).Error
	if err != nil {
		return user, err
	}
	return u, nil
}

func (r *userRepository) DeleteUser(uid uint64) (int64, error) {
	var err error
	user := new(user.User)
	err = r.db.Model(user).Where("id = ?", uid).Delete(&user).Error
	if err != nil {
		return 0, err
	}
	return r.db.RowsAffected, nil
}

func (r *userRepository) GetUserByID(uid uint64) (*user.User, error) {
	var err error
	user := new(user.User)
	err = r.db.Model(user).Where("id = ?", uid).First(&user).Error
	// Handle the specific case first
	if gorm.IsRecordNotFoundError(err) {
		return user, errors.New("User Not Found")
	}
	if err != nil {
		return user, err
	}
	return user, err
}

func (r *userRepository) GetUserByUsername(username string) (*user.User, error) {
	var err error // so that in the end when returning, err == nil ( so the signature doesn't change for any function)
	user := new(user.User)
	logrus.Info("username ", username)
	err = r.db.Model(user).Where("username = ?", username).First(&user).Error
	// Handle the specific case first
	if gorm.IsRecordNotFoundError(err) {
		return user, errors.New("User Not Found")
	}
	if err != nil {
		return user, err
	}
	return user, err
}
