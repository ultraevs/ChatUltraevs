package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthPage(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{"page": "auth"})
}

func ForgotPage(context *gin.Context) {
	context.HTML(http.StatusOK, "forgot.html", gin.H{"page": "forgot"})
}

func NewPass(context *gin.Context) {
	context.HTML(http.StatusOK, "newpass.html", gin.H{"page": "newpass"})
}
