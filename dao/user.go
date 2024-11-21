package dao

import (
	"errors"
	"goIM/common"
	"goIM/global"
	"goIM/models"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type UserInfo struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Gender     string `json:"gender"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	IsLoginOut bool   `json:"is_login_out"`
	LoginTime  string `json:"login_time"`
}

// GetUserList 获取用户列表 所有用户
func GetUserList() ([]UserInfo, error) {
	var users []models.User

	if tx := global.DB.Find(&users); tx.RowsAffected == 0 {
		return nil, errors.New("获取用户列表失败")
	}

	var userList []UserInfo

	for _, user := range users {
		userInfo := UserInfo{
			Id:         user.Id,
			Name:       user.Name,
			Avatar:     user.Avatar,
			Gender:     user.Gender,
			Phone:      user.Phone,
			Email:      user.Email,
			IsLoginOut: user.IsLoginOut,
			LoginTime:  user.LoginTime.Format("2006-01-02 15:04:05"),
		}
		userList = append(userList, userInfo)
	}

	return userList, nil
}

// 查询用户:根据昵称，根据电话，根据邮件

// FindUserByNameAndPwd 昵称和密码查询
func FindUserByNameAdnPwd(name string, password string) (*models.User, error) {
	user := models.User{}
	if tx := global.DB.Where("name = ? AND password = ?", name, password).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("未查询到用户")
	}

	// 登录识别
	t := strconv.Itoa(int(time.Now().Unix()))

	temp := common.Md5encoder(t)

	if tx := global.DB.Model(&user).Update("identity", temp); tx.RowsAffected == 0 {
		return nil, errors.New("更新用户identity失败")
	}

	return &user, nil
}

// FindUserByName 根据用户名查询用户
func FindUserByName(name string) (*models.User, error) {
	user := models.User{}
	if tx := global.DB.Where("name = ?", name).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("未查询到用户")
	}
	return &user, nil
}

// FindUser 根据用户名查询用户 注册判断是否已经注册
func FindUser(name string) (*models.User, error) {
	user := models.User{}
	if tx := global.DB.Where("name = ?", name).First(&user); tx.RowsAffected == 1 {
		return nil, errors.New("当前用户名已经存在")
	}
	return &user, nil
}

// FindUser 根据用户名查询用户 注册判断是否已经注册
func FindUserById(id uint) (*models.User, error) {
	user := models.User{}
	if tx := global.DB.Where("id = ?", id).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("用户名未找到")
	}
	return &user, nil
}

// FindUserByPhone 根据用户电话查询用户
func FindUserByPhone(phone string) (*models.User, error) {
	user := models.User{}
	if tx := global.DB.Where("phone = ?", phone).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("未找到记录")
	}
	return &user, nil
}

// FindUserByEmail 根据用户电话查询用户
func FindUserByEmail(email string) (*models.User, error) {
	user := models.User{}
	if tx := global.DB.Where("email = ?", email).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("未找到记录")
	}
	return &user, nil
}

// CreateUser 新建用户
func CreateUser(user models.User) (*models.User, error) {
	if tx := global.DB.Create(&user); tx.RowsAffected == 0 {
		zap.S().Info("新增用户失败")
		return nil, errors.New("新增用户失败")
	}
	return &user, nil
}

// UpdateUser 更新用户
func UpdateUser(user models.User) (*models.User, error) {
	tx := global.DB.Model(&user).Updates(models.User{
		Name:     user.Name,
		Password: user.Password,
		Gender:   user.Gender,
		Phone:    user.Phone,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Salt:     user.Salt,
	})
	if tx.RowsAffected == 0 {
		zap.S().Info("更新用户失败")
		return nil, errors.New("更新用户失败")
	}
	return &user, nil
}

// DeleteUser 删除用户
func DeleteUser(user models.User) error {
	if tx := global.DB.Delete(&user); tx.RowsAffected == 0 {
		zap.S().Info("删除用户失败")
		return errors.New("删除用户失败")
	}
	return nil
}
