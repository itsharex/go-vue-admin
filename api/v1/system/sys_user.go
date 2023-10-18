package system

import (
	"github/shansec/go-vue-admin/global"
	"github/shansec/go-vue-admin/model/common/response"
	"github/shansec/go-vue-admin/model/system"
	systemReq "github/shansec/go-vue-admin/model/system/request"
	systemRes "github/shansec/go-vue-admin/model/system/response"
	"github/shansec/go-vue-admin/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BaseApi struct{}

// @Tags SysUser
// @Summary 用户注册账号
// @Produce json
// @Param data body systemReq.Register { 用户名、密码、昵称、手机号 }
// @Success 200
// @Router /base/register POST

func (b *BaseApi) Register(c *gin.Context) {
	var register systemReq.Register
	_ = c.ShouldBindJSON(&register)
	if err := utils.Verify(register, utils.RegisterVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := &system.SysUser{
		Username:  register.Username,
		NickName:  register.NickName,
		Password:  register.Password,
		HeaderImg: register.HeaderImg,
		Phone:     register.Phone,
	}
	userResgisterRes, err := userService.Register(*user)
	if err != nil {
		global.MAY_LOGGER.Error("注册失败", zap.Error(err))
		response.FailWithDetailed(systemRes.SysUserResponse{User: userResgisterRes}, "注册失败", c)
	} else {
		response.OkWithDetailed(systemRes.SysUserResponse{User: userResgisterRes}, "注册成功", c)
	}
}

// @Tags SysUser
// @Summary 用户登录
// @Produce json
// @Param data body systemReq.Login { 用户名、密码 }
// @Success 200
// @Router /base/login POST

func (b *BaseApi) Login(c *gin.Context) {
	var login systemReq.Login
	_ = c.ShouldBindJSON(&login)
	if err := utils.Verify(login, utils.LoginVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	u := &system.SysUser{Username: login.Username, Password: login.Password}
	if user, err := userService.Login(u); err != nil {
		global.MAY_LOGGER.Error("登陆失败，用户名不存在或者密码错误", zap.Error(err))
		response.FailWithMessage("用户名不存在或者密码错误", c)
	} else {
		b.TokenNext(c, *user)
	}
}

// TokenNext 登录成功签发 jwt
func (b *BaseApi) TokenNext(c *gin.Context, user system.SysUser) {
	jwt := &utils.JWT{SigningKey: []byte(global.MAY_CONFIG.JWT.SigningKey)}
	claims := jwt.CreateClaims(systemReq.BaseClaims{
		UUID:     user.UUID,
		ID:       user.ID,
		NickName: user.NickName,
		Username: user.Username,
	})
	token, err := jwt.CreateToken(claims)
	if err != nil {
		global.MAY_LOGGER.Error("获取 token 失败", zap.Error(err))
		response.FailWithMessage("获取 token 失败", c)
		return
	}
	response.OkWithDetailed(systemRes.Login{
		User:      user,
		Token:     token,
		ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
	}, "登录成功", c)
}

// @Tags SysUser
// @Summary 修改密码
// @Produce json
// @Param data body systemReq.ChangePassword { 原密码，新密码 }
// @Success 200
// @Router /base/login POST

func (b *BaseApi) ModifyPassword(c *gin.Context) {
	var modifyPassword systemReq.ChangePassword
	_ = c.ShouldBindJSON(&modifyPassword)

	if err := utils.Verify(modifyPassword, utils.ChangePasswordVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUseID(c)
	changePassword := &system.SysUser{MAY_MODEL: global.MAY_MODEL{ID: uid}, Password: modifyPassword.Password}
	if _, err := userService.ChangePassword(changePassword, modifyPassword.NewPassword); err != nil {
		global.MAY_LOGGER.Error("修改失败", zap.Error(err))
		response.FailWithMessage("修改失败，原密码不正确", c)
	} else {
		response.OkWithMessage("修改成功，请重新登录", c)
	}
}

// @Tags SysUser
// @Summary 获取用户信息
// @Produce json
// @Param data body systemRes.SysUserResponse { uuid }
// @Success 200
// @Router /user/getUserInfo GET

func (b *BaseApi) GetUserInfo(c *gin.Context) {
	uuid := utils.GetUseUuid(c)
	if user, err := userService.GetUserInformation(uuid); err != nil {
		global.MAY_LOGGER.Error("获取用户信息失败", zap.Error(err))
		response.FailWithMessage("获取用户信息失败，用户信息不存在", c)
	} else {
		response.OkWithDetailed(systemRes.SysUserResponse{
			User: *user,
		}, "获取用户信息成功", c)
	}
}
