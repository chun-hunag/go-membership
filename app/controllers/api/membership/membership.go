package membership

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-membership/app/repositories/postgres"
	"go-membership/app/services"
	"go-membership/app/utilis"
	"go-membership/configs"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var (
	secret          = configs.GetJwtSecret()
	userRepository  = postgres.NewUserRepository()
	tokenRepository = postgres.NewTokenRepository()
	memberService   = services.NewMemberService()
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

	statusRegister := memberService.Register(registerForm.Name, registerForm.Email, registerForm.Password)

	switch statusRegister {
	case services.StatusRegisterEmailRegistered:
		response["message"] = "The email had been registered."
		c.JSON(http.StatusBadRequest, response)
		break
	case services.StatusRegisterEmailRegisterFailed:
		response["message"] = "The registration is failed."
		c.JSON(http.StatusBadRequest, response)
		break
	case services.StatusRegisterEncryptedPasswordFailed:
		response["message"] = "Some thing wrong."
		c.JSON(http.StatusBadRequest, response)
		break
	case services.StatusRegisterInternalError:
		response["message"] = "The system has some error."
		c.JSON(http.StatusInternalServerError, response)
		break
	case services.StatusRegisterSuccess:
		response["message"] = "Registered."
		c.JSON(http.StatusOK, response)
	}
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

	token, loginStatus := memberService.Login(loginForm.Email, loginForm.Password)
	switch loginStatus {
	case services.StatusLoginEmailNotExist:
		response["message"] = "User has this email is not existed."
		c.JSON(http.StatusBadRequest, response)
		break
	case services.StatusLoginPasswordInvalid:
		response["message"] = "Password is not correct."
		c.JSON(http.StatusBadRequest, response)
		break
	case services.StatusLoginInternalError:
		c.JSON(http.StatusInternalServerError, response)
		break
	case services.StatusLoginSignJwtFailed:
		response["message"] = "Some thing wrong."
		c.JSON(http.StatusBadRequest, response)
		break
	case services.StatusLoginInsertTokenFailed:
		response["message"] = "Some thing wrong."
		c.JSON(http.StatusBadRequest, response)
		break
	case services.StatusLoginSuccess:
		response["message"] = "Login success."
		response["token"] = token
		c.JSON(http.StatusBadRequest, response)
	}
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
