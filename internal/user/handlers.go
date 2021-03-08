package user

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"

	"github.com/niklod/highload-social-network/config"
	"github.com/niklod/highload-social-network/internal/user/city"
	"github.com/niklod/highload-social-network/internal/user/interest"
	"github.com/niklod/highload-social-network/internal/user/post"
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
	Feed              post.Feed
}

type UserHandler struct {
	userService     *Service
	cityService     *city.Service
	interestService *interest.Service
	postService     *post.Service
	sessionStore    *sessions.CookieStore
}

func NewHandler(
	userService *Service,
	cityService *city.Service,
	postService *post.Service,
	sessionStore *sessions.CookieStore,
	interestService *interest.Service,
) *UserHandler {
	return &UserHandler{
		userService:     userService,
		cityService:     cityService,
		postService:     postService,
		sessionStore:    sessionStore,
		interestService: interestService,
	}
}

func (u *UserHandler) HandleUserRegistrate(c *gin.Context) {
	user := getUser(c)
	if user != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/user/%s", user.Login))
		return
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
	var handlerErrors []interface{}

	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		log.Printf("get session user handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	req := &UserCreateRequest{}
	if err := c.ShouldBind(&req); err != nil {
		handlerErrors = append(handlerErrors, err.Error())
		c.HTML(http.StatusOK, "registrate", ViewData{Errors: handlerErrors})
		return
	}

	if err := req.Validate(); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			handlerErrors = append(handlerErrors, fieldError{err: e}.String())
		}
	}

	if len(handlerErrors) > 0 {
		for _, e := range handlerErrors {
			session.AddFlash(e)
		}

		if err := session.Save(c.Request, c.Writer); err != nil {
			log.Printf("saving session: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Redirect(http.StatusSeeOther, "/registrate")
		return
	}

	_, err = u.userService.Create(req.ConverIntoUser())
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExist) {
			session.AddFlash("Пользователь с таким логином уже существует")

			if err := session.Save(c.Request, c.Writer); err != nil {
				log.Printf("saving session: %v", err)
				c.Status(http.StatusInternalServerError)
				return
			}
			c.Redirect(http.StatusFound, "/registrate")
			return
		}

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

	c.Redirect(http.StatusFound, "/login")
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

	c.Redirect(http.StatusFound, "login")
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
		c.Redirect(http.StatusFound, fmt.Sprintf("/user/%s", user.Login))
		return
	}
	c.HTML(http.StatusOK, "login", ViewData{Messages: messages})
}

func (u *UserHandler) HandleUserLoginSubmit(c *gin.Context) {
	var req UserLoginRequest
	var handlerErrors []interface{}

	if err := c.ShouldBind(&req); err != nil {
		handlerErrors = append(handlerErrors, err.Error())
		c.HTML(http.StatusBadRequest, "login", ViewData{Errors: handlerErrors})
		return
	}

	if err := req.Validate(); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			handlerErrors = append(handlerErrors, fieldError{err: e}.String())
		}
	}

	if len(handlerErrors) > 0 {
		c.HTML(http.StatusUnprocessableEntity, "login", ViewData{Errors: handlerErrors})
		return
	}

	user, err := u.userService.GetUserByLogin(req.Login)
	if err != nil {
		log.Printf("get user by id handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if user == nil || !u.userService.CheckPasswordsEquality(req.Password, user.Password) {
		handlerErrors = append(handlerErrors, "Указан неверный логин или пароль")
		c.HTML(http.StatusForbidden, "login", ViewData{Errors: handlerErrors})
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

	c.Redirect(http.StatusFound, fmt.Sprintf("/user/%s", user.Login))
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

	userPosts, err := u.postService.PostsByUserId(user.ID)
	if err != nil {
		log.Printf("user detail, getting posts: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	user.Posts = userPosts

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
		c.Redirect(http.StatusUnauthorized, "login")
		return
	}

	user, err := u.userService.GetUserByLogin(userLogin)
	if err != nil {
		log.Printf("find user by login: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	err = u.userService.AddFriend(authUser.ID, user.ID)
	if err != nil {
		log.Printf("addint to friends: %v", err)
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

	c.Redirect(http.StatusSeeOther, redirectLocation)
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

	c.Redirect(http.StatusSeeOther, redirectLocation)
}

func (u *UserHandler) HandleAddPost(c *gin.Context) {
	authUser := getUser(c)
	postBody := c.PostForm("post")

	if authUser == nil {
		c.Redirect(http.StatusUnauthorized, "login")
		return
	}

	post := &post.Post{
		Body: postBody,
		Author: post.Author{
			ID:        authUser.ID,
			FirstName: authUser.FirstName,
			LastName:  authUser.Lastname,
			Login:     authUser.Login,
		},
	}

	err := u.postService.Add(post)
	if err != nil {
		log.Printf("addint post: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		log.Printf("get session user handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Printf("save session with flashes: %v", err)
	}

	redirectLocation := fmt.Sprintf("/user/%s", authUser.Login)

	c.Redirect(http.StatusSeeOther, redirectLocation)
}

func (u *UserHandler) HandleFeed(c *gin.Context) {
	authUser := getUser(c)

	session, err := u.sessionStore.Get(c.Request, config.SessionName)
	if err != nil {
		log.Printf("user detail, getting session: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if authUser == nil {
		c.Redirect(http.StatusUnauthorized, "/login")
		return
	}

	feed, err := u.postService.UserFeed(authUser.ID)
	if err != nil {
		log.Printf("feed page, getting user feed: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	data := ViewData{
		Messages:          session.Flashes(),
		AuthenticatedUser: authUser,
		Feed:              feed,
	}

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Printf("save session with flashes: %v", err)
	}

	c.HTML(http.StatusOK, "user_feed", data)
}

func (u *UserHandler) HandleUsersList(c *gin.Context) {
	authUser := getUser(c)

	req := UserSearchRequest{}

	if err := c.ShouldBind(&req); err != nil {
		log.Printf("gettings user list in handler: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	users, err := u.userService.userRepo.GetByFirstAndLastName(req.FirstName, req.LastName)
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

	var user User

	user, ok = val.(User)
	if !ok {
		return nil
	}
	return &user
}
