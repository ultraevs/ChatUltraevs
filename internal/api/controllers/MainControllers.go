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

func NewPassPage(context *gin.Context) {
	context.HTML(http.StatusOK, "newpass.html", gin.H{"page": "newpass"})
}

func MainPage(context *gin.Context) {
	context.HTML(http.StatusOK, "home.html", gin.H{"page": "home"})
}
