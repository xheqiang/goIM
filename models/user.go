package models

import (
	"time"
)

type User struct {
	Model
	Name          string    `gorm:"type:varchar(25)" json:"name"`
	Password      string    `gorm:"type:varchar(100)" json:"password"`
	Avatar        string    `gorm:"type:varchar(100)" json:"avatar"`
	Gender        string    `gorm:"type:varchar(6); default:male;comment 'male表示男,famale表示女'" json:"gender"`
	Phone         string    `gorm:"type:varchar(25)" valid:"matches(^1[3-9]{1}\\d{9}$)" json:"phone"`
	Email         string    `gorm:"type:varchar(25)" valid:"email" json:"email"`
	Identity      string    `gorm:"type:varchar(100)" json:"identity"`
	ClientIp      string    `gorm:"type:varchar(100)" valid:"ipv4" json:"client_ip"`
	ClientPort    string    `gorm:"type:varchar(25)" json:"client_port"`
	Salt          string    `gorm:"type:varchar(100)" json:"salt"` //加密盐值
	LoginTime     time.Time `gorm:"type:datetime" json:"login_time"`
	HeartBeatTime time.Time `gorm:"type:datetime" json:"heart_beat_time"`
	LoginOutTime  time.Time `gorm:"type:datetime" json:"login_out_time"`
	IsLoginOut    bool      `gorm:"type:tinyint(1)" json:"is_login_out"`
	DeviceInfo    string    `gorm:"type:varchar(100)" json:"device_info"` // 登录设备ID
}

func (User) TableName() string {
	return "user"
}
