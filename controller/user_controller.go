package controller

/* リクエストの受け付けとレスポンスの生成 */

import (
	"bulletin-board-rest-api/model"
	"bulletin-board-rest-api/usecase"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	GetUserName(c echo.Context) error
	GetUserInfo(c echo.Context) error
	UpdateUserName(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

// usecaseを「Dependency Injection」するための関数（コンストラクタ）
// usecaseのインスタンスを受け取る
func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

func (uc *userController) SignUp(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil { //リクエストボディをuserにバインド（User型に変換して格納）
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userRes, err := uc.uu.SignUp(user) //usecaseのSignUpメソッドを呼び出し
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes) //Created(201)のステータスと、作成したユーザー情報を返す
}

func (uc *userController) LogIn(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	jwtToken, err := uc.uu.Login(user) //usecaseのLoginメソッドを呼び出し（JWTtokenが入る）
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// JSONとしてJWTトークンを返す
	return c.JSON(http.StatusOK, echo.Map{
		"token": jwtToken,
	})
}

func (uc *userController) LogOut(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (uc *userController) GetUserName(c echo.Context) error {
	// JWTのclaimsからユーザーIDを取得
	user := c.Get("user").(*jwt.Token) // jwtをデコードした内容を取得
	claims := user.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64)) // float64をuintにキャスト

	// ユーザーIDを元にユーザー名を取得
	username, err := uc.uu.GetUserName(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, username) //* ここでUserNameを取得してJSON形式で返す！
}

func (uc *userController) GetUserInfo(c echo.Context) error {
	// JWTのclaimsからユーザーIDを取得
	user := c.Get("user").(*jwt.Token) // jwtをデコードした内容を取得
	claims := user.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64)) // float64をuintにキャスト

	// ユーザーIDを元にユーザー名を取得
	userRes, err := uc.uu.GetUserInfo(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	response := map[string]interface{}{
		"email":     userRes.Email,
		"user_name": userRes.UserName,
	}

	return c.JSON(http.StatusOK, response) //* ここでUserNameを取得してJSON形式で返す！
}

func (uc *userController) UpdateUserName(c echo.Context) error {
	user := c.Get("user").(*jwt.Token) // jwtをデコードした内容を取得
	claims := user.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	req := model.UpdateUserNameRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := uc.uu.UpdateUserName(userId, req.UserName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
