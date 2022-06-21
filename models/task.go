package models

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:255;not null" json:"title"`
	Description string `gorm:"size:255;not null" json:"description"`
	Creator     User   `json:"creator" gorm:"foreignkey:CreatorID"`
	CreatorID   uint   `json:"-"`
	Assignees   []User `gorm:"many2many:user_tasks;"`
	Daily       bool   `json:"daily"`
	OnTime      string `json:"on_time"`
	OnDay       string `json:"on_day"`
}
