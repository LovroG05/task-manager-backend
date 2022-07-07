package models

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:255;not null" json:"title"`
	Description string `gorm:"size:255;not null" json:"description"`
	Creator     User   `json:"creator" gorm:"foreignkey:CreatorID"`
	CreatorID   uint   `json:"-"`
	Assignees   []User `gorm:"many2many:user_tasks;not null;constraint:OnDelete:SET NULL"`
	Time        string `json:"time"`
	Days        uint8  `json:"days"`
	Last        int    `json:"last_performer"`
	OneTime     bool   `json:"one_time"`
}
