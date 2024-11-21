package service

import (
	"fmt"
	"goIM/common"
	"goIM/dao"
	"goIM/middlewear"
	"goIM/models"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type userService struct{}

var UserService = new(userService)

func (userService *userService) List(ctx *gin.Context) {
	userList, err := dao.GetUserList()
	if err != nil {
		jsonOutPut(ctx, -1, "获取用户列表失败")
		return
	}
	jsonOutPut(ctx, 0, "获取用户列表成功", userList)
}

func (userService *userService) Register(ctx *gin.Context) {
	user := models.User{}
	user.Name = ctx.Request.FormValue("name")
	password := ctx.Request.FormValue("password")
	repassword := ctx.Request.FormValue("repassword")

	if user.Name == "" || password == "" || repassword == "" {
		jsonOutPut(ctx, -1, "用户名或密码不能为空")
		return
	}

	if password != repassword {
		jsonOutPut(ctx, -1, "两次密码不一致")
		return
	}

	// 查询用户是否存在
	userInfo, _ := dao.FindUserByName(user.Name)
	if userInfo != nil {
		jsonOutPut(ctx, -1, "用户已存在")
		return
	}

	// 生成盐值
	salt := fmt.Sprintf("%d", rand.Int31())

	// 生成密码
	user.Password = common.SaltPassword(password, salt)
	user.Salt = salt
	user.LoginTime = time.Now()
	user.LoginOutTime = time.Now()
	user.HeartBeatTime = time.Now()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := dao.CreateUser(user)
	if err != nil {
		jsonOutPut(ctx, -1, "注册失败")
		return
	}

	data := map[string]interface{}{
		"name":     user.Name,
		"password": password,
	}

	jsonOutPut(ctx, 0, "注册成功", data)
}

func (userService *userService) Login(ctx *gin.Context) {
	name := ctx.Request.FormValue("name")
	password := ctx.Request.FormValue("password")

	data, err := dao.FindUserByName(name)
	if err != nil {
		jsonOutPut(ctx, -1, "登录失败, 请检查用户名密码")
		return
	}
	if data.Name == "" {
		jsonOutPut(ctx, -1, "用户不存在")
		return
	}

	//数据库密码保存是使用md5密文的， 验证密码时，将用户传递密码再次加密，然后进行对比
	if !common.CheckPassword(data.Password, data.Salt, password) {
		jsonOutPut(ctx, -1, "密码错误")
		return
	}

	token, err := middlewear.GenerateToken(data.Id, "yk")
	if err != nil {
		zap.S().Info("生成token失败", err)
		jsonOutPut(ctx, -1, "生成token失败")
		return
	}
	rep := make(map[string]interface{})
	rep["token"] = token
	rep["userId"] = data.Id
	jsonOutPut(ctx, 0, "登录成功", rep)
}

func (userService *userService) Update(ctx *gin.Context) {
	user := models.User{}

	uid, err := strconv.Atoi(ctx.Request.FormValue("userId"))
	if err != nil {
		jsonOutPut(ctx, -1, "参数错误")
		return
	}

	user.Id = uint(uid)
	name := ctx.Request.FormValue("name")
	password := ctx.Request.FormValue("password")
	email := ctx.Request.FormValue("email")
	phone := ctx.Request.FormValue("phone")
	avatar := ctx.Request.FormValue("avatar")
	gender := ctx.Request.FormValue("gender")

	if name != "" {
		user.Name = name
	}
	if password != "" {
		salt := fmt.Sprintf("%d", rand.Int31())
		user.Salt = salt
		user.Password = common.SaltPassword(password, salt)
	}
	if email != "" {
		user.Email = email
	}
	if phone != "" {
		user.Phone = phone
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if gender != "" {
		user.Gender = gender
	}

	rep, err := dao.UpdateUser(user)
	if err != nil {
		zap.S().Info("用户信息更新失败：", err)
		jsonOutPut(ctx, -1, "用户信息更新失败")
		return
	}

	jsonOutPut(ctx, 0, "用户信息修改成功", rep)
}

func (userService *userService) Delete(ctx *gin.Context) {
	user := models.User{}

	uid, err := strconv.Atoi(ctx.Request.FormValue("userId"))
	if err != nil {
		zap.S().Info("类型转换失败", err)
		jsonOutPut(ctx, -1, "参数错误")
		return
	}
	user.Id = uint(uid)
	err = dao.DeleteUser(user)
	if err != nil {
		zap.S().Info("注销用户失败", err)
		jsonOutPut(ctx, 0, "注销用户失败")
		return
	}

	jsonOutPut(ctx, 0, "注销用户成功")

}
