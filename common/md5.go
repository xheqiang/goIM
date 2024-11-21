package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// Md5encoder 加密后返回小写值
func Md5encoder(code string) string {
	m := md5.New()
	io.WriteString(m, code)
	return hex.EncodeToString(m.Sum(nil))
}

// Md5StrToUpper 加密后返回大写
func Md5StrToUpper(code string) string {
	return strings.ToUpper(Md5encoder(code))
}

// SaltPassWord 密码加盐
func SaltPassword(pwd string, salt string) string {
	saltPwd := fmt.Sprintf("%s$%s", Md5encoder(pwd), salt)
	return saltPwd
}

// CheckPassword 核验密码
func CheckPassword(encodePwd, salt, pwd string) bool {
	return encodePwd == SaltPassword(pwd, salt)
}
