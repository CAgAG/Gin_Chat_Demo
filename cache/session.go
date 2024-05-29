package cache

import (
	"Chat_demo/pkg/utils"
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	session_redis "github.com/gin-contrib/sessions/redis"
	// "github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type SessionValue struct {
	UserName string
	PassWord string
}

// 设置 redis中的 session 的过期时间, 单位 s
var RedisSessionLiveTime int = 60 * 60 * 24 * 7

func NewSession() (sessions.Store, error) {
	// 注册结构体，使其可以跨路由存取
	gob.Register(SessionValue{})

	// 使用 redis 来保存 session
	store, err := session_redis.NewStore(10, "tcp", RedisAddr, "", []byte("secret"))
	// store := cookie.NewStore([]byte("secret"))

	// store, err := MyRedisStore(10, "tcp", RedisAddr, "", []byte("secret"))
	store.Options(sessions.Options{MaxAge: RedisSessionLiveTime}) // session 存活时间
	return store, err
}

// 密码 ===========================================
func (sv *SessionValue) SetPasswordAuth(context *gin.Context) error {
	if sv.UserName == "" {
		return &utils.MyError{Message: "用户名不能为空"}
	}

	session := sessions.Default(context)
	auth, err := utils.HashStr(sv.PassWord)
	if err != nil {
		return &utils.MyError{Message: "加密失败"}
	}

	session.Set(fmt.Sprintf("%s-Auth", sv.UserName), auth)
	err = session.Save()
	return err
}

func (sv *SessionValue) GetAuth(context *gin.Context) string {
	session := sessions.Default(context)
	user_name := session.Get(fmt.Sprintf("%s-Auth", sv.UserName))
	if user_name == nil {
		return ""
	}
	return user_name.(string)
}

func (sv *SessionValue) AuthIsExist(context *gin.Context) bool {
	session := sessions.Default(context)
	user_name := session.Get(fmt.Sprintf("%s-Auth", sv.UserName))
	if user_name == nil {
		return false
	}
	return true
}

func (sv *SessionValue) DelAuth(context *gin.Context) error {
	session := sessions.Default(context)
	session.Delete(fmt.Sprintf("%s-Auth", sv.UserName))
	err := session.Save()
	return err
}

// 用户名 ==================================
func (sv *SessionValue) GetUserName(context *gin.Context) string {
	session := sessions.Default(context)
	user_name := session.Get("UserName")
	if user_name == nil {
		return ""
	}
	return user_name.(string)
}

func (sv *SessionValue) SetUserName(context *gin.Context) error {
	if sv.UserName == "" {
		return &utils.MyError{Message: "用户名不能为空"}
	}

	session := sessions.Default(context)
	session.Set("UserName", sv.UserName)
	err := session.Save()
	return err
}

func (sv *SessionValue) DelUserName(context *gin.Context) error {
	session := sessions.Default(context)
	session.Delete("UserName")
	err := session.Save()
	return err
}
