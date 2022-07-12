package inputs

type TaskLogInput struct {
	TaskID int    `json:"task_id" binding:"required"`
	Time   string `json:"time"`
	Notes  string `json:"notes"`
}
