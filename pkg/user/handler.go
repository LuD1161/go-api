package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/LuD1161/restructuring-tnbt/pkg/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

// Handler : Handler for User
type Handler interface {
	Login(c *gin.Context)
	GetUserByID(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userHandler struct {
	userService Service
	log         *logrus.Logger
}

// NewHandler : Returns handler for new user service
func NewHandler(userService Service, log *logrus.Logger) Handler {
	return &userHandler{
		userService,
		log,
	}
}

// @Summary Get User by ID
// @Description Get the user details by ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param   id     path    int     true        "User ID"
// @Param Authorization header string true "JWT header starting with the Bearer"
// @Success 200 {object} UserInfoPayload
// @Router /user/{id} [get]
func (h *userHandler) GetUserByID(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	h.log.Infof("in GetUserByID")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.GetUserByID").Error(),
		})
		return
	}
	user, err := h.userService.GetUserByID(uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.GetUserByID").Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user.UserInfoPayload)
}

// FIXME : Probably need to define a USER Response struct -> Done

// CreateUser godoc
// @Summary Create a user
// @Description Create a user
// @Tags User
// @Accept  json
// @Produce  json
// @Param json body CreateUserPayload true "Create User"
// @Success 200 {object} User
// @Router /user/ [post]
// CreateUser : Creates new user
func (h *userHandler) CreateUser(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.CreateUser").Error(),
		})
		return
	}
	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.CreateUser").Error(),
		})
		return
	}
	// FIXME : How to get the Prepare functionality here, currently Prepare sets u.ID = u.ID
	// instead of u.ID = 0
	// user.Prepare()
	v := validator.New()
	err = v.Struct(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.CreateUser").Error(),
		})
		return
	}
	userCreated, err := h.userService.CreateUser(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.CreateUser").Error(),
		})
		return
	}
	c.Header("Location", fmt.Sprintf("%s%s/%d", c.Request.Host, c.Request.RequestURI, userCreated.ID))
	c.JSON(http.StatusOK, userCreated.UserInfoPayload)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user
// @Tags User
// @Accept  json
// @Produce  json
// @Param json body UpdateUserPayload true "Can only update current user's password as of now"
// @Param Authorization header string true "JWT header starting with the Bearer"
// @Success 200 {object} UserInfoPayload
// @Router /user/ [put]
// UpdateUser : Updates new user
func (h *userHandler) UpdateUser(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.UpdateUser").Error(),
		})
		return
	}
	user := new(User)
	updateUser := new(UpdateUserPayload)
	// err = json.Unmarshal(body, &user)
	err = json.Unmarshal(body, &updateUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.UpdateUser").Error(),
		})
		return
	}
	tokenID, err := auth.ExtractTokenID(c.Request)
	// Change this error to unauthorized
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.UpdateUser").Error(),
		})
		return
	}
	// FIXME : Check whether it's needed or not
	user.Password = updateUser.Password
	h.userService.Prepare(user)
	v := validator.New()
	err = v.Struct(updateUser)
	if err != nil {
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.UpdateUser").Error(),
		})
		return
	}
	user.ID = tokenID
	updatedUser, err := h.userService.UpdateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.UpdateUser").Error(),
		})
		return
	}
	c.Header("Location", fmt.Sprintf("%s%s/%d", c.Request.Host, c.Request.RequestURI, updatedUser.ID))
	c.JSON(http.StatusOK, user.UserInfoPayload)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user
// @Tags User
// @Accept  json
// @Produce  json
// @Param json body User true "Only own user can be deleted"
// @Success 200 {object} User

// DeleteUser : Deletes new user
func (h *userHandler) DeleteUser(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.DeleteUser").Error(),
		})
		return
	}
	tokenID, err := auth.ExtractTokenID(c.Request)
	// Change this error to unauthorized
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.DeleteUser").Error(),
		})
		return
	}
	if tokenID != uid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.DeleteUser").Error(),
		})
		return
	}
	status, err := h.userService.DeleteUser(uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.DeleteUser").Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}

// Login godoc
// @Summary Login
// @Description Login to get a JWToken
// @Tags Login
// @Accept  json
// @Produce  json
// @Param json body LoginPayload true "Login to get the JWToken"
// @Success 200 {string} string "JWToken here"
// @Router /login [post]
// Login : Login to get a new JWT
func (h *userHandler) Login(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.Login").Error(),
		})
		return
	}
	user := new(User)
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.Login").Error(),
		})
		return
	}

	// user.Prepare()
	// err = user.Validate("login")
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnprocessableEntity, err)
	// 	return
	// }
	token, err := h.userService.Login(user.Username, user.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": errors.Wrap(err, "pkg.user.handler.Login").Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
