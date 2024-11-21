package service

import (
	"goIM/dao"
	"goIM/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type groupService struct{}

var GroupService = new(groupService)

func (g *groupService) CreateGroup(ctx *gin.Context) {
	owner := ctx.Request.FormValue("userId")
	ownerId, err := strconv.Atoi(owner)

	if err != nil {
		zap.S().Info("ownerId 类型转换失败", err)
		jsonOutPut(ctx, -1, "参数错误")
		return
	}

	typeStr := ctx.Request.FormValue("type")
	typeId, err := strconv.Atoi(typeStr)
	if err != nil {
		zap.S().Info("type 类型转换失败", err)
		jsonOutPut(ctx, -1, "参数错误")
		return
	}

	img := ctx.Request.FormValue("img")
	name := ctx.Request.FormValue("name")
	desc := ctx.Request.FormValue("desc")

	if ownerId == 0 || name == "" || typeId == 0 {
		jsonOutPut(ctx, -1, "参数错误")
		return
	}

	group := models.Group{}
	group.OwnerId = uint(ownerId)
	group.Type = typeId
	group.Name = name

	if img != "" {
		group.Image = img
	}
	if desc != "" {
		group.Desc = desc
	}

	code, err := dao.CreateGroup(group)

	if code != 0 {
		jsonOutPut(ctx, -1, string(err.Error()))
		return
	}

	jsonOutPut(ctx, 0, "创建群组成功")
}

func (g *groupService) GroupList(ctx *gin.Context) {
	owner := ctx.Request.FormValue("userId")
	ownerId, err := strconv.Atoi(owner)

	if err != nil {
		zap.S().Info("类型转换失败", err)
		jsonOutPut(ctx, -1, "参数错误")
		return
	}

	if ownerId == 0 {
		jsonOutPut(ctx, -1, "未登录 请重新登录")
		return
	}
	groupList, err := dao.GetGroupList(uint(ownerId))
	if err != nil {
		zap.S().Info("获取群组列表失败\n", err)
		jsonOutPut(ctx, -1, "获取群组列表失败")
		return
	}
	jsonOutPut(ctx, 0, "获取群组列表成功", groupList)
}

func (g *groupService) JoinGroup(ctx *gin.Context) {
	groupName := ctx.Request.FormValue("group_name")
	if groupName == "" {
		zap.S().Info("参数错误，群聊名称不能为空\n")
		jsonOutPut(ctx, -1, "请输入群聊名称")
		return
	}
	userIdStr := ctx.Request.FormValue("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		zap.S().Info("类型转换失败", err, "\n")
		jsonOutPut(ctx, -1, "参数错误")
		return
	}
	if userId == 0 {
		zap.S().Info("账号未登录 请重新登录\n")
		jsonOutPut(ctx, -1, "账号未登录 请重新登录")
		return
	}

	code, err := dao.JoinGroup(uint(userId), groupName)
	if code != 0 {
		zap.S().Info("加入群组错误 请重试:", err.Error(), "\n")
		jsonOutPut(ctx, -1, err.Error())
		return
	}

	jsonOutPut(ctx, 0, "加入群聊成功")
}
