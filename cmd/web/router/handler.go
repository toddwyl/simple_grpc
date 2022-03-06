package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"simple_grpc/cmd/web/http_service"
	errcode "simple_grpc/internal/error"
	"simple_grpc/internal/facade"
	"simple_grpc/internal/resp"
)

func ping(context *gin.Context) {
	context.JSON(200, gin.H{"message": "pong"})
}

func login(c *gin.Context) {
	response := resp.NewResponse(c)
	param := facade.LoginUserRequest{}
	err := c.ShouldBind(&param)
	if err != nil {
		logrus.Errorf("app.Login errs: %v", err)
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}
	svc := http_service.New(c.Request.Context())
	loginResponse, err := svc.Login(&param)
	if err != nil {
		logrus.Errorf("app.Login errs: %v", err)
		response.ToErrorResponse(errcode.ErrorUserLogin)
		return
	}
	c.SetCookie("session_id", loginResponse.SessionID, 3600, "/", "", false, true)
	response.ToResponse("Login Succeed.", loginResponse)
	return
}

func register(c *gin.Context) {
	response := resp.NewResponse(c)
	param := facade.RegisterUserRequest{}
	err := c.ShouldBind(&param)
	if err != nil {
		logrus.Errorf("app.Register errs: %v", err)
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}
	svc := http_service.New(c.Request.Context())
	loginResponse, err := svc.Register(&param)
	if err != nil {
		logrus.Errorf("app.Register errs: %v", err)
		response.ToErrorResponse(errcode.ErrorUserRegister)
		return
	}
	response.ToResponse("Register Succeed.", loginResponse)
	return
}

func getUserInfo(c *gin.Context) {
	response := resp.NewResponse(c)
	param := facade.GetUserRequest{}

	err := c.ShouldBind(&param)
	if err != nil {
		logrus.Errorf("app.Register errs: %v", err)
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := http_service.New(c.Request.Context())
	getResponse, err := svc.GetUser(&param)
	if err != nil {
		logrus.Errorf("app.Get errs: %v", err)
		response.ToErrorResponse(errcode.ErrorUserGet)
		return
	}
	response.ToResponse("Get Succeed.", getResponse)
	return
}

func editUserInfo(c *gin.Context) {
	response := resp.NewResponse(c)
	param := facade.EditUserRequest{}
	err := c.ShouldBind(&param)
	if err != nil {
		logrus.Errorf("app.Edit errs: %v", err)
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	//// Get sessionID set by Auth middleware
	//sessionId, _ := c.Get("session_id")
	//param.SessionId = fmt.Sprintf("%v", sessionId)

	svc := http_service.New(c.Request.Context())
	editResponse, err := svc.EditUser(&param)
	if err != nil {
		logrus.Errorf("app.Edit errs: %v", err)
		response.ToErrorResponse(errcode.ErrorUserEdit)
		return
	}
	response.ToResponse("Edit Succeed.", editResponse)
	return
}
