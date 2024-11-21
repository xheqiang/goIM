package models

type Group struct {
	Model
	Name    string `gorm:"type:varchar(25)" json:"name"`
	OwnerId uint   `gorm:"type:int(11); unsigned" json:"owner_id"`
	Type    int    `gorm:"type:int(11); unsigned" json:"type"`
	Image   string `gorm:"type:varchar(50);" json:"image"`
	Desc    string `gorm:"type:varchar(100);" json:"desc"`
}

func (r *Group) TableName() string {
	return "group"
}
