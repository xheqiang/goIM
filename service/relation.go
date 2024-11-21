package service

import (
	"goIM/dao"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type userInfo struct {
	Name     string
	Avatar   string
	Gender   string
	Phone    string
	Email    string
	Identity string
}

type relationService struct{}

var RelationService = new(relationService)

func (relationService *relationService) FriendList(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Request.FormValue("userId"))
	users, err := dao.FriendList(uint(id))
	if err != nil {
		zap.S().Info("获取好友列表失败", err)
		jsonOutPut(ctx, -1, "获取好友列表失败")
		return
	}

	if users == nil {
		zap.S().Info("好友列表为空", err)
		jsonOutPut(ctx, -1, "好友列表为空")
		return
	}

	friendList := make([]userInfo, 0)

	for _, v := range *users {
		friend := userInfo{
			Name:     v.Name,
			Avatar:   v.Avatar,
			Gender:   v.Gender,
			Phone:    v.Phone,
			Email:    v.Email,
			Identity: v.Identity,
		}
		friendList = append(friendList, friend)
	}
	jsonOutPut(ctx, 0, "获取好友列表成功", friendList)
}

func (relationService *relationService) AddFriend(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Request.FormValue("userId"))
	if err != nil {
		zap.S().Info("参数错误", err)
		jsonOutPut(ctx, -1, "参数错误")
	}

	target := ctx.Request.FormValue("target")
	if targetId, err := strconv.Atoi(target); err == nil {
		_, err := dao.AddFriend(uint(userId), uint(targetId))
		if err != nil {
			zap.S().Info("添加好友失败", err)
			jsonOutPut(ctx, -1, "添加好友失败")
			return
		}
	} else {
		_, err := dao.AddFriendByName(uint(userId), target)
		if err != nil {
			zap.S().Info("添加好友失败", err)
			jsonOutPut(ctx, -1, "添加好友失败")
			return
		}
	}

	jsonOutPut(ctx, 0, "添加好友成功")
}
