package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
}

type userController struct {
	userUsecase usecase.IUserUsecase
}

func NewUserController(u usecase.IUserUsecase) IUserController {
	return &userController{userUsecase: u}
}

func (u *userController) SignUp(c echo.Context) error {
	user := model.User{}
	// リクエストボディをバインド
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// ユーザーを作成
	userResponse, err := u.userUsecase.SignUp(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, userResponse)
}

func (u *userController) LogIn(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	token, err := u.userUsecase.LogIn(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/" // 全てのパスでアクセス可能
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = True
	cookie.HttpOnly = true                  // クライアントからのアクセスを禁止
	cookie.SameSite = http.SameSiteNoneMode // クロスサイトでもCookieを送信
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

func (u *userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}
