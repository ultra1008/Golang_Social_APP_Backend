package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/niklod/highload-social-network/config"
	"github.com/niklod/highload-social-network/user/city"
	"github.com/niklod/highload-social-network/user/interest"

	"github.com/gin-gonic/gin"
)

const (
	userSessionKey = "user"
)

type ViewData struct {
	Citys             []city.City
	Errors            []interface{}
	Messages          []interface{}
	Interests         []interest.Interest
	User              *User
	AuthenticatedUser *User
	UsersAreFriends   bool
}

type UserHandler struct {
	userService     *Service
	cityService     *city.Service
	interestService *interest.Service
	sessionStore    *sessions.CookieStore
}

func NewHandler(
	userService *Service,
	cityService *city.Service,
	sessionStore *sessions.CookieStore,
	interestService *interest.Service,
) *UserHandler {
	return &UserHandler{
		userService:     userService,
		cityService:     cityService,
		sessionStore:    sessionStore,
		interestService: interestService,
	}
}

func (u *UserHandler) HandleUserRegistrate(c *gin.Context) {
	user := getUser(c)
	if user != nil {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/user/%s", user.Login))
	}

	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		log.Printf("get session user handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	messages := session.Flashes()

	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("saving session: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	interests, err := u.interestService.Interests()
	if err != nil {
		log.Printf("gettings interests on registration page: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "registrate", ViewData{Interests: interests, Errors: messages})
}

func (u *UserHandler) HandleUserRegistrateSubmit(c *gin.Context) {
	var errors []interface{}

	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		log.Printf("get session user handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	req := &UserCreateRequest{}
	if err := c.ShouldBind(&req); err != nil {
		errors = append(errors, err.Error())
		c.HTML(http.StatusOK, "registrate", ViewData{Errors: errors})
		return
	}

	if err := c.Request.ParseForm(); err != nil {
		log.Printf("parse form user handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	req.Interests = c.Request.Form["inputInterests"]

	if err := req.Validate(); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, fieldError{err: e}.String())
		}
	}

	if len(errors) > 0 {
		for _, e := range errors {
			session.AddFlash(e)
		}

		if err := session.Save(c.Request, c.Writer); err != nil {
			log.Printf("saving session: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/registrate")
		return
	}

	_, err = u.userService.Create(req.ConverIntoUser())
	if err != nil {
		fmt.Printf("user creation: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	session.AddFlash("Регистрация успешно пройдена")

	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("saving session: %v", err)
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

	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		log.Printf("user detail, getting session: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	messages := session.Flashes()

	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("user detail, saving session: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if user != nil {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/user/%s", user.Login))
	}
	c.HTML(http.StatusOK, "login", ViewData{Messages: messages})
}

func (u *UserHandler) HandleUserLoginSubmit(c *gin.Context) {
	var req UserLoginRequest
	var errors []interface{}

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

	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		log.Printf("user detail, getting session: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

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

	userFriends, err := u.userService.Friends(user.ID)
	if err != nil {
		log.Printf("user detail, getting friends: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	user.Friends = userFriends

	userInterests, err := u.interestService.InterestsByUserId(user.ID)
	if err != nil {
		log.Printf("user detail, getting interests: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	user.Interests = userInterests

	user.Sanitize()

	if authUser != nil {
		authUserFriends, err := u.userService.Friends(authUser.ID)
		if err != nil {
			log.Printf("user detail, getting friends: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		authUser.Friends = authUserFriends
	}

	data := ViewData{
		Messages:          session.Flashes(),
		User:              user,
		AuthenticatedUser: authUser,
		UsersAreFriends:   u.userService.IsUsersAreFriends(authUser, user),
	}

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Printf("save session with flashes: %v", err)
	}

	c.HTML(http.StatusOK, "user_detail", data)
}

func (u *UserHandler) HandleAddFriend(c *gin.Context) {
	authUser := getUser(c)
	userLogin := c.Param("login")

	if authUser == nil {
		c.HTML(http.StatusUnauthorized, "login", nil)
		return
	}

	user, err := u.userService.GetUserByLogin(userLogin)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	err = u.userService.AddFriend(authUser.ID, user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Пользователь %s %s успешно добавлен в друзья", user.FirstName, user.Lastname)

	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		log.Printf("get session user handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	session.AddFlash(msg)
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Printf("save session with flashes: %v", err)
	}

	redirectLocation := fmt.Sprintf("/user/%s", user.Login)

	c.Redirect(http.StatusMovedPermanently, redirectLocation)
}

func (u *UserHandler) HandleDeleteFriend(c *gin.Context) {
	authUser := getUser(c)
	userLogin := c.Param("login")

	if authUser == nil {
		c.HTML(http.StatusUnauthorized, "login", nil)
		return
	}

	user, err := u.userService.GetUserByLogin(userLogin)
	if err != nil {
		log.Printf("get user by login in handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	err = u.userService.DeleteFriend(authUser.ID, user.ID)
	if err != nil {
		log.Printf("delete friend in handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Пользователь %s %s успешно удален из друзей", user.FirstName, user.Lastname)

	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		log.Printf("get session user handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	session.AddFlash(msg)

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Printf("save session with flashes: %v", err)
	}

	redirectLocation := fmt.Sprintf("/user/%s", user.Login)

	c.Redirect(http.StatusMovedPermanently, redirectLocation)
}

func (u *UserHandler) HandleUsersList(c *gin.Context) {
	authUser := getUser(c)

	users, err := u.userService.Users()
	if err != nil {
		log.Printf("gettings user list in handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "user_list", struct {
		Users             []User
		AuthenticatedUser *User
	}{users, authUser})
}

func (u *UserHandler) AuthMiddleware(c *gin.Context) {
	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, exist := session.Values[userSessionKey]
	if exist && session.Options.MaxAge > 0 {
		c.Set(userSessionKey, user)
	}
}

func getUser(c *gin.Context) *User {
	val, ok := c.Get(userSessionKey)
	if !ok {
		return nil
	}

	var user = User{}

	user, ok = val.(User)
	if !ok {
		return nil
	}
	return &user
}
