package usecase

/* クエストに関連するビジネスロジックを実装する部分 */

import (
	"bulletin-board-rest-api/model"
	"bulletin-board-rest-api/repository"
	"bulletin-board-rest-api/validator"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	Login(user model.User) (string, error) //JWTを返すためにstring型
	GetUserName(userId uint) (string, error)
	GetUserInfo(userId uint) (model.UserResponse, error)
	UpdateUserName(userId uint, userName string) error
}

type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	return &userUsecase{ur, uv}
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	if err := uu.uv.ValidateUserSignUp(user); err != nil {
		return model.UserResponse{}, err
	}
	//パスワードのハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10) //GenerateFromPassword関数により、パスワードをハッシュ化
	if err != nil {
		return model.UserResponse{}, err
	}
	newUser := model.User{Email: user.Email, Password: string(hash), UserName: user.UserName} //ハッシュ化したパスワードをnewUserに格納
	if err := uu.ur.CreateUser(&newUser); err != nil {                                        //引数のnewUserをDBに保存
		return model.UserResponse{}, err
	}
	resUser := model.UserResponse{ //レスポンス用のUserResponse型の変数を作成
		ID:       newUser.ID,
		Email:    newUser.Email,
		UserName: newUser.UserName,
	}
	return resUser, nil
}

func (uu *userUsecase) Login(user model.User) (string, error) {
	if err := uu.uv.ValidateUserLogIn(user); err != nil {
		return "", err //Loginメソッドの戻り値に適するように空の文字列とエラーを返す
	}
	//クライアントからのEmailがDB内に存在するかを確認
	storedUser := model.User{} //DBから取得したユーザー情報を格納するための変数
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	// ハッシュ化されたパスと元のパスの一致を比較
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	// JWTトークンの作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ //JWTの生成
		"user_id": storedUser.ID,                             //ユーザーIDの設定
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), //TODO: 有効期限の設定
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET"))) //著名済みの文字列を生成 <- トークンとして使うことで、信頼性↑
	if err != nil {
		return "", err
	}
	return tokenString, nil //JWTを返す
}

func (uu *userUsecase) GetUserName(userId uint) (string, error) {
	User := model.User{}
	if err := uu.ur.GetUserByID(&User, userId); err != nil {
		return "", err
	}
	return User.UserName, nil
}

func (uu *userUsecase) GetUserInfo(userId uint) (model.UserResponse, error) {
	User := model.User{}
	if err := uu.ur.GetUserByID(&User, userId); err != nil {
		return model.UserResponse{}, err
	}

	resUser := model.UserResponse{
		UserName: User.UserName,
		Email:    User.Email,
	}
	return resUser, nil
}

func (uu *userUsecase) UpdateUserName(userId uint, userName string) error {
	if err := uu.ur.UpdateUserName(userId, userName); err != nil {
		return err
	}
	return nil
}
