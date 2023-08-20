package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"password-lock/models"
	"password-lock/service"
	"password-lock/utils"
)

func (c Controller) RegisterUser(ctx *gin.Context) {
	var user models.User

	// decoding json message to user model
	err := json.NewDecoder(ctx.Request.Body).Decode(&user)
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	_, err = c.service.RegisterUser(user)
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	secretKey, err := utils.GenerateRandomStringURLSafe()
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	err = service.SendEmail(user.EmailAddress, secretKey)
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	SendResponse(ctx, Response{
		Status:  http.StatusCreated,
		Message: "user registered successfully",
	})
	return

}

func (c Controller) Login(ctx *gin.Context) {

	var credentials models.User

	// decoding json message to user model
	err := json.NewDecoder(ctx.Request.Body).Decode(&credentials)
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	user, err := c.service.Authenticate(credentials)
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		})
		return
	}

	sessionKey, err := c.service.GenerateAndSaveSessionKey(user)
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	ctx.SetCookie("session", sessionKey, 600, "/", "", true, true)

	SendResponse(ctx, Response{
		Status:  http.StatusOK,
		Message: "successfully logged in",
	})
	return

}

func (c Controller) Logout(ctx *gin.Context) {
	sessionId, err := ctx.Cookie("session")
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	c.service.TerminateSession(sessionId)

	ctx.SetCookie("session", "", -1, "/", "", true, true)

	SendResponse(ctx, Response{
		Status:  http.StatusOK,
		Message: "successfully logged out",
	})
	return
}
