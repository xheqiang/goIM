package models

type Relation struct {
	Model
	OwnerId  uint   `gorm:"type:int(11);unsigned" json:"owner_id"`  // 谁的关系信息
	TargetId uint   `gorm:"type:int(11);unsigned" json:"target_id"` // 对应的谁
	Type     int    `gorm:"type:tinyint(3);" json:"type"`           // 关系类型 1:好友 2: 群
	Desc     string `gorm:"type:varchar(200);" json:"sesc"`         // 群的描述
}

func (r *Relation) TableName() string {
	return "relation"
}
