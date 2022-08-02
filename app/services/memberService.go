package services

import (
	"github.com/dgrijalva/jwt-go"
	"go-membership/app/models"
	"go-membership/app/repositories/postgres"
	"go-membership/configs"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var (
	secret          = configs.GetJwtSecret()
	userRepository  = postgres.NewUserRepository()
	tokenRepository = postgres.NewTokenRepository()
)

type RegisterStatus int
type LoginStatus int

const (
	StatusRegisterSuccess                 RegisterStatus = 0
	StatusRegisterEmailRegistered         RegisterStatus = 1
	StatusRegisterEmailRegisterFailed     RegisterStatus = 2
	StatusRegisterEncryptedPasswordFailed RegisterStatus = 3
	StatusRegisterInternalError           RegisterStatus = 4

	StatusLoginSuccess           LoginStatus = 0
	StatusLoginEmailNotExist     LoginStatus = 1
	StatusLoginPasswordInvalid   LoginStatus = 2
	StatusLoginInternalError     LoginStatus = 3
	StatusLoginSignJwtFailed     LoginStatus = 4
	StatusLoginInsertTokenFailed LoginStatus = 5
)

type RegisterError struct {
	message string
}

type MemberService struct {
}

func NewMemberService() *MemberService {
	return &MemberService{}
}

func (m *MemberService) Register(name, email, password string) RegisterStatus {
	// is user registered?
	count, err := userRepository.CountByEmail(email)
	if err != nil {
		log.Println(err)
		return StatusRegisterInternalError
	}

	if count > 0 {
		return StatusRegisterEmailRegistered
	}

	encryptedPassword, err := m.cryptPassword(password)
	if err != nil {
		log.Println("CryptPassword Failed:", err)
		return StatusRegisterEmailRegisterFailed
	}

	// register
	users := models.NewUser(name, email, encryptedPassword)
	err = userRepository.Insert(users)
	if err != nil {
		log.Println(err)
		return StatusRegisterInternalError
	}

	count, err = userRepository.CountByEmail(email)

	if err != nil {
		log.Println(err)
		return StatusRegisterInternalError
	}

	if count != 1 {
		return StatusRegisterEncryptedPasswordFailed
	}
	return StatusRegisterSuccess
}

func (m *MemberService) cryptPassword(password string) (string, error) {
	pwd := []byte(password)
	encryptedPwd, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	password = string(encryptedPwd)
	return password, err
}

func (m *MemberService) isPasswordCorrect(inputPassword, userPassword string) bool {
	byteHashed := []byte(userPassword)
	byteInputPassword := []byte(inputPassword)
	err := bcrypt.CompareHashAndPassword(byteHashed, byteInputPassword)
	if err != nil {
		return false
	}
	return true
}

func (m *MemberService) Login(email, password string) (string, LoginStatus) {

	user, err := userRepository.SelectByEmail(email)

	if err != nil {
		return "", StatusLoginInternalError
	}

	if user.ID == 0 {
		return "", StatusLoginEmailNotExist
	}

	// valid password
	if !m.isPasswordCorrect(password, user.Password) {
		return "", StatusLoginPasswordInvalid
	}

	// search present jwt token
	oneDayAfterTime := time.Now().AddDate(0, 0, 1)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   user.Email,
		ExpiresAt: oneDayAfterTime.Unix(),
	})

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println("Signed JwtToken Failed:", err)
		return "", StatusLoginSignJwtFailed
	}

	err = tokenRepository.Insert(models.NewToken(*user, tokenStr, oneDayAfterTime))
	if err != nil {
		log.Println(err)
		return "", StatusLoginInternalError
	}

	tokenModel, err := tokenRepository.SelectUnExpiredByUserId(int(user.ID))
	if err != nil {
		log.Println(err)
		return "", StatusLoginInternalError
	}

	if tokenModel.ID == 0 {
		log.Println("Insert token failed.")
		return "", StatusLoginInsertTokenFailed
	}

	return tokenStr, StatusLoginSuccess
}
