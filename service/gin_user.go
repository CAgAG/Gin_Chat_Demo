package service

import (
	"Chat_demo/cache"
	"Chat_demo/model"
	"Chat_demo/pkg/status"
	"Chat_demo/pkg/utils"
	"Chat_demo/serializer"
	"fmt"
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

// 使用context.ShouldBind来绑定JSON请求数据到一个结构体
type UserBaseService struct {
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
} // 用户输入的信息，会通过json解析到这个结构体

type UserSetPasswordService struct {
	UserName    string `form:"user_name" json:"user_name"`
	OldPassword string `form:"old_password" json:"old_password"`
	NewPassword string `form:"new_password" json:"new_password"`
}

func (service *UserBaseService) Register() serializer.Response {
	var new_user model.User
	var count int64

	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).First(&new_user).Count(&count)
	if count != 0 {
		return serializer.Response{
			Status: status.USER_EXIST,
			Msg:    status.TransCode(status.USER_EXIST),
		}
	}

	new_user = model.User{
		UserName: service.UserName,
	}
	if err := new_user.SetPassword(service.Password); err != nil {
		return serializer.Response{
			Status: status.USER_SET_PASSWORD_FAIL,
			Msg:    status.TransCode(status.USER_SET_PASSWORD_FAIL),
		}
	}
	model.DB.Create(&new_user)
	return serializer.Response{
		Status: status.USER_REGISTER_SUCCESS,
		Msg:    status.TransCode(status.USER_REGISTER_SUCCESS),
	}
}

func (service *UserBaseService) Login(context *gin.Context) serializer.Response {
	var login_user model.User
	var count int64

	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).First(&login_user).Count(&count)
	if count == 0 {
		return serializer.Response{
			Status: status.USER_NOT_EXIST,
			Msg:    status.TransCode(status.USER_NOT_EXIST),
		}
	}
	if !login_user.CheckPassword(service.Password) {
		return serializer.Response{
			Status: status.USER_INCORRECT_PASSWORD,
			Msg:    status.TransCode(status.USER_INCORRECT_PASSWORD),
		}
	}

	// Set session ================================================
	sv := cache.SessionValue{UserName: service.UserName, PassWord: service.Password}
	err := sv.SetPasswordAuth(context)
	if err != nil {
		return serializer.Response{
			Status: status.USER_SET_SESSION_FAIL,
			Msg:    status.TransCode(status.USER_SET_SESSION_FAIL),
		}
	}
	logging.Info("session password: " + fmt.Sprintf("%v", sv.GetAuth(context)))
	err = sv.SetUserName(context)
	if err != nil {
		return serializer.Response{
			Status: status.USER_SET_SESSION_FAIL,
			Msg:    status.TransCode(status.USER_SET_SESSION_FAIL),
		}
	}
	logging.Info("session user: " + fmt.Sprintf("%v", sv.GetUserName(context)))
	// Set session ================================================

	return serializer.Response{
		Status: status.USER_LOGIN_SUCCESS,
		Msg:    status.TransCode(status.USER_LOGIN_SUCCESS),
	}
}

func (service *UserSetPasswordService) SetPassWord(context *gin.Context) serializer.Response {
	// session
	sv := cache.SessionValue{UserName: service.UserName}
	logging.Info("session: " + fmt.Sprintf("%v", sv.GetAuth(context)))
	if sv.GetAuth(context) == "" {
		return serializer.Response{
			Status: status.USER_NOT_LOGIN,
			Msg:    status.TransCode(status.USER_NOT_LOGIN),
		}
	}

	if !utils.HashCompare(sv.GetAuth(context), service.OldPassword) {
		return serializer.Response{
			Status: status.USER_INCORRECT_PASSWORD,
			Msg:    status.TransCode(status.USER_INCORRECT_PASSWORD),
		}
	}

	update_user := model.User{UserName: service.UserName}
	err := update_user.SetPassword(service.NewPassword)
	if err != nil {
		return serializer.Response{
			Status: status.USER_SET_PASSWORD_FAIL,
			Msg:    status.TransCode(status.USER_SET_PASSWORD_FAIL),
		}
	}

	model.DB.Model(&model.User{}).Where("user_name=?", update_user.UserName).Update("password", update_user.Password)
	// 更新 Auth
	sv.PassWord = service.NewPassword
	err = sv.SetPasswordAuth(context)
	if err != nil {
		return serializer.Response{
			Status: status.USER_SET_SESSION_FAIL,
			Msg:    status.TransCode(status.USER_SET_SESSION_FAIL),
		}
	}

	return serializer.Response{
		Status: status.USER_SET_PASSWORD_SUCCESS,
		Msg:    status.TransCode(status.USER_SET_PASSWORD_SUCCESS),
	}
}

func (service *UserBaseService) Logout(context *gin.Context) serializer.Response {
	var new_user model.User
	var count int64

	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).First(&new_user).Count(&count)
	if count == 0 {
		return serializer.Response{
			Status: status.USER_NOT_EXIST,
			Msg:    status.TransCode(status.USER_NOT_EXIST),
		}
	}

	// 用户是否登录
	sv := cache.SessionValue{UserName: service.UserName}
	logging.Info(sv.AuthIsExist(context))
	if !sv.AuthIsExist(context) {
		return serializer.Response{
			Status: status.USER_NOT_LOGIN,
			Msg:    status.TransCode(status.USER_NOT_LOGIN),
		}
	}

	err := sv.DelAuth(context)
	if err != nil {
		return serializer.Response{
			Status: status.USER_LOGOUT_FAIL,
			Msg:    status.TransCode(status.USER_LOGOUT_FAIL),
		}
	}
	logging.Info(sv.AuthIsExist(context))
	return serializer.Response{
		Status: status.USER_LOGOUT_SUCCESS,
		Msg:    status.TransCode(status.USER_LOGOUT_SUCCESS),
	}
}
