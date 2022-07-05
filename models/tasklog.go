package models

type TaskLog struct {
	ID     uint `gorm:"primaryKey"`
	TaskID int
	Task   Task
	UserID int
	User   User
	Time   string
	Notes  string
}
