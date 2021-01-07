package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/niklod/highload-social-network/config"
	"github.com/niklod/highload-social-network/user/city"

	"github.com/gin-gonic/gin"
)

const (
	userSessionKey = "user"
)

type ViewData struct {
	Citys             []city.City
	Errors            []string
	Messages          []string
	User              *User
	AuthenticatedUser *User
}

type UserHandler struct {
	userService  *Service
	cityService  *city.Service
	sessionStore *sessions.CookieStore
}

func NewHandler(userService *Service, cityService *city.Service, ss *sessions.CookieStore) *UserHandler {
	return &UserHandler{
		userService:  userService,
		cityService:  cityService,
		sessionStore: ss,
	}
}

func (u *UserHandler) HandleUserRegistrate(c *gin.Context) {
	user := getUser(c)
	if user != nil {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/user/%s", user.Login))
	}
	c.HTML(http.StatusOK, "registrate", nil)
}

func (u *UserHandler) HandleUserRegistrateSubmit(c *gin.Context) {
	var errors []string

	req := &UserCreateRequest{}
	if err := c.ShouldBind(&req); err != nil {
		errors = append(errors, err.Error())
		c.HTML(http.StatusOK, "registrate", ViewData{Errors: errors})
		return
	}

	if err := req.Validate(); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, fieldError{err: e}.String())
		}
	}

	if len(errors) > 0 {
		c.HTML(http.StatusOK, "registrate", ViewData{Errors: errors})
		return
	}

	_, err := u.userService.Create(req.ConverIntoUser())
	if err != nil {
		fmt.Printf("user creation: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/login")
}

func (u *UserHandler) HandleUserLogout(c *gin.Context) {
	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		log.Printf("get session user handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	session.Values[userSessionKey] = User{}
	session.Options.MaxAge = -1

	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("saving session: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "login", nil)
}

func (u *UserHandler) HandleUserLogin(c *gin.Context) {
	user := getUser(c)
	if user != nil {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/user/%s", user.Login))
	}
	c.HTML(http.StatusOK, "login", nil)
}

func (u *UserHandler) HandleUserLoginSubmit(c *gin.Context) {
	var req UserLoginRequest
	var errors []string

	if err := c.ShouldBind(&req); err != nil {
		errors = append(errors, err.Error())
		c.HTML(http.StatusOK, "login", ViewData{Errors: errors})
		return
	}

	if err := req.Validate(); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, fieldError{err: e}.String())
		}
	}

	if len(errors) > 0 {
		c.HTML(http.StatusOK, "login", ViewData{Errors: errors})
		return
	}

	user, err := u.userService.GetUserByLogin(req.Login)
	if err != nil {
		log.Printf("get user by id handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if user == nil || !u.userService.CheckPasswordsEquality(req.Password, user.Password) {
		errors = append(errors, "Указан неверный логин или пароль")
		c.HTML(http.StatusOK, "login", ViewData{Errors: errors})
		return
	}

	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		log.Printf("get session user handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	session.Values[userSessionKey] = *user
	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("saving session: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/user/%s", user.Login))
}

func (u *UserHandler) HandleUserDetail(c *gin.Context) {
	authUser := getUser(c)
	userLogin := c.Param("login")

	user, err := u.userService.GetUserByLogin(userLogin)
	if err != nil {
		log.Printf("user detail, getting user: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if user == nil {
		c.Status(http.StatusNotFound)
		return
	}

	user.Sanitize()

	c.HTML(http.StatusOK, "user_detail", ViewData{User: user, AuthenticatedUser: authUser})
}

func (u *UserHandler) AuthMiddleware(c *gin.Context) {
	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, exist := session.Values[userSessionKey]
	if exist {
		c.Set(userSessionKey, user)
	}
}

func getUser(c *gin.Context) *User {
	val, _ := c.Get(userSessionKey)
	var user = User{}

	user, ok := val.(User)
	if !ok {
		return nil
	}
	return &user
}
