package routes

import (
	"github.com/gin-gonic/gin"
	"go-membership/app/controllers/api/membership"
)

func Listen() {
	router := gin.New()
	router.POST("/api/register", membership.Registration)
	router.POST("/api/login", membership.Login)
	router.POST("/api/send-password-reset-email", membership.SendPasswordResetEmail)
	router.POST("/api/reset-password", membership.ResetPassword)
	router.Run()
}
