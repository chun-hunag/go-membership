package membership

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-membership/app/models"
	"go-membership/app/repositories/postgres"
	"go-membership/app/utilis"
	"go-membership/configs"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var (
	secret          = configs.GetJwtSecret()
	userRepository  = postgres.NewUserRepository()
	tokenRepository = postgres.NewTokenRepository()
)

type registerForm struct {
	Name            string `form:"name"`
	Email           string `form:"email"`
	Password        string `form:"password"`
	PasswordConfirm string `form:"passwordConfirm"`
}

func (r *registerForm) IsPasswordSame() bool {
	return r.Password == r.PasswordConfirm
}

func (r *registerForm) Validate() (bool, []string) {
	invalids := make([]string, 0)
	if len(r.Name) < 4 {
		invalids = append(invalids, "Name is less than 4 characters.")
	}

	if !utilis.IsEmailValid(r.Email) {
		invalids = append(invalids, "Email is invalid.")
	}

	if len(r.Password) < 8 || len(r.PasswordConfirm) < 8 {
		invalids = append(invalids, "Password or PasswordConfirm is less than 8 words.")
	}

	if !r.IsPasswordSame() {
		invalids = append(invalids, "Different passwords entered twice.")
	}

	isValid := len(invalids) == 0
	return isValid, invalids
}

func (r *registerForm) cryptPassword() (bool, error) {
	pwd := []byte(r.Password)
	encryptedPwd, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)

	if err != nil {
		return false, err
	}

	r.Password = string(encryptedPwd)
	return true, err
}

func Registration(c *gin.Context) {
	var registerForm registerForm
	c.Bind(&registerForm)
	response := make(map[string]interface{})

	isValid, invalids := registerForm.Validate()
	if !isValid {
		response["message"] = "Input value is valid."
		response["Invalids"] = invalids
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// is user registered
	if userRepository.CountByEmail(registerForm.Email) > 0 {
		response["message"] = "The email had been registered."
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// crypt password
	if isSuccess, err := registerForm.cryptPassword(); !isSuccess {
		response["message"] = "Some thing wrong."
		log.Println("CryptPassword Failed:", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// register
	users := models.NewUser(registerForm.Name, registerForm.Email, registerForm.Password)
	userRepository.Insert(users)
	if userRepository.CountByEmail(registerForm.Email) != 1 {
		response["message"] = "The registration is failed."
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response["message"] = "Registered."
	c.JSON(http.StatusOK, response)
}

type loginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (l *loginForm) Validate() (bool, []string) {
	invalids := make([]string, 0)
	if !utilis.IsEmailValid(l.Email) {
		invalids = append(invalids, "Email is invalid.")
	}

	isValid := len(invalids) == 0
	return isValid, invalids
}

func (l *loginForm) isPasswordCorrect(password string) bool {
	byteHashed := []byte(password)
	byteInputPassword := []byte(l.Password)
	err := bcrypt.CompareHashAndPassword(byteHashed, byteInputPassword)
	if err != nil {
		return false
	}
	return true
}

func Login(c *gin.Context) {
	var loginForm loginForm
	c.Bind(&loginForm)
	response := make(map[string]interface{})

	isValid, invalids := loginForm.Validate()
	if !isValid {
		response["message"] = "Input value is valid."
		response["Invalids"] = invalids
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// get user from database
	user := userRepository.SelectByEmail(loginForm.Email)
	if user.ID == 0 {
		response["message"] = "User has this email is not existed."
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// valid password
	if !loginForm.isPasswordCorrect(user.Password) {
		response["message"] = "Password is not correct."
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// search present jwt token
	oneDayAfterTime := time.Now().AddDate(0, 0, 1)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   user.Email,
		ExpiresAt: oneDayAfterTime.Unix(),
	})

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		response["message"] = "Some thing wrong."
		log.Println("Signed JwtToken Failed:", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	tokenRepository.Insert(models.NewToken(*user, tokenStr, oneDayAfterTime))
	tokenModel := tokenRepository.SelectUnExpiredByUserId(int(user.ID))
	if tokenModel.ID == 0 {
		response["message"] = "Some thing wrong."
		log.Println("Insert token failed.")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response["message"] = "Login success."
	response["token"] = tokenStr
	c.JSON(http.StatusBadRequest, response)
}

func SendPasswordResetEmail(c *gin.Context) {
	// valid user exist or not

	// update remember token

	// send mailed

	fmt.Println("send user password reset email")
}

func ResetPassword(c *gin.Context) {
	// get user reme

	fmt.Println("user resentPassword")
}
