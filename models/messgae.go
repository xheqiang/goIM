package models

type Message struct {
	Id        uint   `gorm:"primaryKey;type:int(11);unique" json:"id"`
	FromId    int64  `gorm:"type:int(11);unsigned" json:"from_id"`
	TargetId  int64  `gorm:"type:int(11);unsigned" json:"target_id"`
	Type      int    `gorm:"type:tinyint(3);unsigned" json:"type"`
	Media     int    `gorm:"type:tinyint(3);unsigned" json:"media"`
	Content   string `gorm:"type:text" json:"content"`
	Pic       string `gorm:"type:varchar(100)" json:"pic"`
	Url       string `gorm:"type:varchar(100)" json:"url"`
	Desc      string `gorm:"type:varchar(200)" json:"Desc"`
	Size      int    `gorm:"type:int(11);unsigned" json:"size"`
	CreatedAt string `gorm:"type:datetime" json:"created_at"`
	UpdatedAt string `gorm:"type:datetime" json:"updated_at"`
}

func (m *Message) TableName() string {
	return "message"
}
