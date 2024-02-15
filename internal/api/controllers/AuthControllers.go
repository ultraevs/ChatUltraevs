package controller

import (
	"app/internal/database"
	"app/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

// UserCreate создает нового пользователя.
// @Summary Создать нового пользователя
// @Description Создает нового пользователя с предоставленным email, паролем и именем.
// @Accept json
// @Produce json
// @Param request body model.UserCreateRequest true "Запрос на создание пользователя"
// @Success 200 {object} model.CodeResponse "Пользователь успешно создан"
// @Failure 400 {object} model.ErrorResponse "Не удалось создать пользователя"
// @Tags Auth
// @Router /v1/create [post]
func UserCreate(context *gin.Context) {
	var body struct {
		Email    string
		Password string
		Name     string
	}
	if context.Bind(&body) != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	hashPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
	}
	_, err = database.Db.Exec("INSERT INTO chat_users (email, password, name) VALUES ($1, $2, $3)", body.Email, string(hashPass), body.Name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login Вход в аккаунт.
// @Summary Логин
// @Description Авторизует пользователя с предоставленным email и паролем.
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Запрос на авторизацию пользователя"
// @Success 200 {object} model.CodeResponse "Пользователь авторизован"
// @Failure 400 {object} model.ErrorResponse "Не удалось авторизовать пользователя"
// @Tags Auth
// @Router /v1/login [post]
func Login(context *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	if context.Bind(&body) != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	fmt.Println(body)
	var form model.LoginRequest
	err := database.Db.QueryRow("SELECT email, password FROM chat_users WHERE email = $1", body.Email).Scan(&form.Email, &form.Password)
	if err == nil {
		fmt.Println(bcrypt.CompareHashAndPassword([]byte("$2a$10$QQIbvr7VuxKX5YNG1Upu/uQp61zFsvjz3pRLgP4hKWI5ve3s19I1m"), []byte("test")))
		if err := bcrypt.CompareHashAndPassword([]byte(form.Password), []byte(body.Password)); err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
			return
		}
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "User isn't founded"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": form.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	context.SetSameSite(http.SameSiteLaxMode)
	var domain string
	if os.Getenv("SETUP_TYPE") == "local" {
		domain = "localhost"
	} else {
		domain = "chat.ultraevs.ru"
	}
	context.SetCookie("Authorization", tokenString, 3600*24*30, "/", domain, false, false)
	context.Set("Authorization", tokenString)
	context.JSON(http.StatusOK, gin.H{"response": "success"})
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) model.EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, os.Getenv("SMTP_HOST"))
	return e.Send(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), smtpAuth)
}

// ForgotPassword Восстановление аккаунта.
// @Summary Восстановление
// @Description Инициирует восстановление пароля по email.
// @Accept json
// @Produce json
// @Param request body model.ForgotRequest true "Запрос на инициацию восстановления пользователя"
// @Success 200 {object} model.CodeResponse "Процесс восстановления начат"
// @Failure 400 {object} model.ErrorResponse "Не удалось начать процесс восстановления"
// @Tags Auth
// @Router /v1/forgot [post]
func ForgotPassword(context *gin.Context) {
	var body struct {
		Email string
	}
	if context.Bind(&body) != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	sender := NewGmailSender("Chat Ultraevs", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"))
	fmt.Println(body.Email)
	subject := "Восстановление пароля"
	token := body.Email
	var link string
	if os.Getenv("SETUP_TYPE") == "local" {
		link = "http://localhost:8083/newpass?token=" + token
	} else {
		link = "https://chat.ultraevs.ru/newpass?token=" + token
	}
	content := fmt.Sprintf(`
	<h1>Здравствуйте, вы сделали запрос на восстановление пароля. Чтобы сменить пароль, перейдите по ссылке:</h1>
	<a href="%s">Восстановить пароль</a>
	`, link)
	to := []string{body.Email}
	err := sender.SendEmail(subject, content, to, nil, nil)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error with sending email"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"success": "Email sent"})
}

//// NewPassword Создание нового пароля.
//// @Summary Новый пароль
//// @Description Новый пароль для восстановления пароля по email.
//// @Accept json
//// @Produce json
//// @Param token query string true "Токен восстановления пароля"
//// @Success 200 {object} model.CodeResponse "Процесс восстановления завершен"
//// @Failure 400 {object} model.ErrorResponse "Не удалось завершить процесс восстановления"
//// @Tags Auth
//// @Router /v1/newpass [get]
//func NewPassword(context *gin.Context) {
//	receivedEmail := context.Query("token")
//	if receivedEmail == "" {
//		context.JSON(http.StatusInternalServerError, gin.H{"error": "Email is missing"})
//		return
//	}
//	context.JSON(http.StatusOK, gin.H{"email": receivedEmail})
//}

// PostNewPassword Добавление нового пароля.
// @Summary Добавление нового пароля в БД
// @Description Сохранение нового пароля .
// @Accept json
// @Produce json
// @Param request body model.NewPassword true "Новый пароль пользователя"
// @Success 200 {object} model.CodeResponse "Пароль сохранен"
// @Failure 400 {object} model.ErrorResponse "Не удалось сохранить пароль"
// @Tags Auth
// @Router /v1/newpass [post]
func PostNewPassword(context *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	if context.Bind(&body) != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate password hash"})
		return
	}

	var existingPassword string
	err = database.Db.QueryRow("SELECT password FROM chat_users WHERE email = $1", body.Email).Scan(&existingPassword)
	if err == nil {
	} else {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	_, err = database.Db.Exec("UPDATE chat_users SET password = $1 WHERE email = $2", string(hashPass), body.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}
	sender := NewGmailSender("Chat Ultraevs", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"))

	subject := "Восстановление пароля"
	content := fmt.Sprintf(`
	<p>Здравствуйте, пароль на вашем аккануте был изменен.</p>
	<p>Если это были не вы - обратитесь в поддержку нашего сайта: smyakneksbimisis@gmail.com</p>
	`)
	to := []string{body.Email}
	err = sender.SendEmail(subject, content, to, nil, nil)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error with sending email"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
