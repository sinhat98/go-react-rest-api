package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user *model.User) (model.UserResponse, error)
	LogIn(user *model.User) (string, error)
}

type userUsecase struct {
	userRepository repository.IUserRepository
	userValidator  validator.IUserValidator
}

func NewUserUsecase(userRepository repository.IUserRepository, userValidator validator.IUserValidator) IUserUsecase {
	return &userUsecase{userRepository: userRepository, userValidator: userValidator}
}

func (u *userUsecase) SignUp(user *model.User) (model.UserResponse, error) {
	if err := u.userValidator.UserValidate(*user); err != nil {
		return model.UserResponse{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10) // パスワードをハッシュ化
	if err != nil {
		return model.UserResponse{}, err
	}
	newUser := model.User{
		Email:    user.Email,
		Password: string(hash),
	}
	if err := u.userRepository.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	return resUser, nil
}

func (u *userUsecase) LogIn(user *model.User) (string, error) {
	if err := u.userValidator.UserValidate(*user); err != nil {
		return "", err
	}
	storedUser := model.User{}
	if err := u.userRepository.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(), // 有効期限を12時間に設定
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET"))) // トークンを生成
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
