package router

import (
	"github.com/gin-gonic/gin"
)

func InitGin() *gin.Engine {
	r := gin.New()
	initRouter(r)
	return r
}

func initRouter(r *gin.Engine) {
	r.GET("/ping", ping)
	r.POST("/user/login", login)
	r.POST("user/register", register)
	r.GET("/user/info", getUserInfo)
	r.GET("/user/edit", editUserInfo)
	//r.POST("/user/logout", logout)
	//r.GET("/user/find", findByUsername)
	//r.POST("/user/profile", profilePicUpdate)
	//r.POST("/user/nick", nickNameUpdate)
}
