package dao

import (
	"errors"
	"goIM/global"
	"goIM/models"
)

func CreateGroup(group models.Group) (int, error) {
	com := models.Group{}

	// 查询是否已经存在
	if tx := global.DB.Where("name =?", group.Name).First(&com); tx.RowsAffected == 1 {
		return -1, errors.New("当前群已经存在")
	}

	tx := global.DB.Begin()
	if t := tx.Create(&group); t.RowsAffected == 0 {
		tx.Rollback()
		return -1, errors.New("创建群失败")
	}

	relation := models.Relation{}
	relation.TargetId = group.OwnerId
	relation.TargetId = group.Id
	relation.Type = 2
	if t := tx.Create(&relation); t.RowsAffected == 0 {
		tx.Rollback()
		return -1, errors.New("群记录创建失败")
	}
	tx.Commit()
	return 0, nil
}

func GetGroupList(ownerId uint) (*[]models.Group, error) {
	// 获取我加入的群
	relation := make([]models.Relation, 0)
	if tx := global.DB.Where("owner_id = ? and type = ?", ownerId, 2).Find(&relation); tx.RowsAffected == 0 {
		return nil, errors.New("没有加入任何群组")
	}

	groupIds := make([]uint, 0)
	for _, v := range relation {
		groupIds = append(groupIds, v.TargetId)
	}

	group := make([]models.Group, 0)
	if tx := global.DB.Where("id in ?", groupIds).Find(&group); tx.RowsAffected == 0 {
		return nil, errors.New("群聊信息获取失败")
	}
	return &group, nil
}

func JoinGroup(ownerId uint, gName string) (int, error) {
	group := models.Group{}
	if tx := global.DB.Where("name = ?", gName).First(&group); tx.RowsAffected == 0 {
		return -1, errors.New("群记录不存在")
	}

	// 是否已经加入群聊
	relation := models.Relation{}
	if tx := global.DB.Where("owner_id = ? and target_id = ? and type = ?", ownerId, group.Id, 2).First(&relation); tx.RowsAffected != 0 {
		return -1, errors.New("已经加入该群聊")
	}

	relation = models.Relation{}
	relation.OwnerId = ownerId
	relation.TargetId = group.Id
	relation.Type = 2

	if tx := global.DB.Create(&relation); tx.RowsAffected == 0 {
		return -1, errors.New("加入群聊失败")
	}
	return 0, nil
}
