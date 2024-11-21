package dao

import (
	"errors"
	"goIM/global"
	"goIM/models"

	"go.uber.org/zap"
)

// FriendList 获取好友列表
func FriendList(userId uint) (*[]models.User, error) {
	relation := make([]models.Relation, 0)
	if tx := global.DB.Where("owner_id = ? and type = ?", userId, 1).Find(&relation); tx.RowsAffected == 0 {
		zap.S().Info("未找到Relation数据")
		return nil, errors.New("未查到好友关系")
	}
	userIds := make([]uint, 0)
	for _, v := range relation {
		userIds = append(userIds, v.TargetId)
	}
	users := make([]models.User, 0)
	if tx := global.DB.Where("id in ?", userIds).Find(&users); tx.RowsAffected == 0 {
		zap.S().Info("未找到Relation好友数据")
		return nil, errors.New("未查到好友")
	}

	return &users, nil
}

// 通过昵称添加好友
func AddFriendByName(userId uint, targetName string) (int, error) {
	user, err := FindUserByName(targetName)
	if err != nil {
		return -1, errors.New("未找到用户")
	}

	if user.Id == 0 {
		zap.S().Info("未查询到用户")
		return -1, errors.New("未查询到用户")
	}

	return AddFriend(userId, user.Id)
}

// AddFriend 加好友
func AddFriend(userId, targetId uint) (int, error) {
	if userId == targetId {
		return -2, errors.New("不能添加自己为好友")
	}
	targetUser, err := FindUserById(uint(targetId))
	if err != nil {
		return -1, errors.New("未找到用户")
	}

	if targetUser.Id == 0 {
		zap.S().Info("未找到用户")
		return -1, errors.New("未找到用户")
	}

	relation := models.Relation{}

	if tx := global.DB.Where("owner_id = ? and target_id = ? and type =?", userId, targetId, 1).First(&relation); tx.RowsAffected == 1 {
		zap.S().Info("好友已经存在")
		return 0, errors.New("好友已经存在")
	}

	if tx := global.DB.Where("owner_id = ? and target_id = ? and type =?", targetId, userId, 1).First(&relation); tx.RowsAffected == 1 {
		zap.S().Info("好友已经存在")
		return 0, errors.New("好友已经存在")
	}

	// 开启事务 提交好友关系
	tx := global.DB.Begin()

	relation.OwnerId = uint(userId)
	relation.TargetId = uint(targetId)
	relation.Type = 1

	if t := tx.Create(&relation); t.RowsAffected == 0 {
		zap.S().Info("好友添加失败")

		// 事务回滚
		tx.Rollback()
		return -1, errors.New("好友添加失败")
	}

	relation = models.Relation{}
	relation.OwnerId = uint(targetId)
	relation.TargetId = uint(userId)
	relation.Type = 1

	if t := tx.Create(&relation); t.RowsAffected == 0 {
		zap.S().Info("好友添加失败")

		// 事务回滚
		tx.Rollback()
		return -1, errors.New("好友添加失败")
	}

	tx.Commit()
	return 1, nil
}

func FindUses(groupId uint) (*[]uint, error) {
	relation := make([]models.Relation, 0)
	if tx := global.DB.Where("target_id = ? and type = ?", groupId, 2).Find(&relation); tx.RowsAffected == 0 {
		return nil, errors.New("未查到成员信息")
	}

	userIds := make([]uint, 0)
	for _, v := range relation {
		userIds = append(userIds, v.OwnerId)
	}
	return &userIds, nil
}
